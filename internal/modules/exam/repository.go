package exammod

import (
	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Exam Components
func (r *Repository) ListComponents(subjectID uint) ([]domain.ExamComponent, error) {
	var list []domain.ExamComponent
	q := r.db.Order("display_order")
	if subjectID > 0 {
		q = q.Where("subject_id = ?", subjectID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetComponent(id uint) (*domain.ExamComponent, error) {
	var c domain.ExamComponent
	return &c, r.db.First(&c, id).Error
}
func (r *Repository) CreateComponent(c *domain.ExamComponent) error {
	return r.db.Create(c).Error
}
func (r *Repository) UpdateComponent(c *domain.ExamComponent) error {
	return r.db.Save(c).Error
}

// Exam Schedules
func (r *Repository) ListSchedules(termID, subjectID uint) ([]domain.ExamSchedule, error) {
	var list []domain.ExamSchedule
	q := r.db.Order("exam_date")
	if termID > 0 {
		q = q.Where("academic_term_id = ?", termID)
	}
	if subjectID > 0 {
		q = q.Where("subject_id = ?", subjectID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetSchedule(id uint) (*domain.ExamSchedule, error) {
	var s domain.ExamSchedule
	return &s, r.db.First(&s, id).Error
}
func (r *Repository) CreateSchedule(s *domain.ExamSchedule) error {
	return r.db.Create(s).Error
}
func (r *Repository) UpdateSchedule(s *domain.ExamSchedule) error {
	return r.db.Save(s).Error
}

// Results
func (r *Repository) GetStudentResults(studentID, termID uint) ([]domain.Result, error) {
	var list []domain.Result
	q := r.db.Where("student_id = ?", studentID)
	if termID > 0 {
		q = q.Where("academic_term_id = ?", termID)
	}
	return list, q.Order("subject_id").Find(&list).Error
}
func (r *Repository) GetResult(id uint) (*domain.Result, error) {
	var res domain.Result
	return &res, r.db.First(&res, id).Error
}
func (r *Repository) CreateResult(res *domain.Result) error {
	return r.db.Create(res).Error
}
func (r *Repository) UpdateResult(res *domain.Result) error {
	return r.db.Save(res).Error
}
func (r *Repository) BulkCreateResults(results []domain.Result) error {
	return r.db.CreateInBatches(results, 100).Error
}
func (r *Repository) PublishResults(termID uint, publishedBy uint) error {
	return r.db.Exec(`UPDATE exam.results SET published_at = CURRENT_TIMESTAMP, published_by = ? WHERE academic_term_id = ? AND published_at IS NULL`, publishedBy, termID).Error
}

// Component Marks
func (r *Repository) GetComponentMarks(resultID uint) ([]domain.ComponentMarks, error) {
	var list []domain.ComponentMarks
	return list, r.db.Where("result_id = ?", resultID).Find(&list).Error
}
func (r *Repository) CreateComponentMarks(cm *domain.ComponentMarks) error {
	return r.db.Create(cm).Error
}

// Revaluation
func (r *Repository) CreateRevaluation(req *domain.RevaluationRequest) error {
	return r.db.Create(req).Error
}
func (r *Repository) GetRevaluation(id uint) (*domain.RevaluationRequest, error) {
	var req domain.RevaluationRequest
	return &req, r.db.First(&req, id).Error
}
func (r *Repository) ListRevaluations(studentID uint) ([]domain.RevaluationRequest, error) {
	var list []domain.RevaluationRequest
	q := r.db.Order("request_date DESC")
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) UpdateRevaluation(req *domain.RevaluationRequest) error {
	return r.db.Save(req).Error
}

// Supplementary
func (r *Repository) ListSupplementary(termID uint) ([]domain.SupplementaryExam, error) {
	var list []domain.SupplementaryExam
	q := r.db
	if termID > 0 {
		q = q.Where("academic_term_id = ?", termID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) CreateSupplementary(s *domain.SupplementaryExam) error {
	return r.db.Create(s).Error
}

// Academic transcript
func (r *Repository) GetTranscript(studentID uint) ([]TranscriptEntry, error) {
	var list []TranscriptEntry
	return list, r.db.Raw(`
		SELECT r.academic_term_id, t.term_name, t.academic_year,
			s.subject_code, s.subject_name, s.credits,
			r.marks_obtained, r.max_marks, r.grade, r.grade_point, r.is_passed
		FROM exam.results r
		JOIN academic.subjects s ON s.id = r.subject_id
		JOIN academic.academic_terms t ON t.id = r.academic_term_id
		WHERE r.student_id = ? AND r.published_at IS NOT NULL
		ORDER BY t.start_date, s.subject_name
	`, studentID).Scan(&list).Error
}
func (r *Repository) GetSGPA(studentID, termID uint) (*SGPAResult, error) {
	var res SGPAResult
	return &res, r.db.Raw(`
		SELECT 
			SUM(s.credits * r.grade_point) / NULLIF(SUM(s.credits), 0) as sgpa,
			SUM(s.credits) as total_credits,
			SUM(CASE WHEN r.is_passed THEN s.credits ELSE 0 END) as earned_credits
		FROM exam.results r
		JOIN academic.subjects s ON s.id = r.subject_id
		WHERE r.student_id = ? AND r.academic_term_id = ? AND r.published_at IS NOT NULL
	`, studentID, termID).Scan(&res).Error
}
func (r *Repository) GetCGPA(studentID uint) (*SGPAResult, error) {
	var res SGPAResult
	return &res, r.db.Raw(`
		SELECT 
			SUM(s.credits * r.grade_point) / NULLIF(SUM(s.credits), 0) as sgpa,
			SUM(s.credits) as total_credits,
			SUM(CASE WHEN r.is_passed THEN s.credits ELSE 0 END) as earned_credits
		FROM exam.results r
		JOIN academic.subjects s ON s.id = r.subject_id
		WHERE r.student_id = ? AND r.published_at IS NOT NULL
	`, studentID).Scan(&res).Error
}

type TranscriptEntry struct {
	AcademicTermID uint    `json:"academic_term_id"`
	TermName       string  `json:"term_name"`
	AcademicYear   string  `json:"academic_year"`
	SubjectCode    string  `json:"subject_code"`
	SubjectName    string  `json:"subject_name"`
	Credits        float32 `json:"credits"`
	MarksObtained  float64 `json:"marks_obtained"`
	MaxMarks       float64 `json:"max_marks"`
	Grade          string  `json:"grade"`
	GradePoint     float64 `json:"grade_point"`
	IsPassed       bool    `json:"is_passed"`
}

type SGPAResult struct {
	SGPA          float64 `json:"sgpa"`
	TotalCredits  float64 `json:"total_credits"`
	EarnedCredits float64 `json:"earned_credits"`
}
