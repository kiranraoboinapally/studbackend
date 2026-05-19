package exammod

import (
	"context"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
)

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

// Components
func (s *Service) ListComponents(ctx context.Context, subjectID uint) ([]domain.ExamComponent, error) {
	return s.repo.ListComponents(subjectID)
}
func (s *Service) CreateComponent(ctx context.Context, c *domain.ExamComponent) error {
	return s.repo.CreateComponent(c)
}
func (s *Service) UpdateComponent(ctx context.Context, id uint, c *domain.ExamComponent) error {
	existing, err := s.repo.GetComponent(id)
	if err != nil {
		return apperrors.NotFound("component not found")
	}
	c.ID = existing.ID
	return s.repo.UpdateComponent(c)
}

// Schedules
func (s *Service) ListSchedules(ctx context.Context, termID, subjectID uint) ([]domain.ExamSchedule, error) {
	return s.repo.ListSchedules(termID, subjectID)
}
func (s *Service) GetSchedule(ctx context.Context, id uint) (*domain.ExamSchedule, error) {
	sch, err := s.repo.GetSchedule(id)
	if err != nil {
		return nil, apperrors.NotFound("exam schedule not found")
	}
	return sch, nil
}
func (s *Service) CreateSchedule(ctx context.Context, sch *domain.ExamSchedule) error {
	if sch.SubjectID == 0 || sch.AcademicTermID == 0 {
		return apperrors.BadRequest("subject_id and academic_term_id are required")
	}
	return s.repo.CreateSchedule(sch)
}
func (s *Service) UpdateSchedule(ctx context.Context, id uint, sch *domain.ExamSchedule) error {
	existing, err := s.repo.GetSchedule(id)
	if err != nil {
		return apperrors.NotFound("schedule not found")
	}
	sch.ID = existing.ID
	return s.repo.UpdateSchedule(sch)
}

// Results
func (s *Service) GetStudentResults(ctx context.Context, studentID, termID uint) ([]domain.Result, error) {
	return s.repo.GetStudentResults(studentID, termID)
}
func (s *Service) EnterResult(ctx context.Context, res *domain.Result) error {
	if res.StudentID == 0 || res.SubjectID == 0 {
		return apperrors.BadRequest("student_id and subject_id are required")
	}
	// Calculate grade based on marks
	percentage := (res.MarksObtained / res.MaxMarks) * 100
	res.Grade, res.GradePoint, res.IsPassed = calculateGrade(percentage)
	return s.repo.CreateResult(res)
}
func (s *Service) UpdateResult(ctx context.Context, id uint, res *domain.Result) error {
	existing, err := s.repo.GetResult(id)
	if err != nil {
		return apperrors.NotFound("result not found")
	}
	res.ID = existing.ID
	percentage := (res.MarksObtained / res.MaxMarks) * 100
	res.Grade, res.GradePoint, res.IsPassed = calculateGrade(percentage)
	return s.repo.UpdateResult(res)
}
func (s *Service) BulkEnterResults(ctx context.Context, results []domain.Result) error {
	for i := range results {
		if results[i].MaxMarks > 0 {
			percentage := (results[i].MarksObtained / results[i].MaxMarks) * 100
			results[i].Grade, results[i].GradePoint, results[i].IsPassed = calculateGrade(percentage)
		}
	}
	return s.repo.BulkCreateResults(results)
}
func (s *Service) PublishResults(ctx context.Context, termID, publishedBy uint) error {
	return s.repo.PublishResults(termID, publishedBy)
}

// Component Marks
func (s *Service) EnterComponentMarks(ctx context.Context, cm *domain.ComponentMarks) error {
	return s.repo.CreateComponentMarks(cm)
}
func (s *Service) GetComponentMarks(ctx context.Context, resultID uint) ([]domain.ComponentMarks, error) {
	return s.repo.GetComponentMarks(resultID)
}

// Revaluation
func (s *Service) RequestRevaluation(ctx context.Context, req *domain.RevaluationRequest) error {
	if req.ResultID == 0 || req.StudentID == 0 {
		return apperrors.BadRequest("result_id and student_id are required")
	}
	req.RequestDate = time.Now()
	return s.repo.CreateRevaluation(req)
}
func (s *Service) ListRevaluations(ctx context.Context, studentID uint) ([]domain.RevaluationRequest, error) {
	return s.repo.ListRevaluations(studentID)
}
func (s *Service) ProcessRevaluation(ctx context.Context, id uint, reviewedMarks float64, grade, remarks string, processedBy uint) error {
	req, err := s.repo.GetRevaluation(id)
	if err != nil {
		return apperrors.NotFound("revaluation request not found")
	}
	now := time.Now()
	req.ReviewedMarks = reviewedMarks
	req.ReviewedGrade = grade
	req.Remarks = remarks
	req.ProcessedAt = &now
	return s.repo.UpdateRevaluation(req)
}

// Supplementary
func (s *Service) ListSupplementary(ctx context.Context, termID uint) ([]domain.SupplementaryExam, error) {
	return s.repo.ListSupplementary(termID)
}
func (s *Service) CreateSupplementary(ctx context.Context, sup *domain.SupplementaryExam) error {
	return s.repo.CreateSupplementary(sup)
}

// Transcript & GPA
func (s *Service) GetTranscript(ctx context.Context, studentID uint) ([]TranscriptEntry, error) {
	return s.repo.GetTranscript(studentID)
}
func (s *Service) GetSGPA(ctx context.Context, studentID, termID uint) (*SGPAResult, error) {
	return s.repo.GetSGPA(studentID, termID)
}
func (s *Service) GetCGPA(ctx context.Context, studentID uint) (*SGPAResult, error) {
	return s.repo.GetCGPA(studentID)
}

// Grade calculation helper
func calculateGrade(percentage float64) (string, float64, bool) {
	switch {
	case percentage >= 90:
		return "O", 10.0, true
	case percentage >= 80:
		return "A+", 9.0, true
	case percentage >= 70:
		return "A", 8.0, true
	case percentage >= 60:
		return "B+", 7.0, true
	case percentage >= 55:
		return "B", 6.0, true
	case percentage >= 50:
		return "C", 5.0, true
	case percentage >= 40:
		return "P", 4.0, true
	default:
		return "F", 0.0, false
	}
}
