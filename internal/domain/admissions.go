package domain

import (
	"encoding/json"
	"strconv"
	"time"
)

// ─── Admissions ──────────────────────────────────────────────────────────────

type AdmissionCycle struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `gorm:"not null" json:"name"`
	AcademicYear     string     `gorm:"not null;index" json:"academic_year"`
	ProgramID        *uint      `gorm:"index" json:"program_id,omitempty"`
	ApplicationStart time.Time  `gorm:"not null" json:"application_start"`
	ApplicationEnd   time.Time  `gorm:"not null" json:"application_end"`
	EntranceExamDate *time.Time `json:"entrance_exam_date,omitempty"`
	CounselingStart  *time.Time `json:"counseling_start,omitempty"`
	CounselingEnd    *time.Time `json:"counseling_end,omitempty"`
	ApplicationFee   float64    `json:"application_fee"`
	MaxApplications  int        `json:"max_applications"`
	IsOpen           bool       `gorm:"default:true;index" json:"is_open"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (AdmissionCycle) TableName() string { return "admissions.admission_cycles" }

type Applicant struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ApplicationNumber string    `gorm:"unique;not null;index" json:"application_number"`
	CycleID           uint      `gorm:"not null;index" json:"cycle_id"`
	ProgramID         *uint     `gorm:"index" json:"program_id,omitempty"`
	FirstName         string    `gorm:"not null" json:"first_name"`
	LastName          string    `gorm:"not null" json:"last_name"`
	DateOfBirth       time.Time `gorm:"not null" json:"date_of_birth"`
	Email             string    `gorm:"not null;index" json:"email"`
	Phone             string    `json:"phone"`
	Address           string    `json:"address"`
	GenderID          *uint     `gorm:"index" json:"gender_id,omitempty"`
	CategoryID        *uint     `gorm:"index" json:"category_id,omitempty"`
	EntranceScore     float64   `json:"entrance_score"`
	Rank              int       `json:"rank"`
	StatusID          *uint     `gorm:"index" json:"status_id,omitempty"`
	AppliedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"applied_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (Applicant) TableName() string { return "admissions.applicants" }

type ApplicationStatusHistory struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ApplicantID   uint       `gorm:"not null;index" json:"applicant_id"`
	StatusID      uint       `gorm:"not null" json:"status_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (ApplicationStatusHistory) TableName() string { return "admissions.application_status_history" }

type Document struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	ApplicantID          uint       `gorm:"not null;index" json:"applicant_id"`
	DocumentType         string     `gorm:"not null" json:"document_type"`
	FilePath             string     `gorm:"not null" json:"file_path"`
	UploadedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"uploaded_at"`
	VerifiedBy           *uint      `json:"verified_by,omitempty"`
	VerifiedAt           *time.Time `json:"verified_at,omitempty"`
	VerificationStatusID *uint      `gorm:"index" json:"verification_status_id,omitempty"`
}

func (Document) TableName() string { return "admissions.documents" }

type SeatAllocation struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ApplicantID    uint      `gorm:"not null;index" json:"applicant_id"`
	CycleID        uint      `gorm:"not null;index" json:"cycle_id"`
	AllocationRank int       `json:"allocation_rank"`
	StatusID       *uint     `gorm:"index" json:"status_id,omitempty"`
	AllocatedAt    time.Time `json:"allocated_at"`
	CreatedAt      time.Time `json:"created_at"`
}

func (SeatAllocation) TableName() string { return "admissions.seat_allocations" }

type ApplicantStudentMap struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ApplicantID uint      `gorm:"unique;not null;index" json:"applicant_id"`
	StudentID   uint      `gorm:"unique;not null;index" json:"student_id"`
	MappedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"mapped_at"`
}

func (ApplicantStudentMap) TableName() string { return "admissions.applicant_student_map" }

type Waitlist struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ApplicantID uint       `gorm:"not null;index" json:"applicant_id"`
	CycleID     uint       `gorm:"not null;index" json:"cycle_id"`
	Rank        int        `gorm:"not null;index" json:"rank"`
	StatusID    *uint      `gorm:"index" json:"status_id,omitempty"`
	NotifiedAt  *time.Time `json:"notified_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (Waitlist) TableName() string { return "admissions.waitlist" }

type Enquiry struct {
    ID               uint       `gorm:"primaryKey" json:"id"`
    FullName         string     `gorm:"not null" json:"full_name"`
    MobileNumber     string     `gorm:"unique;not null" json:"mobile_number"`
    EmailAddress     string     `gorm:"unique;not null" json:"email_address"`
    Country          string     `json:"country,omitempty"`
    State            string     `json:"state,omitempty"`
    District         string     `json:"district,omitempty"`
    PreferredCampus  *int64     `json:"preferred_campus,omitempty"`
    QualificationType string    `json:"qualification_type,omitempty"`
    ProgramID        *int64     `json:"program_id,omitempty"`
    Program          *int64     `gorm:"-" json:"program,omitempty"` // Alias for ProgramID to accept both field names, ignored by GORM
    Status           string     `gorm:"default:pending" json:"status"`
    OTPVerified      bool       `gorm:"default:false" json:"otp_verified"`
    OTPToken         string     `json:"otp_token,omitempty"`
    OTPSentAt        *time.Time `json:"otp_sent_at,omitempty"`
    OTPExpiresAt     *time.Time `json:"otp_expires_at,omitempty"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}

func (Enquiry) TableName() string { return "admissions.enquiries" }

// UnmarshalJSON handles both string and int64 values for numeric fields
func (e *Enquiry) UnmarshalJSON(data []byte) error {
	type Alias Enquiry
	aux := &struct {
		PreferredCampus interface{} `json:"preferred_campus,omitempty"`
		Program         interface{} `json:"program,omitempty"`
		ProgramID       interface{} `json:"program_id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle PreferredCampus
	if aux.PreferredCampus != nil {
		switch v := aux.PreferredCampus.(type) {
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				e.PreferredCampus = &i
			}
		case float64:
			i := int64(v)
			e.PreferredCampus = &i
		}
	}

	// Handle Program
	if aux.Program != nil {
		switch v := aux.Program.(type) {
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				e.Program = &i
			}
		case float64:
			i := int64(v)
			e.Program = &i
		}
	}

	// Handle ProgramID
	if aux.ProgramID != nil {
		switch v := aux.ProgramID.(type) {
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				e.ProgramID = &i
			}
		case float64:
			i := int64(v)
			e.ProgramID = &i
		}
	}

	return nil
}
