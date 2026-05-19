package domain

import "time"

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
