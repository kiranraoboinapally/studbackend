package studentmod

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

// Service contains student business logic.
type Service struct {
	repo   *Repository
	bus    *eventbus.Bus
	outbox *outbox.Writer
	db     *gorm.DB
}

func NewService(repo *Repository, bus *eventbus.Bus, ob *outbox.Writer, db *gorm.DB) *Service {
	s := &Service{repo: repo, bus: bus, outbox: ob, db: db}
	// Subscribe to admission approval events for auto student creation
	bus.Subscribe(eventbus.EventApplicationApproved, s.HandleApplicationApproved)
	return s
}

// HandleApplicationApproved auto-creates a student profile when an application is approved.
func (s *Service) HandleApplicationApproved(ctx context.Context, evt eventbus.Event) error {
	slog.Info("StudentMod: received ApplicationApproved event", "aggregate_id", evt.AggregateID)

	payloadBytes, _ := evt.Payload.(string)
	var payload eventbus.ApplicationApprovedPayload
	if err := json.Unmarshal([]byte(payloadBytes), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	dob, err := time.Parse("2006-01-02", payload.DateOfBirth)
	if err != nil {
		return fmt.Errorf("invalid date_of_birth: %w", err)
	}

	year := time.Now().Year()
	var student *domain.Student

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&domain.Student{}).Where("program_id = ? AND admission_year = ?", payload.ProgramID, year).Count(&count)
		rollNo := fmt.Sprintf("%dP%dS%03d", year, payload.ProgramID, count+1)
		enrollNo := fmt.Sprintf("NTU%d%s", year, rollNo)

		var activeStatus domain.StatusCode
		tx.Where("module = ? AND code = ?", "student", "ACTIVE").First(&activeStatus)

		student = &domain.Student{
			EnrollmentNumber: enrollNo,
			RollNumber:       rollNo,
			FirstName:        payload.FirstName,
			LastName:         payload.LastName,
			DateOfBirth:      dob,
			GenderID:         payload.GenderID,
			Phone:            payload.Phone,
			Email:            payload.Email,
			CategoryID:       payload.CategoryID,
			ProgramID:        payload.ProgramID,
			AdmissionYear:    year,
			StatusID:         &activeStatus.ID,
			AcademicStanding: "Good",
		}

		if err := tx.Create(student).Error; err != nil {
			return err
		}

		tx.Create(&domain.StudentStatusHistory{
			StudentID:     student.ID,
			StatusID:      activeStatus.ID,
			EffectiveFrom: time.Now(),
			Reason:        "Auto-created from admission approval",
		})

		// Create applicant-student mapping
		tx.Create(&domain.ApplicantStudentMap{
			ApplicantID: payload.ApplicantID,
			StudentID:   student.ID,
			MappedAt:    time.Now(),
		})

		// Get current academic term
		var term domain.AcademicTerm
		tx.Where("is_current = ?", true).First(&term)

		return s.outbox.WriteEvent(tx, "Student", fmt.Sprintf("%d", student.ID),
			eventbus.EventStudentEnrolled,
			eventbus.StudentEnrolledPayload{
				StudentID:  student.ID,
				UserID:     student.UserID,
				ProgramID:  student.ProgramID,
				TermID:     term.ID,
				RollNumber: student.RollNumber,
			},
		)
	})

	if txErr != nil {
		slog.Error("StudentMod: failed to auto-create student from admission", "error", txErr)
		return txErr
	}

	slog.Info("StudentMod: auto-created student from admission", "student_id", student.ID, "applicant_id", payload.ApplicantID)
	return nil
}

// ─── Enroll Student ──────────────────────────────────────────────────────────

type EnrollRequest struct {
	UserID          uint   `json:"user_id"`
	ProgramID       uint   `json:"program_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	DateOfBirth     string `json:"date_of_birth"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	GenderID        *uint  `json:"gender_id"`
	CategoryID      *uint  `json:"category_id"`
	AdmissionYear   int    `json:"admission_year"`
	AdmissionQuota  string `json:"admission_quota"`
}

func (s *Service) EnrollStudent(ctx context.Context, req EnrollRequest) (*domain.Student, error) {
	if req.FirstName == "" || req.Email == "" || req.ProgramID == 0 {
		return nil, apperrors.BadRequest("first_name, email, and program_id are required")
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, apperrors.Validation("date_of_birth must be YYYY-MM-DD format")
	}

	year := req.AdmissionYear
	if year == 0 {
		year = time.Now().Year()
	}

	var student *domain.Student

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&domain.Student{}).Where("program_id = ? AND admission_year = ?", req.ProgramID, year).Count(&count)
		rollNo := fmt.Sprintf("%dP%dS%03d", year, req.ProgramID, count+1)
		enrollNo := fmt.Sprintf("NTU%d%s", year, rollNo)

		var activeStatus domain.StatusCode
		tx.Where("module = ? AND code = ?", "student", "ACTIVE").First(&activeStatus)

		student = &domain.Student{
			UserID:           req.UserID,
			EnrollmentNumber: enrollNo,
			RollNumber:       rollNo,
			FirstName:        req.FirstName,
			LastName:         req.LastName,
			DateOfBirth:      dob,
			GenderID:         req.GenderID,
			Phone:            req.Phone,
			Email:            req.Email,
			CategoryID:       req.CategoryID,
			ProgramID:        req.ProgramID,
			AdmissionYear:    year,
			AdmissionQuota:   req.AdmissionQuota,
			StatusID:         &activeStatus.ID,
			AcademicStanding: "Good",
		}

		if err := tx.Create(student).Error; err != nil {
			return err
		}

		tx.Create(&domain.StudentStatusHistory{
			StudentID:     student.ID,
			StatusID:      activeStatus.ID,
			EffectiveFrom: time.Now(),
			Reason:        "Initial enrollment",
		})

		var term domain.AcademicTerm
		tx.Where("is_current = ?", true).First(&term)

		return s.outbox.WriteEvent(tx, "Student", fmt.Sprintf("%d", student.ID),
			eventbus.EventStudentEnrolled,
			eventbus.StudentEnrolledPayload{
				StudentID:  student.ID,
				UserID:     student.UserID,
				ProgramID:  student.ProgramID,
				TermID:     term.ID,
				RollNumber: student.RollNumber,
			},
		)
	})

	if txErr != nil {
		return nil, apperrors.Internal("enrollment failed", txErr)
	}
	return student, nil
}

// ─── Get Student ─────────────────────────────────────────────────────────────

func (s *Service) GetByID(ctx context.Context, id uint) (*domain.Student, error) {
	student, err := s.repo.FindByID(id)
	if err != nil {
		return nil, apperrors.NotFound("student")
	}
	return student, nil
}

func (s *Service) GetByUserID(ctx context.Context, userID uint) (*domain.Student, error) {
	student, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, apperrors.NotFound("student")
	}
	return student, nil
}

// ─── List Students ───────────────────────────────────────────────────────────

func (s *Service) List(ctx context.Context, page, pageSize int) ([]domain.Student, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListAll(page, pageSize)
}

// ─── Dashboard ───────────────────────────────────────────────────────────────

func (s *Service) GetDashboard(ctx context.Context, studentID uint) (map[string]interface{}, error) {
	return s.repo.GetDashboard(studentID)
}

// ─── Guardians ───────────────────────────────────────────────────────────────

func (s *Service) AddGuardian(ctx context.Context, g *domain.Guardian) error {
	return s.repo.CreateGuardian(g)
}

func (s *Service) GetGuardians(ctx context.Context, studentID uint) ([]domain.Guardian, error) {
	return s.repo.GetGuardians(studentID)
}

// ─── Grievances ──────────────────────────────────────────────────────────────

func (s *Service) FileGrievance(ctx context.Context, g *domain.Grievance) error {
	return s.repo.CreateGrievance(g)
}

func (s *Service) GetGrievances(ctx context.Context, studentID uint) ([]domain.Grievance, error) {
	return s.repo.GetGrievances(studentID)
}

// ─── Update Student ──────────────────────────────────────────────────────────

func (s *Service) UpdateStudent(ctx context.Context, id uint, student *domain.Student) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return apperrors.NotFound("student")
	}
	student.ID = existing.ID
	student.EnrollmentNumber = existing.EnrollmentNumber
	student.RollNumber = existing.RollNumber
	return s.repo.Update(student)
}

// ─── Medical Records ─────────────────────────────────────────────────────────

func (s *Service) GetMedicalRecord(ctx context.Context, studentID uint) (*domain.MedicalRecord, error) {
	rec, err := s.repo.GetMedicalRecord(studentID)
	if err != nil {
		return nil, apperrors.NotFound("medical record")
	}
	return rec, nil
}

func (s *Service) CreateMedicalRecord(ctx context.Context, m *domain.MedicalRecord) error {
	return s.repo.CreateMedicalRecord(m)
}

// ─── Alumni ──────────────────────────────────────────────────────────────────

func (s *Service) GetAlumni(ctx context.Context, studentID uint) (*domain.Alumni, error) {
	var a domain.Alumni
	err := s.db.Where("student_id = ?", studentID).First(&a).Error
	if err != nil {
		return nil, apperrors.NotFound("alumni record")
	}
	return &a, nil
}

func (s *Service) CreateAlumni(ctx context.Context, a *domain.Alumni) error {
	return s.db.Create(a).Error
}
