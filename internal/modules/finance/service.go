package financemod

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
	"university-erp-backend/internal/platform/eventbus"
	"university-erp-backend/internal/platform/outbox"

	"gorm.io/gorm"
)

type Service struct {
	repo   *Repository
	bus    *eventbus.Bus
	outbox *outbox.Writer
	db     *gorm.DB
}

func NewService(repo *Repository, bus *eventbus.Bus, ob *outbox.Writer, db *gorm.DB) *Service {
	s := &Service{repo: repo, bus: bus, outbox: ob, db: db}
	bus.Subscribe(eventbus.EventStudentEnrolled, s.HandleStudentEnrolled)
	bus.Subscribe(eventbus.EventTermRegistered, s.HandleTermRegistered)
	bus.Subscribe(eventbus.EventBookOverdue, s.HandleBookOverdue)
	bus.Subscribe(eventbus.EventPayrollProcessed, s.HandlePayrollProcessed)
	return s
}

// HandleStudentEnrolled auto-generates an invoice when a student enrolls.
func (s *Service) HandleStudentEnrolled(ctx context.Context, evt eventbus.Event) error {
	slog.Info("FinanceMod: received StudentEnrolled event", "aggregate_id", evt.AggregateID)

	payloadBytes, _ := evt.Payload.(string)
	var payload eventbus.StudentEnrolledPayload
	if err := json.Unmarshal([]byte(payloadBytes), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return s.generateInvoiceForEnrollment(payload.StudentID, payload.TermID, payload.ProgramID, payload.RollNumber)
}

// HandleTermRegistered generates invoices when a student registers for a new term.
func (s *Service) HandleTermRegistered(ctx context.Context, evt eventbus.Event) error {
	slog.Info("FinanceMod: received TermRegistered event", "aggregate_id", evt.AggregateID)

	payloadBytes, _ := evt.Payload.(string)
	var payload eventbus.TermRegisteredPayload
	if err := json.Unmarshal([]byte(payloadBytes), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Check if invoice already exists for this student+term
	var existing int64
	s.db.Model(&domain.Invoice{}).Where("student_id = ? AND academic_term_id = ?", payload.StudentID, payload.AcademicTermID).Count(&existing)
	if existing > 0 {
		slog.Info("FinanceMod: invoice already exists for student+term, skipping", "student_id", payload.StudentID, "term_id", payload.AcademicTermID)
		return nil
	}

	var student domain.Student
	if err := s.db.First(&student, payload.StudentID).Error; err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	return s.generateInvoiceForEnrollment(payload.StudentID, payload.AcademicTermID, payload.ProgramID, student.RollNumber)
}

// HandleBookOverdue creates a fine invoice when a library book is overdue.
func (s *Service) HandleBookOverdue(ctx context.Context, evt eventbus.Event) error {
	slog.Info("FinanceMod: received BookOverdue event", "aggregate_id", evt.AggregateID)

	payloadBytes, _ := evt.Payload.(string)
	var payload eventbus.BookOverduePayload
	if err := json.Unmarshal([]byte(payloadBytes), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return s.generateFineInvoice(payload.StudentID, payload.FineAmount, "Library overdue fine", payload.CirculationID)
}

// HandlePayrollProcessed creates payment vouchers in the finance ledger when payroll is processed.
func (s *Service) HandlePayrollProcessed(ctx context.Context, evt eventbus.Event) error {
	slog.Info("FinanceMod: received PayrollProcessed event", "aggregate_id", evt.AggregateID)

	payloadBytes, _ := evt.Payload.(string)
	var payload eventbus.PayrollProcessedPayload
	if err := json.Unmarshal([]byte(payloadBytes), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	slog.Info("FinanceMod: payroll voucher recorded",
		"employee_id", payload.EmployeeID,
		"net_pay", payload.NetPay,
		"month", payload.Month,
	)
	return nil
}

func (s *Service) generateInvoiceForEnrollment(studentID, termID, programID uint, rollNumber string) error {
	structures, err := s.repo.GetFeeStructuresForProgram(programID)
	if err != nil || len(structures) == 0 {
		slog.Warn("FinanceMod: no fee structure found for program", "program_id", programID)
		return nil
	}

	var totalAmount float64
	var items []domain.InvoiceItem
	for _, fs := range structures {
		head, _ := s.repo.GetFeeHead(fs.FeeHeadID)
		desc := "Fee"
		if head != nil {
			desc = head.Name
		}
		items = append(items, domain.InvoiceItem{
			FeeHeadID:   fs.FeeHeadID,
			Description: desc,
			Quantity:    1,
			UnitAmount:  fs.Amount,
			Amount:      fs.Amount,
		})
		totalAmount += fs.Amount
	}

	var unPaidStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "finance", "UNPAID").First(&unPaidStatus)

	invNo := fmt.Sprintf("INV-%d-%s", time.Now().Year(), rollNumber)
	invoice := domain.Invoice{
		StudentID:      studentID,
		InvoiceNumber:  invNo,
		AcademicTermID: termID,
		DueDate:        time.Now().AddDate(0, 1, 0),
		TotalAmount:    totalAmount,
		PaidAmount:     0,
		StatusID:       &unPaidStatus.ID,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].InvoiceID = invoice.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return s.outbox.WriteEvent(tx, "Invoice", fmt.Sprintf("%d", invoice.ID),
			eventbus.EventInvoiceGenerated,
			eventbus.InvoiceGeneratedPayload{
				InvoiceID:     invoice.ID,
				StudentID:     invoice.StudentID,
				TotalAmount:   invoice.TotalAmount,
				InvoiceNumber: invoice.InvoiceNumber,
			},
		)
	})

	if err != nil {
		slog.Error("FinanceMod: failed to generate invoice", "student_id", studentID, "error", err)
		return err
	}
	slog.Info("FinanceMod: generated invoice", "invoice_number", invNo, "student_id", studentID)
	return nil
}

func (s *Service) generateFineInvoice(studentID uint, amount float64, reason string, sourceID uint) error {
	var unPaidStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "finance", "UNPAID").First(&unPaidStatus)

	var term domain.AcademicTerm
	s.db.Where("is_current = true").First(&term)

	invNo := fmt.Sprintf("FINE-%d-%d", time.Now().Unix(), studentID)
	invoice := domain.Invoice{
		StudentID:      studentID,
		InvoiceNumber:  invNo,
		AcademicTermID: term.ID,
		DueDate:        time.Now().AddDate(0, 0, 7),
		TotalAmount:    amount,
		PaidAmount:     0,
		StatusID:       &unPaidStatus.ID,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}
		item := domain.InvoiceItem{
			InvoiceID:   invoice.ID,
			Description: reason,
			Quantity:    1,
			UnitAmount:  amount,
			Amount:      amount,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "Fine", fmt.Sprintf("%d", invoice.ID),
			eventbus.EventFineGenerated,
			eventbus.FineGeneratedPayload{
				StudentID:    studentID,
				Amount:       amount,
				Reason:       reason,
				SourceModule: "library",
				SourceID:     sourceID,
			},
		)
	})
}

// Fee Heads
func (s *Service) ListFeeHeads(ctx context.Context) ([]domain.FeeHead, error) {
	return s.repo.ListFeeHeads()
}
func (s *Service) CreateFeeHead(ctx context.Context, h *domain.FeeHead) error {
	if h.Name == "" || h.Code == "" {
		return apperrors.BadRequest("fee head name and code are required")
	}
	return s.repo.CreateFeeHead(h)
}
func (s *Service) UpdateFeeHead(ctx context.Context, id uint, h *domain.FeeHead) error {
	existing, err := s.repo.GetFeeHead(id)
	if err != nil {
		return apperrors.NotFound("fee head not found")
	}
	h.ID = existing.ID
	return s.repo.UpdateFeeHead(h)
}

// Fee Structures
func (s *Service) ListFeeStructures(ctx context.Context, programID uint, academicYear string) ([]domain.FeeStructure, error) {
	return s.repo.ListFeeStructures(programID, academicYear)
}
func (s *Service) CreateFeeStructure(ctx context.Context, fs *domain.FeeStructure) error {
	if fs.ProgramID == 0 || fs.FeeHeadID == 0 {
		return apperrors.BadRequest("program_id and fee_head_id are required")
	}
	return s.repo.CreateFeeStructure(fs)
}
func (s *Service) UpdateFeeStructure(ctx context.Context, id uint, fs *domain.FeeStructure) error {
	fs.ID = id
	return s.repo.UpdateFeeStructure(fs)
}

// Invoices
func (s *Service) GetStudentInvoices(ctx context.Context, studentID uint) ([]domain.Invoice, error) {
	return s.repo.GetStudentInvoices(studentID)
}
func (s *Service) GetInvoiceWithItems(ctx context.Context, id uint) (map[string]interface{}, error) {
	inv, err := s.repo.GetInvoiceByID(id)
	if err != nil {
		return nil, apperrors.NotFound("invoice not found")
	}
	items, _ := s.repo.GetInvoiceItems(id)
	return map[string]interface{}{
		"invoice": inv,
		"items":   items,
	}, nil
}
func (s *Service) ListAllInvoices(ctx context.Context, page, pageSize int) ([]domain.Invoice, int64, error) {
	return s.repo.ListAllInvoices(page, pageSize)
}
func (s *Service) GenerateInvoiceForStudent(ctx context.Context, studentID, termID uint, programID uint) (*domain.Invoice, error) {
	structures, err := s.repo.GetFeeStructuresForProgram(programID)
	if err != nil || len(structures) == 0 {
		return nil, apperrors.BadRequest("no fee structures found for program")
	}
	var total float64
	var items []domain.InvoiceItem
	for _, fs := range structures {
		head, _ := s.repo.GetFeeHead(fs.FeeHeadID)
		desc := "Fee"
		if head != nil {
			desc = head.Name
		}
		items = append(items, domain.InvoiceItem{
			FeeHeadID:   fs.FeeHeadID,
			Description: desc,
			Quantity:    1,
			UnitAmount:  fs.Amount,
			Amount:      fs.Amount,
		})
		total += fs.Amount
	}
	var unPaidStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "finance", "UNPAID").First(&unPaidStatus)
	inv := &domain.Invoice{
		StudentID:      studentID,
		InvoiceNumber:  fmt.Sprintf("INV-%d-%d-%d", time.Now().Year(), studentID, termID),
		AcademicTermID: termID,
		DueDate:        time.Now().AddDate(0, 1, 0),
		TotalAmount:    total,
		StatusID:       &unPaidStatus.ID,
	}
	if err := s.repo.CreateInvoice(inv, items); err != nil {
		return nil, err
	}
	return inv, nil
}

// Payments
type PaymentRequest struct {
	InvoiceID     uint    `json:"invoice_id"`
	StudentID     uint    `json:"student_id"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
	PaymentModeID *uint   `json:"payment_mode_id"`
}

func (s *Service) ProcessPayment(ctx context.Context, req PaymentRequest) (*domain.Payment, error) {
	inv, err := s.repo.GetInvoiceByID(req.InvoiceID)
	if err != nil {
		return nil, apperrors.NotFound("invoice not found")
	}
	if req.Amount <= 0 {
		return nil, apperrors.BadRequest("payment amount must be positive")
	}
	if inv.PaidAmount+req.Amount > inv.TotalAmount+0.01 {
		return nil, apperrors.BadRequest("payment amount exceeds invoice balance")
	}
	payment := &domain.Payment{
		InvoiceID:     req.InvoiceID,
		StudentID:     req.StudentID,
		Amount:        req.Amount,
		TransactionID: req.TransactionID,
		PaymentModeID: req.PaymentModeID,
		PaymentDate:   time.Now(),
	}
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		inv.PaidAmount += req.Amount
		status := "PARTIAL"
		if inv.PaidAmount >= inv.TotalAmount {
			status = "PAID"
		}
		var sc domain.StatusCode
		if err := tx.Where("module = ? AND code = ?", "finance", status).First(&sc).Error; err == nil {
			inv.StatusID = &sc.ID
		}
		if err := tx.Save(inv).Error; err != nil {
			return err
		}
		alloc := domain.PaymentAllocation{
			PaymentID:       payment.ID,
			InvoiceID:       payment.InvoiceID,
			AllocatedAmount: payment.Amount,
		}
		if err := tx.Create(&alloc).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "Payment", fmt.Sprintf("%d", payment.ID),
			eventbus.EventPaymentCompleted,
			eventbus.PaymentCompletedPayload{
				PaymentID:     payment.ID,
				InvoiceID:     payment.InvoiceID,
				StudentID:     payment.StudentID,
				Amount:        payment.Amount,
				TransactionID: payment.TransactionID,
			},
		)
	})
	if txErr != nil {
		return nil, apperrors.Internal("payment processing failed", txErr)
	}
	return payment, nil
}

func (s *Service) ListPayments(ctx context.Context, studentID uint) ([]domain.Payment, error) {
	return s.repo.ListPayments(studentID)
}
func (s *Service) GetPayment(ctx context.Context, id uint) (*domain.Payment, error) {
	p, err := s.repo.GetPayment(id)
	if err != nil {
		return nil, apperrors.NotFound("payment not found")
	}
	return p, nil
}

// Scholarships
func (s *Service) ListScholarships(ctx context.Context) ([]domain.Scholarship, error) {
	return s.repo.ListScholarships()
}
func (s *Service) CreateScholarship(ctx context.Context, sc *domain.Scholarship) error {
	if sc.Name == "" {
		return apperrors.BadRequest("scholarship name is required")
	}
	return s.repo.CreateScholarship(sc)
}
func (s *Service) UpdateScholarship(ctx context.Context, id uint, sc *domain.Scholarship) error {
	existing, err := s.repo.GetScholarship(id)
	if err != nil {
		return apperrors.NotFound("scholarship not found")
	}
	sc.ID = existing.ID
	return s.repo.UpdateScholarship(sc)
}
func (s *Service) GetStudentScholarships(ctx context.Context, studentID uint) ([]domain.StudentScholarship, error) {
	return s.repo.GetStudentScholarships(studentID)
}
func (s *Service) AssignScholarship(ctx context.Context, ss *domain.StudentScholarship) error {
	if ss.StudentID == 0 || ss.ScholarshipID == 0 {
		return apperrors.BadRequest("student_id and scholarship_id are required")
	}
	return s.repo.AssignScholarship(ss)
}

// Discounts
func (s *Service) GetStudentDiscounts(ctx context.Context, studentID uint) ([]domain.StudentDiscount, error) {
	return s.repo.GetStudentDiscounts(studentID)
}
func (s *Service) ApplyDiscount(ctx context.Context, d *domain.StudentDiscount) error {
	if d.StudentID == 0 || d.Amount <= 0 {
		return apperrors.BadRequest("student_id and positive amount are required")
	}
	return s.repo.CreateDiscount(d)
}

// Installments
func (s *Service) GetStudentInstallments(ctx context.Context, studentID, termID uint) ([]domain.InstallmentPlan, error) {
	return s.repo.GetStudentInstallments(studentID, termID)
}
func (s *Service) CreateInstallmentPlan(ctx context.Context, ip *domain.InstallmentPlan) error {
	if ip.StudentID == 0 || ip.Amount <= 0 {
		return apperrors.BadRequest("student_id and amount are required")
	}
	return s.repo.CreateInstallment(ip)
}

// Refunds
func (s *Service) RequestRefund(ctx context.Context, ref *domain.Refund) error {
	if ref.PaymentID == 0 || ref.Amount <= 0 {
		return apperrors.BadRequest("payment_id and amount are required")
	}
	return s.repo.CreateRefund(ref)
}
func (s *Service) ApproveRefund(ctx context.Context, id, approvedBy uint) error {
	ref, err := s.repo.GetRefund(id)
	if err != nil {
		return apperrors.NotFound("refund not found")
	}
	now := time.Now()
	ref.ApprovedBy = &approvedBy
	ref.ProcessedAt = &now
	return s.repo.UpdateRefund(ref)
}
func (s *Service) ListRefunds(ctx context.Context, studentID uint) ([]domain.Refund, error) {
	return s.repo.ListRefunds(studentID)
}

// Summary
func (s *Service) GetSummary(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetFinanceSummary()
}
