package studentmod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

// Repository handles student database operations.
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(student *domain.Student) error {
	return r.db.Create(student).Error
}

func (r *Repository) FindByID(id uint) (*domain.Student, error) {
	var s domain.Student
	err := r.db.First(&s, id).Error
	return &s, err
}

func (r *Repository) FindByUserID(userID uint) (*domain.Student, error) {
	var s domain.Student
	err := r.db.Where("user_id = ?", userID).First(&s).Error
	return &s, err
}

func (r *Repository) FindByEnrollment(enrollment string) (*domain.Student, error) {
	var s domain.Student
	err := r.db.Where("enrollment_number = ?", enrollment).First(&s).Error
	return &s, err
}

func (r *Repository) ListByProgram(programID uint, page, pageSize int) ([]domain.Student, int64, error) {
	var students []domain.Student
	var total int64
	query := r.db.Where("program_id = ?", programID)
	query.Model(&domain.Student{}).Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id ASC").Find(&students).Error
	return students, total, err
}

func (r *Repository) ListAll(page, pageSize int) ([]domain.Student, int64, error) {
	var students []domain.Student
	var total int64
	r.db.Model(&domain.Student{}).Count(&total)
	err := r.db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id ASC").Find(&students).Error
	return students, total, err
}

func (r *Repository) Update(student *domain.Student) error {
	return r.db.Save(student).Error
}

func (r *Repository) CreateGuardian(g *domain.Guardian) error {
	return r.db.Create(g).Error
}

func (r *Repository) GetGuardians(studentID uint) ([]domain.Guardian, error) {
	var guardians []domain.Guardian
	err := r.db.Where("student_id = ?", studentID).Find(&guardians).Error
	return guardians, err
}

func (r *Repository) CreateMedicalRecord(m *domain.MedicalRecord) error {
	return r.db.Create(m).Error
}

func (r *Repository) GetMedicalRecord(studentID uint) (*domain.MedicalRecord, error) {
	var rec domain.MedicalRecord
	err := r.db.Where("student_id = ?", studentID).First(&rec).Error
	return &rec, err
}

func (r *Repository) CreateStatusHistory(h *domain.StudentStatusHistory) error {
	return r.db.Create(h).Error
}

func (r *Repository) GetStatusHistory(studentID uint) ([]domain.StudentStatusHistory, error) {
	var history []domain.StudentStatusHistory
	err := r.db.Where("student_id = ?", studentID).Order("effective_from DESC").Find(&history).Error
	return history, err
}

func (r *Repository) CreateGrievance(g *domain.Grievance) error {
	return r.db.Create(g).Error
}

func (r *Repository) GetGrievances(studentID uint) ([]domain.Grievance, error) {
	var grievances []domain.Grievance
	err := r.db.Where("student_id = ?", studentID).Order("created_at DESC").Find(&grievances).Error
	return grievances, err
}

// Dashboard aggregates student data for a dashboard view.
func (r *Repository) GetDashboard(studentID uint) (map[string]interface{}, error) {
	var student domain.Student
	if err := r.db.First(&student, studentID).Error; err != nil {
		return nil, err
	}

	var program domain.Program
	r.db.First(&program, student.ProgramID)

	var guardians []domain.Guardian
	r.db.Where("student_id = ?", studentID).Find(&guardians)

	var enrollmentCount int64
	r.db.Model(&domain.StudentEnrollment{}).Where("student_id = ?", studentID).Count(&enrollmentCount)

	var pendingInvoices int64
	r.db.Model(&domain.Invoice{}).Where("student_id = ? AND paid_amount < total_amount", studentID).Count(&pendingInvoices)

	return map[string]interface{}{
		"student":          student,
		"program":          program,
		"guardians":        guardians,
		"enrollment_count": enrollmentCount,
		"pending_invoices": pendingInvoices,
	}, nil
}
