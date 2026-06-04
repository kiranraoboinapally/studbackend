package admissionsmod

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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
	return s
}

func (s *Service) ListCycles(ctx context.Context) ([]domain.AdmissionCycle, error) {
	return s.repo.ListCycles()
}
func (s *Service) GetOpenCycles(ctx context.Context) ([]domain.AdmissionCycle, error) {
	return s.repo.GetOpenCycles()
}
func (s *Service) GetCycle(ctx context.Context, id uint) (*domain.AdmissionCycle, error) {
	c, err := s.repo.GetCycle(id)
	if err != nil {
		return nil, apperrors.NotFound("admission cycle not found")
	}
	return c, nil
}
func (s *Service) CreateCycle(ctx context.Context, c *domain.AdmissionCycle) error {
	if c.Name == "" {
		return apperrors.BadRequest("cycle name is required")
	}
	c.IsOpen = true
	return s.repo.CreateCycle(c)
}
func (s *Service) UpdateCycle(ctx context.Context, id uint, c *domain.AdmissionCycle) error {
        existing, err := s.repo.GetCycle(id)
        if err != nil {
                return apperrors.NotFound("admission cycle not found")
        }
        c.ID = id
        c.CreatedAt = existing.CreatedAt
        return s.repo.UpdateCycle(c)
}

// Enquiries
func (s *Service) SubmitEnquiry(ctx context.Context, e *domain.Enquiry) (*domain.Enquiry, error) {
        // Copy Program to ProgramID if Program is set but ProgramID is not
        if e.Program != nil && e.ProgramID == nil {
                e.ProgramID = e.Program
        }

        // Check if enquiry already exists with same mobile or email
        existing, _ := s.repo.GetEnquiryByMobile(e.MobileNumber)
        if existing != nil {
                // If OTP is not verified, allow updating the enquiry and resending OTP
                if !existing.MobileOTPVerified || !existing.EmailOTPVerified {
                        // Update existing enquiry with new details
                        existing.FullName = e.FullName
                        existing.Country = e.Country
                        existing.State = e.State
                        existing.District = e.District
                        existing.PreferredCampus = e.PreferredCampus
                        existing.QualificationType = e.QualificationType
                        if e.ProgramID != nil {
                                existing.ProgramID = e.ProgramID
                        }
                        // Generate separate OTPs for mobile and email
                        now := time.Now()
                        expiresAt := now.Add(15 * time.Minute)
                        
                        // Generate mobile OTP
                        if !existing.MobileOTPVerified {
                                mobileOTP := generateOTP()
                                existing.MobileOTPToken = mobileOTP
                                existing.MobileOTPSentAt = &now
                                existing.MobileOTPExpiresAt = &expiresAt
                                log.Printf("MOBILE OTP RESENT for %s: %s", existing.MobileNumber, mobileOTP)
                        }
                        
                        // Generate email OTP
                        if !existing.EmailOTPVerified {
                                emailOTP := generateOTP()
                                existing.EmailOTPToken = emailOTP
                                existing.EmailOTPSentAt = &now
                                existing.EmailOTPExpiresAt = &expiresAt
                                log.Printf("EMAIL OTP RESENT for %s: %s", existing.EmailAddress, emailOTP)
                        }
                        
                        existing.Status = "pending"
                        if err := s.repo.UpdateEnquiry(existing); err != nil {
                                log.Println("UPDATE ENQUIRY FAILED:", err)
                                return nil, err
                        }
                        return existing, nil
                }
                return nil, apperrors.Conflict("enquiry with this mobile number already exists")
        }
        existing, _ = s.repo.GetEnquiryByEmail(e.EmailAddress)
        if existing != nil {
                // If OTP is not verified, allow updating the enquiry and resending OTP
                if !existing.MobileOTPVerified || !existing.EmailOTPVerified {
                        // Update existing enquiry with new details
                        existing.FullName = e.FullName
                        existing.Country = e.Country
                        existing.State = e.State
                        existing.District = e.District
                        existing.PreferredCampus = e.PreferredCampus
                        existing.QualificationType = e.QualificationType
                        if e.ProgramID != nil {
                                existing.ProgramID = e.ProgramID
                        }
                        // Generate separate OTPs for mobile and email
                        now := time.Now()
                        expiresAt := now.Add(15 * time.Minute)
                        
                        // Generate mobile OTP
                        if !existing.MobileOTPVerified {
                                mobileOTP := generateOTP()
                                existing.MobileOTPToken = mobileOTP
                                existing.MobileOTPSentAt = &now
                                existing.MobileOTPExpiresAt = &expiresAt
                                log.Printf("MOBILE OTP RESENT for %s: %s", existing.MobileNumber, mobileOTP)
                        }
                        
                        // Generate email OTP
                        if !existing.EmailOTPVerified {
                                emailOTP := generateOTP()
                                existing.EmailOTPToken = emailOTP
                                existing.EmailOTPSentAt = &now
                                existing.EmailOTPExpiresAt = &expiresAt
                                log.Printf("EMAIL OTP RESENT for %s: %s", existing.EmailAddress, emailOTP)
                        }
                        
                        existing.Status = "pending"
                        if err := s.repo.UpdateEnquiry(existing); err != nil {
                                log.Println("UPDATE ENQUIRY FAILED:", err)
                                return nil, err
                        }
                        return existing, nil
                }
                return nil, apperrors.Conflict("enquiry with this email already exists")
        }

        // Generate separate OTPs for mobile and email for new enquiry
        now := time.Now()
        expiresAt := now.Add(15 * time.Minute)
        
        // Generate mobile OTP
        mobileOTP := generateOTP()
        e.MobileOTPToken = mobileOTP
        e.MobileOTPSentAt = &now
        e.MobileOTPExpiresAt = &expiresAt
        e.MobileOTPVerified = false
        log.Printf("MOBILE OTP GENERATED for %s: %s", e.MobileNumber, mobileOTP)
        
        // Generate email OTP
        emailOTP := generateOTP()
        e.EmailOTPToken = emailOTP
        e.EmailOTPSentAt = &now
        e.EmailOTPExpiresAt = &expiresAt
        e.EmailOTPVerified = false
        log.Printf("EMAIL OTP GENERATED for %s: %s", e.EmailAddress, emailOTP)
        
        e.Status = "pending"
        e.OTPVerified = false

        if err := s.repo.CreateEnquiry(e); err != nil {
                log.Println("CREATE ENQUIRY FAILED:", err)
                return nil, err
        }
        return e, nil
}

func generateOTP() string {
        return fmt.Sprintf("%06d", 100000+rand.Intn(900000))
}

func (s *Service) ResendOTP(ctx context.Context, id uint) (*domain.Enquiry, error) {
        e, err := s.repo.GetEnquiryByID(id)
        if err != nil {
                return nil, apperrors.NotFound("enquiry not found")
        }
        if e.MobileOTPVerified && e.EmailOTPVerified {
                return nil, apperrors.BadRequest("Both OTPs already verified")
        }
        
        // Generate separate OTPs for mobile and email
        now := time.Now()
        expiresAt := now.Add(15 * time.Minute)
        
        // Generate mobile OTP if not verified
        if !e.MobileOTPVerified {
                mobileOTP := generateOTP()
                e.MobileOTPToken = mobileOTP
                e.MobileOTPSentAt = &now
                e.MobileOTPExpiresAt = &expiresAt
                log.Printf("MOBILE OTP RESENT for %s: %s", e.MobileNumber, mobileOTP)
        }
        
        // Generate email OTP if not verified
        if !e.EmailOTPVerified {
                emailOTP := generateOTP()
                e.EmailOTPToken = emailOTP
                e.EmailOTPSentAt = &now
                e.EmailOTPExpiresAt = &expiresAt
                log.Printf("EMAIL OTP RESENT for %s: %s", e.EmailAddress, emailOTP)
        }
        
        e.Status = "pending"
        if err := s.repo.UpdateEnquiry(e); err != nil {
                log.Println("UPDATE ENQUIRY FAILED:", err)
                return nil, err
        }
        return e, nil
}

func (s *Service) GetEnquiryByID(ctx context.Context, id uint) (*domain.Enquiry, error) {
        e, err := s.repo.GetEnquiryByID(id)
        if err != nil {
                return nil, apperrors.NotFound("enquiry not found")
        }
        return e, nil
}

func (s *Service) GetEnquiryByMobile(ctx context.Context, mobileNumber string) (*domain.Enquiry, error) {
        e, err := s.repo.GetEnquiryByMobile(mobileNumber)
        if err != nil {
                return nil, apperrors.NotFound("enquiry not found")
        }
        return e, nil
}

func (s *Service) GetEnquiryByEmail(ctx context.Context, emailAddress string) (*domain.Enquiry, error) {
        e, err := s.repo.GetEnquiryByEmail(emailAddress)
        if err != nil {
                return nil, apperrors.NotFound("enquiry not found")
        }
        return e, nil
}

func (s *Service) GetLatestUnverifiedEnquiry(ctx context.Context) (*domain.Enquiry, error) {
        e, err := s.repo.GetLatestUnverifiedEnquiry()
        if err != nil {
                return nil, apperrors.NotFound("no unverified enquiry found")
        }
        return e, nil
}
func (s *Service) VerifyOTP(ctx context.Context, id uint, otp string, otpType string) error {
        e, err := s.repo.GetEnquiryByID(id)
        if err != nil {
                return apperrors.NotFound("enquiry not found")
        }
        
        var storedOTP string
        var expiresAt *time.Time
        
        // Determine which OTP to verify based on type
        if otpType == "mobile" {
                storedOTP = e.MobileOTPToken
                expiresAt = e.MobileOTPExpiresAt
        } else if otpType == "email" {
                storedOTP = e.EmailOTPToken
                expiresAt = e.EmailOTPExpiresAt
        } else {
                // For backward compatibility, try mobile OTP first
                storedOTP = e.MobileOTPToken
                expiresAt = e.MobileOTPExpiresAt
        }
        
        if storedOTP != otp {
                return apperrors.BadRequest("invalid OTP")
        }
        if expiresAt != nil && time.Now().After(*expiresAt) {
                return apperrors.BadRequest("OTP expired")
        }
        
        // Mark the specific OTP as verified
        if otpType == "mobile" {
                e.MobileOTPVerified = true
        } else if otpType == "email" {
                e.EmailOTPVerified = true
        }
        
        // Update overall status if both are verified
        if e.MobileOTPVerified && e.EmailOTPVerified {
                e.OTPVerified = true
                e.Status = "verified"
        }
        
        return s.repo.UpdateEnquiry(e)
}
func (s *Service) ListEnquiries(ctx context.Context, page, pageSize int) ([]domain.Enquiry, int64, error) {
        return s.repo.ListEnquiries(page, pageSize)
}
func (s *Service) CloseCycle(ctx context.Context, id uint) error {
	existing, err := s.repo.GetCycle(id)
	if err != nil {
		return apperrors.NotFound("cycle not found")
	}
	existing.IsOpen = false
	return s.repo.UpdateCycle(existing)
}

func (s *Service) ListApplicants(ctx context.Context, cycleID uint, page, pageSize int) ([]domain.Applicant, int64, error) {
	return s.repo.ListApplicants(cycleID, page, pageSize)
}
func (s *Service) GetApplicant(ctx context.Context, id uint) (*domain.Applicant, error) {
	a, err := s.repo.GetApplicant(id)
	if err != nil {
		return nil, apperrors.NotFound("applicant not found")
	}
	return a, nil
}
func (s *Service) Submit(ctx context.Context, a *domain.Applicant) error {
	if a.FirstName == "" || a.Email == "" {
		return apperrors.BadRequest("first name and email are required")
	}
	count, _ := s.repo.CountApplicationNumber(a.CycleID)
	a.ApplicationNumber = fmt.Sprintf("APP-%d-%04d", time.Now().Year(), count+1)
	a.AppliedAt = time.Now()

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(a).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "Applicant", fmt.Sprintf("%d", a.ID),
			eventbus.EventApplicationSubmitted,
			map[string]interface{}{
				"applicant_id":      a.ID,
				"application_number": a.ApplicationNumber,
				"email":             a.Email,
				"cycle_id":          a.CycleID,
			},
		)
	})
	return txErr
}
func (s *Service) UpdateApplicant(ctx context.Context, id uint, a *domain.Applicant) error {
	existing, err := s.repo.GetApplicant(id)
	if err != nil {
		return apperrors.NotFound("applicant not found")
	}
	a.ID = existing.ID
	a.ApplicationNumber = existing.ApplicationNumber
	return s.repo.UpdateApplicant(a)
}

// UpdateStatus changes the applicant status and, if approved, publishes an outbox event
// that triggers automatic student profile creation via the event bus.
func (s *Service) UpdateStatus(ctx context.Context, id uint, statusID uint, changedBy uint) error {
	applicant, err := s.repo.GetApplicant(id)
	if err != nil {
		return apperrors.NotFound("applicant not found")
	}

	// Find the "APPROVED" status code
	var approvedStatus domain.StatusCode
	s.db.Where("module = ? AND code = ?", "admission", "APPROVED").First(&approvedStatus)

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Applicant{}).Where("id = ?", id).Update("status_id", statusID).Error; err != nil {
			return err
		}
		hist := &domain.ApplicationStatusHistory{
			ApplicantID:   id,
			StatusID:      statusID,
			EffectiveFrom: time.Now(),
		}
		if err := tx.Create(hist).Error; err != nil {
			return err
		}

		// If approved, emit ApplicationApproved event for auto student creation
		if approvedStatus.ID != 0 && statusID == approvedStatus.ID {
			programID := uint(0)
			if applicant.ProgramID != nil {
				programID = *applicant.ProgramID
			}
			return s.outbox.WriteEvent(tx, "Applicant", fmt.Sprintf("%d", id),
				eventbus.EventApplicationApproved,
				eventbus.ApplicationApprovedPayload{
					ApplicantID: id,
					ProgramID:   programID,
					Email:       applicant.Email,
					FirstName:   applicant.FirstName,
					LastName:    applicant.LastName,
					DateOfBirth: applicant.DateOfBirth.Format("2006-01-02"),
					GenderID:    applicant.GenderID,
					CategoryID:  applicant.CategoryID,
					Phone:       applicant.Phone,
					CycleID:     applicant.CycleID,
				},
			)
		}
		return nil
	})
	return txErr
}

func (s *Service) GetStatusHistory(ctx context.Context, id uint) ([]domain.ApplicationStatusHistory, error) {
	return s.repo.GetStatusHistory(id)
}

// Documents
func (s *Service) GetDocuments(ctx context.Context, applicantID uint) ([]domain.Document, error) {
	return s.repo.GetApplicantDocuments(applicantID)
}
func (s *Service) UploadDocument(ctx context.Context, d *domain.Document) error {
	if d.ApplicantID == 0 || d.DocumentType == "" {
		return apperrors.BadRequest("applicant_id and document_type are required")
	}
	d.UploadedAt = time.Now()
	return s.repo.CreateDocument(d)
}
func (s *Service) VerifyDocument(ctx context.Context, docID, verifiedBy uint) error {
	return s.repo.VerifyDocument(docID, verifiedBy)
}

// Seat Allocation
func (s *Service) AllocateSeat(ctx context.Context, sa *domain.SeatAllocation) error {
	if sa.ApplicantID == 0 || sa.CycleID == 0 {
		return apperrors.BadRequest("applicant_id and cycle_id are required")
	}
	sa.AllocatedAt = time.Now()

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(sa).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "SeatAllocation", fmt.Sprintf("%d", sa.ID),
			eventbus.EventSeatAllocated,
			map[string]interface{}{
				"allocation_id": sa.ID,
				"applicant_id": sa.ApplicantID,
				"cycle_id":     sa.CycleID,
				"rank":         sa.AllocationRank,
			},
		)
	})
	return txErr
}
func (s *Service) GetSeatAllocation(ctx context.Context, applicantID uint) (*domain.SeatAllocation, error) {
	sa, err := s.repo.GetSeatAllocation(applicantID)
	if err != nil {
		return nil, apperrors.NotFound("seat allocation not found")
	}
	return sa, nil
}
func (s *Service) ListSeatAllocations(ctx context.Context, cycleID uint) ([]domain.SeatAllocation, error) {
	return s.repo.ListSeatAllocations(cycleID)
}

// Waitlist
func (s *Service) AddToWaitlist(ctx context.Context, w *domain.Waitlist) error {
	return s.repo.AddToWaitlist(w)
}
func (s *Service) GetWaitlist(ctx context.Context, cycleID uint) ([]domain.Waitlist, error) {
	return s.repo.GetWaitlist(cycleID)
}

// Conversion - admission to student
func (s *Service) ConvertToStudent(ctx context.Context, applicantID, studentID uint) error {
	m := &domain.ApplicantStudentMap{
		ApplicantID: applicantID,
		StudentID:   studentID,
		MappedAt:    time.Now(),
	}
	return s.repo.CreateApplicantStudentMap(m)
}

func (s *Service) GetCycleStats(ctx context.Context, cycleID uint) (map[string]int64, error) {
	return s.repo.GetCycleStats(cycleID)
}
