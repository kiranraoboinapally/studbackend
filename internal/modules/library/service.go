package librarymod

import (
	"context"
	"fmt"
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
	return &Service{repo: repo, bus: bus, outbox: ob, db: db}
}

// Authors
func (s *Service) ListAuthors(ctx context.Context) ([]domain.Author, error) {
	return s.repo.ListAuthors()
}
func (s *Service) GetAuthor(ctx context.Context, id uint) (*domain.Author, error) {
	a, err := s.repo.GetAuthor(id)
	if err != nil {
		return nil, apperrors.NotFound("author not found")
	}
	return a, nil
}
func (s *Service) CreateAuthor(ctx context.Context, a *domain.Author) error {
	if a.Name == "" {
		return apperrors.BadRequest("author name is required")
	}
	return s.repo.CreateAuthor(a)
}
func (s *Service) UpdateAuthor(ctx context.Context, id uint, a *domain.Author) error {
	existing, err := s.repo.GetAuthor(id)
	if err != nil {
		return apperrors.NotFound("author not found")
	}
	a.ID = existing.ID
	return s.repo.UpdateAuthor(a)
}

// Books
func (s *Service) ListBooks(ctx context.Context, search string) ([]domain.Book, error) {
	return s.repo.ListBooks(search)
}
func (s *Service) GetBook(ctx context.Context, id uint) (*domain.Book, error) {
	b, err := s.repo.GetBook(id)
	if err != nil {
		return nil, apperrors.NotFound("book not found")
	}
	return b, nil
}
func (s *Service) CreateBook(ctx context.Context, b *domain.Book) error {
	if b.Title == "" {
		return apperrors.BadRequest("book title is required")
	}
	return s.repo.CreateBook(b)
}
func (s *Service) UpdateBook(ctx context.Context, id uint, b *domain.Book) error {
	existing, err := s.repo.GetBook(id)
	if err != nil {
		return apperrors.NotFound("book not found")
	}
	b.ID = existing.ID
	return s.repo.UpdateBook(b)
}

// Book Copies
func (s *Service) ListBookCopies(ctx context.Context, bookID uint) ([]domain.BookCopy, error) {
	return s.repo.ListBookCopies(bookID)
}
func (s *Service) CreateBookCopy(ctx context.Context, bc *domain.BookCopy) error {
	if bc.BookID == 0 {
		return apperrors.BadRequest("book_id is required")
	}
	return s.repo.CreateBookCopy(bc)
}

// Book Authors
func (s *Service) AddBookAuthor(ctx context.Context, ba *domain.BookAuthor) error {
	return s.repo.AddBookAuthor(ba)
}
func (s *Service) RemoveBookAuthor(ctx context.Context, bookID, authorID uint) error {
	return s.repo.RemoveBookAuthor(bookID, authorID)
}
func (s *Service) GetBookAuthors(ctx context.Context, bookID uint) ([]domain.Author, error) {
	return s.repo.GetBookAuthors(bookID)
}

// Digital Resources
func (s *Service) ListDigitalResources(ctx context.Context, search string) ([]domain.DigitalResource, error) {
	return s.repo.ListDigitalResources(search)
}
func (s *Service) CreateDigitalResource(ctx context.Context, dr *domain.DigitalResource) error {
	if dr.Title == "" {
		return apperrors.BadRequest("title is required")
	}
	return s.repo.CreateDigitalResource(dr)
}

// Circulations - Issue book with outbox event
func (s *Service) IssueBook(ctx context.Context, bookCopyID, studentID uint, issuedBy *uint) (*domain.Circulation, error) {
	if bookCopyID == 0 || studentID == 0 {
		return nil, apperrors.BadRequest("book_copy_id and student_id are required")
	}

	// Check if copy is available
	copy, err := s.repo.GetBookCopy(bookCopyID)
	if err != nil {
		return nil, apperrors.NotFound("book copy not found")
	}

	// Check for existing active circulation
	existing, _ := s.repo.GetActiveCirculation(bookCopyID, 0)
	if existing != nil {
		return nil, apperrors.Conflict("book copy is already issued")
	}

	var issuedStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "library", "ISSUED").First(&issuedStatus)

	dueDate := time.Now().AddDate(0, 0, 14) // 14 days
	circ := &domain.Circulation{
		BookCopyID: bookCopyID,
		StudentID:  studentID,
		IssuedDate: time.Now(),
		DueDate:    dueDate,
		StatusID:   &issuedStatus.ID,
		IssuedBy:   issuedBy,
	}

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(circ).Error; err != nil {
			return err
		}
		// Update book available copies
		tx.Model(&domain.Book{}).Where("id = ?", copy.BookID).Update("available_copies", gorm.Expr("available_copies - 1"))
		// Update copy status
		var checkedOutStatus domain.StatusCode
		tx.Where("module = ? AND code = ?", "library", "CHECKED_OUT").First(&checkedOutStatus)
		copy.StatusID = &checkedOutStatus.ID
		tx.Save(copy)

		return s.outbox.WriteEvent(tx, "Circulation", fmt.Sprintf("%d", circ.ID),
			eventbus.EventBookIssued,
			map[string]interface{}{
				"circulation_id": circ.ID,
				"book_copy_id":   bookCopyID,
				"student_id":     studentID,
				"due_date":       dueDate,
			},
		)
	})
	if txErr != nil {
		return nil, txErr
	}
	return circ, nil
}

// ReturnBook handles book return with fine calculation and overdue event
func (s *Service) ReturnBook(ctx context.Context, circulationID uint) (*domain.Circulation, error) {
	circ, err := s.repo.GetCirculation(circulationID)
	if err != nil {
		return nil, apperrors.NotFound("circulation not found")
	}
	if circ.ReturnedDate != nil {
		return nil, apperrors.Conflict("book already returned")
	}

	now := time.Now()
	circ.ReturnedDate = &now

	// Calculate fine if overdue
	daysOverdue := 0
	if now.After(circ.DueDate) {
		daysOverdue = int(now.Sub(circ.DueDate).Hours() / 24)
		circ.FineAmount = float64(daysOverdue) * 5.0 // 5 per day
	}

	var returnedStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "library", "RETURNED").First(&returnedStatus)
	circ.StatusID = &returnedStatus.ID

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(circ).Error; err != nil {
			return err
		}
		// Update book available copies
		var bc domain.BookCopy
		if err := tx.First(&bc, circ.BookCopyID).Error; err == nil {
			tx.Model(&domain.Book{}).Where("id = ?", bc.BookID).Update("available_copies", gorm.Expr("available_copies + 1"))
			var availableStatus domain.StatusCode
			tx.Where("module = ? AND code = ?", "library", "AVAILABLE").First(&availableStatus)
			bc.StatusID = &availableStatus.ID
			tx.Save(&bc)
		}

		// Create fine record if overdue
		if circ.FineAmount > 0 {
			fine := &domain.LibraryFine{
				CirculationID: circ.ID,
				Amount:        circ.FineAmount,
				Reason:        fmt.Sprintf("Overdue by %d days", daysOverdue),
			}
			if err := tx.Create(fine).Error; err != nil {
				return err
			}

			// Emit overdue event for finance module to create fine invoice
			return s.outbox.WriteEvent(tx, "Circulation", fmt.Sprintf("%d", circ.ID),
				eventbus.EventBookOverdue,
				eventbus.BookOverduePayload{
					CirculationID: circ.ID,
					StudentID:     circ.StudentID,
					BookCopyID:    circ.BookCopyID,
					FineAmount:    circ.FineAmount,
					DaysOverdue:   daysOverdue,
				},
			)
		}

		return s.outbox.WriteEvent(tx, "Circulation", fmt.Sprintf("%d", circ.ID),
			eventbus.EventBookReturned,
			map[string]interface{}{
				"circulation_id": circ.ID,
				"student_id":     circ.StudentID,
				"book_copy_id":   circ.BookCopyID,
			},
		)
	})
	if txErr != nil {
		return nil, txErr
	}
	return circ, nil
}

// MarkOverdue scans for overdue books and emits events
func (s *Service) MarkOverdue(ctx context.Context) error {
	overdues, err := s.repo.ListOverdueCirculations()
	if err != nil {
		return err
	}

	for _, circ := range overdues {
		daysOverdue := int(time.Now().Sub(circ.DueDate).Hours() / 24)
		fineAmount := float64(daysOverdue) * 5.0

		s.db.Model(&circ).Updates(map[string]interface{}{
			"fine_amount": fineAmount,
		})

		s.bus.PublishAsync(ctx, eventbus.Event{
			Type:          eventbus.EventBookOverdue,
			AggregateType: "Circulation",
			AggregateID:   fmt.Sprintf("%d", circ.ID),
			Payload: eventbus.BookOverduePayload{
				CirculationID: circ.ID,
				StudentID:     circ.StudentID,
				BookCopyID:    circ.BookCopyID,
				FineAmount:    fineAmount,
				DaysOverdue:   daysOverdue,
			},
		})
	}
	return nil
}

func (s *Service) ListActiveCirculations(ctx context.Context, studentID uint) ([]domain.Circulation, error) {
	return s.repo.ListActiveCirculations(studentID)
}

// Reservations
func (s *Service) CreateReservation(ctx context.Context, res *domain.Reservation) error {
	if res.BookID == 0 || res.StudentID == 0 {
		return apperrors.BadRequest("book_id and student_id are required")
	}
	res.ReservedFrom = time.Now()
	return s.repo.CreateReservation(res)
}
func (s *Service) ListReservations(ctx context.Context, bookID, studentID uint) ([]domain.Reservation, error) {
	return s.repo.ListReservations(bookID, studentID)
}

// Fines
func (s *Service) ListFines(ctx context.Context, studentID uint) ([]domain.LibraryFine, error) {
	return s.repo.ListFines(studentID)
}
func (s *Service) PayFine(ctx context.Context, fineID uint) error {
	var fine domain.LibraryFine
	if err := s.db.First(&fine, fineID).Error; err != nil {
		return apperrors.NotFound("fine not found")
	}
	now := time.Now()
	fine.PaidDate = &now
	return s.repo.UpdateFine(&fine)
}

// Purchase Requests
func (s *Service) CreatePurchaseRequest(ctx context.Context, pr *domain.PurchaseRequest) error {
	if pr.Title == "" {
		return apperrors.BadRequest("title is required")
	}
	return s.repo.CreatePurchaseRequest(pr)
}
func (s *Service) ListPurchaseRequests(ctx context.Context) ([]domain.PurchaseRequest, error) {
	return s.repo.ListPurchaseRequests(nil)
}
func (s *Service) ApprovePurchaseRequest(ctx context.Context, id, approvedBy uint) error {
	pr, err := s.repo.ListPurchaseRequests(nil)
	_ = pr
	if err != nil {
		return apperrors.NotFound("purchase request not found")
	}
	now := time.Now()
	var req domain.PurchaseRequest
	s.db.First(&req, id)
	req.ApprovedBy = &approvedBy
	req.ApprovedAt = &now
	return s.repo.UpdatePurchaseRequest(&req)
}

// Stats
func (s *Service) GetStats(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetLibraryStats()
}
