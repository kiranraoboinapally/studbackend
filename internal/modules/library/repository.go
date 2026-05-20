package librarymod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Authors
func (r *Repository) ListAuthors() ([]domain.Author, error) {
	var list []domain.Author
	return list, r.db.Order("name").Find(&list).Error
}
func (r *Repository) GetAuthor(id uint) (*domain.Author, error) {
	var a domain.Author
	return &a, r.db.First(&a, id).Error
}
func (r *Repository) CreateAuthor(a *domain.Author) error {
	return r.db.Create(a).Error
}
func (r *Repository) UpdateAuthor(a *domain.Author) error {
	return r.db.Save(a).Error
}

// Books
func (r *Repository) ListBooks(search string) ([]domain.Book, error) {
	var list []domain.Book
	q := r.db.Order("title")
	if search != "" {
		q = q.Where("title ILIKE ? OR isbn ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetBook(id uint) (*domain.Book, error) {
	var b domain.Book
	return &b, r.db.First(&b, id).Error
}
func (r *Repository) CreateBook(b *domain.Book) error {
	return r.db.Create(b).Error
}
func (r *Repository) UpdateBook(b *domain.Book) error {
	return r.db.Save(b).Error
}

// Book Copies
func (r *Repository) ListBookCopies(bookID uint) ([]domain.BookCopy, error) {
	var list []domain.BookCopy
	q := r.db.Order("copy_number")
	if bookID > 0 {
		q = q.Where("book_id = ?", bookID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetBookCopy(id uint) (*domain.BookCopy, error) {
	var bc domain.BookCopy
	return &bc, r.db.First(&bc, id).Error
}
func (r *Repository) CreateBookCopy(bc *domain.BookCopy) error {
	return r.db.Create(bc).Error
}
func (r *Repository) UpdateBookCopy(bc *domain.BookCopy) error {
	return r.db.Save(bc).Error
}

// Book Authors
func (r *Repository) AddBookAuthor(ba *domain.BookAuthor) error {
	return r.db.Save(ba).Error
}
func (r *Repository) RemoveBookAuthor(bookID, authorID uint) error {
	return r.db.Where("book_id = ? AND author_id = ?", bookID, authorID).Delete(&domain.BookAuthor{}).Error
}
func (r *Repository) GetBookAuthors(bookID uint) ([]domain.Author, error) {
	var list []domain.Author
	return list, r.db.Raw(`SELECT a.* FROM library.authors a JOIN library.book_authors ba ON ba.author_id = a.id WHERE ba.book_id = ?`, bookID).Scan(&list).Error
}

// Digital Resources
func (r *Repository) ListDigitalResources(search string) ([]domain.DigitalResource, error) {
	var list []domain.DigitalResource
	q := r.db.Order("title")
	if search != "" {
		q = q.Where("title ILIKE ?", "%"+search+"%")
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetDigitalResource(id uint) (*domain.DigitalResource, error) {
	var dr domain.DigitalResource
	return &dr, r.db.First(&dr, id).Error
}
func (r *Repository) CreateDigitalResource(dr *domain.DigitalResource) error {
	return r.db.Create(dr).Error
}

// Circulations
func (r *Repository) CreateCirculation(c *domain.Circulation) error {
	return r.db.Create(c).Error
}
func (r *Repository) GetCirculation(id uint) (*domain.Circulation, error) {
	var c domain.Circulation
	return &c, r.db.First(&c, id).Error
}
func (r *Repository) GetActiveCirculation(bookCopyID, studentID uint) (*domain.Circulation, error) {
	var c domain.Circulation
	return &c, r.db.Where("book_copy_id = ? AND student_id = ? AND returned_date IS NULL", bookCopyID, studentID).First(&c).Error
}
func (r *Repository) UpdateCirculation(c *domain.Circulation) error {
	return r.db.Save(c).Error
}
func (r *Repository) ListActiveCirculations(studentID uint) ([]domain.Circulation, error) {
	var list []domain.Circulation
	q := r.db.Where("returned_date IS NULL").Order("due_date")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) ListOverdueCirculations() ([]domain.Circulation, error) {
	var list []domain.Circulation
	return list, r.db.Where("returned_date IS NULL AND due_date < CURRENT_DATE AND fine_amount = 0").Find(&list).Error
}

// Reservations
func (r *Repository) CreateReservation(res *domain.Reservation) error {
	return r.db.Create(res).Error
}
func (r *Repository) ListReservations(bookID, studentID uint) ([]domain.Reservation, error) {
	var list []domain.Reservation
	q := r.db.Order("reserved_from DESC")
	if bookID > 0 {
		q = q.Where("book_id = ?", bookID)
	}
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateReservation(res *domain.Reservation) error {
	return r.db.Save(res).Error
}

// Fines
func (r *Repository) CreateFine(f *domain.LibraryFine) error {
	return r.db.Create(f).Error
}
func (r *Repository) ListFines(studentID uint) ([]domain.LibraryFine, error) {
	var list []domain.LibraryFine
	q := r.db.Order("created_at DESC")
	if studentID > 0 {
		q = q.Where("circulation_id IN (SELECT id FROM library.circulations WHERE student_id = ?)", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateFine(f *domain.LibraryFine) error {
	return r.db.Save(f).Error
}

// Purchase Requests
func (r *Repository) CreatePurchaseRequest(pr *domain.PurchaseRequest) error {
	return r.db.Create(pr).Error
}
func (r *Repository) ListPurchaseRequests(statusID *uint) ([]domain.PurchaseRequest, error) {
	var list []domain.PurchaseRequest
	q := r.db.Order("created_at DESC")
	if statusID != nil {
		q = q.Where("status_id = ?", *statusID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdatePurchaseRequest(pr *domain.PurchaseRequest) error {
	return r.db.Save(pr).Error
}

// Stats
func (r *Repository) GetLibraryStats() (map[string]interface{}, error) {
	var totalBooks, totalCopies, activeCirculations, overdueCount int64
	r.db.Model(&domain.Book{}).Count(&totalBooks)
	r.db.Model(&domain.BookCopy{}).Count(&totalCopies)
	r.db.Model(&domain.Circulation{}).Where("returned_date IS NULL").Count(&activeCirculations)
	r.db.Model(&domain.Circulation{}).Where("returned_date IS NULL AND due_date < CURRENT_DATE").Count(&overdueCount)
	return map[string]interface{}{
		"total_books":          totalBooks,
		"total_copies":         totalCopies,
		"active_circulations":  activeCirculations,
		"overdue_circulations": overdueCount,
	}, nil
}
