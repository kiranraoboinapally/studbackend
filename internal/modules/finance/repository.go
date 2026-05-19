package financemod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Fee Heads
func (r *Repository) ListFeeHeads() ([]domain.FeeHead, error) {
	var list []domain.FeeHead
	return list, r.db.Order("name").Find(&list).Error
}
func (r *Repository) GetFeeHead(id uint) (*domain.FeeHead, error) {
	var h domain.FeeHead
	return &h, r.db.First(&h, id).Error
}
func (r *Repository) CreateFeeHead(h *domain.FeeHead) error {
	return r.db.Create(h).Error
}
func (r *Repository) UpdateFeeHead(h *domain.FeeHead) error {
	return r.db.Save(h).Error
}

// Fee Structures
func (r *Repository) ListFeeStructures(programID uint, academicYear string) ([]domain.FeeStructure, error) {
	var list []domain.FeeStructure
	q := r.db.Where("is_active = true")
	if programID > 0 {
		q = q.Where("program_id = ?", programID)
	}
	if academicYear != "" {
		q = q.Where("academic_year = ?", academicYear)
	}
	return list, q.Order("semester_number, fee_head_id").Find(&list).Error
}
func (r *Repository) GetFeeStructuresForProgram(programID uint) ([]domain.FeeStructure, error) {
	var list []domain.FeeStructure
	return list, r.db.Where("program_id = ? AND is_active = true", programID).Find(&list).Error
}
func (r *Repository) CreateFeeStructure(fs *domain.FeeStructure) error {
	return r.db.Create(fs).Error
}
func (r *Repository) UpdateFeeStructure(fs *domain.FeeStructure) error {
	return r.db.Save(fs).Error
}

// Invoices
func (r *Repository) GetStudentInvoices(studentID uint) ([]domain.Invoice, error) {
	var list []domain.Invoice
	return list, r.db.Where("student_id = ?", studentID).Order("created_at DESC").Find(&list).Error
}
func (r *Repository) GetInvoiceByID(id uint) (*domain.Invoice, error) {
	var inv domain.Invoice
	return &inv, r.db.First(&inv, id).Error
}
func (r *Repository) GetInvoiceItems(invoiceID uint) ([]domain.InvoiceItem, error) {
	var list []domain.InvoiceItem
	return list, r.db.Where("invoice_id = ?", invoiceID).Find(&list).Error
}
func (r *Repository) CreateInvoice(invoice *domain.Invoice, items []domain.InvoiceItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].InvoiceID = invoice.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *Repository) ListAllInvoices(page, pageSize int) ([]domain.Invoice, int64, error) {
	var list []domain.Invoice
	var total int64
	r.db.Model(&domain.Invoice{}).Count(&total)
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	return list, total, r.db.Offset((page-1)*pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
}

// Payments
func (r *Repository) CreatePayment(payment *domain.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		var inv domain.Invoice
		if err := tx.First(&inv, payment.InvoiceID).Error; err != nil {
			return err
		}
		inv.PaidAmount += payment.Amount
		status := "PARTIAL"
		if inv.PaidAmount >= inv.TotalAmount {
			status = "PAID"
		}
		var sc domain.StatusCode
		if err := tx.Where("module = ? AND code = ?", "finance", status).First(&sc).Error; err == nil {
			inv.StatusID = &sc.ID
		}
		if err := tx.Save(&inv).Error; err != nil {
			return err
		}
		alloc := domain.PaymentAllocation{
			PaymentID:       payment.ID,
			InvoiceID:       payment.InvoiceID,
			AllocatedAmount: payment.Amount,
		}
		return tx.Create(&alloc).Error
	})
}
func (r *Repository) ListPayments(studentID uint) ([]domain.Payment, error) {
	var list []domain.Payment
	q := r.db.Order("payment_date DESC")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetPayment(id uint) (*domain.Payment, error) {
	var p domain.Payment
	return &p, r.db.First(&p, id).Error
}

// Scholarships
func (r *Repository) ListScholarships() ([]domain.Scholarship, error) {
	var list []domain.Scholarship
	return list, r.db.Order("name").Find(&list).Error
}
func (r *Repository) GetScholarship(id uint) (*domain.Scholarship, error) {
	var s domain.Scholarship
	return &s, r.db.First(&s, id).Error
}
func (r *Repository) CreateScholarship(s *domain.Scholarship) error {
	return r.db.Create(s).Error
}
func (r *Repository) UpdateScholarship(s *domain.Scholarship) error {
	return r.db.Save(s).Error
}
func (r *Repository) GetStudentScholarships(studentID uint) ([]domain.StudentScholarship, error) {
	var list []domain.StudentScholarship
	return list, r.db.Where("student_id = ?", studentID).Find(&list).Error
}
func (r *Repository) AssignScholarship(ss *domain.StudentScholarship) error {
	return r.db.Create(ss).Error
}

// Discounts
func (r *Repository) GetStudentDiscounts(studentID uint) ([]domain.StudentDiscount, error) {
	var list []domain.StudentDiscount
	return list, r.db.Where("student_id = ?", studentID).Find(&list).Error
}
func (r *Repository) CreateDiscount(d *domain.StudentDiscount) error {
	return r.db.Create(d).Error
}

// Installments
func (r *Repository) GetStudentInstallments(studentID, termID uint) ([]domain.InstallmentPlan, error) {
	var list []domain.InstallmentPlan
	q := r.db.Where("student_id = ?", studentID)
	if termID > 0 {
		q = q.Where("academic_term_id = ?", termID)
	}
	return list, q.Order("due_date").Find(&list).Error
}
func (r *Repository) CreateInstallment(ip *domain.InstallmentPlan) error {
	return r.db.Create(ip).Error
}

// Refunds
func (r *Repository) CreateRefund(ref *domain.Refund) error {
	return r.db.Create(ref).Error
}
func (r *Repository) GetRefund(id uint) (*domain.Refund, error) {
	var ref domain.Refund
	return &ref, r.db.First(&ref, id).Error
}
func (r *Repository) UpdateRefund(ref *domain.Refund) error {
	return r.db.Save(ref).Error
}
func (r *Repository) ListRefunds(studentID uint) ([]domain.Refund, error) {
	var list []domain.Refund
	q := r.db.Order("created_at DESC")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}

// Finance summary
func (r *Repository) GetFinanceSummary() (map[string]interface{}, error) {
	var totalInvoiced, totalCollected float64
	r.db.Model(&domain.Invoice{}).Select("COALESCE(SUM(total_amount),0)").Scan(&totalInvoiced)
	r.db.Model(&domain.Invoice{}).Select("COALESCE(SUM(paid_amount),0)").Scan(&totalCollected)
	var pendingCount int64
	r.db.Model(&domain.Invoice{}).Where("paid_amount < total_amount").Count(&pendingCount)
	return map[string]interface{}{
		"total_invoiced":    totalInvoiced,
		"total_collected":   totalCollected,
		"outstanding":       totalInvoiced - totalCollected,
		"pending_invoices":  pendingCount,
	}, nil
}
