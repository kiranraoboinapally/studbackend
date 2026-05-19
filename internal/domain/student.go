package domain

import "time"

// ─── Student Profile ─────────────────────────────────────────────────────────

type Student struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uint      `gorm:"unique;not null;index" json:"user_id"`
	EnrollmentNumber    string    `gorm:"unique;not null;index" json:"enrollment_number"`
	RollNumber          string    `gorm:"unique;index" json:"roll_number"`
	FirstName           string    `gorm:"not null" json:"first_name"`
	LastName            string    `gorm:"not null" json:"last_name"`
	DateOfBirth         time.Time `gorm:"not null" json:"date_of_birth"`
	GenderID            *uint     `gorm:"index" json:"gender_id,omitempty"`
	Phone               string    `json:"phone"`
	Email               string    `gorm:"not null;index" json:"email"`
	AlternateEmail      string    `json:"alternate_email"`
	Address             string    `json:"address"`
	City                string    `json:"city"`
	State               string    `json:"state"`
	PostalCode          string    `json:"postal_code"`
	Nationality         string    `gorm:"default:'Indian'" json:"nationality"`
	CategoryID          *uint     `gorm:"index" json:"category_id,omitempty"`
	ProgramID           uint      `gorm:"not null;index" json:"program_id"`
	AdmissionYear       int       `gorm:"not null;index" json:"admission_year"`
	AdmissionQuota      string    `json:"admission_quota"`
	IsHostelRequired    bool      `gorm:"default:false" json:"is_hostel_required"`
	IsTransportRequired bool      `gorm:"default:false" json:"is_transport_required"`
	StatusID            *uint     `gorm:"index" json:"status_id,omitempty"`
	AcademicStanding    string    `gorm:"type:varchar(30);default:'Good'" json:"academic_standing"`
	ProfilePhoto        string    `json:"profile_photo"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (Student) TableName() string { return "student.students" }

type StudentStatusHistory struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	StudentID     uint       `gorm:"not null;index" json:"student_id"`
	StatusID      uint       `gorm:"not null" json:"status_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Reason        string     `json:"reason"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (StudentStatusHistory) TableName() string { return "student.student_status_history" }

type Guardian struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	StudentID  uint   `gorm:"not null;index" json:"student_id"`
	Name       string `gorm:"not null" json:"name"`
	Relation   string `json:"relation"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Occupation string `json:"occupation"`
	IsPrimary  bool   `gorm:"default:false" json:"is_primary"`
}

func (Guardian) TableName() string { return "student.guardians" }

type MedicalRecord struct {
	ID                    uint       `gorm:"primaryKey" json:"id"`
	StudentID             uint       `gorm:"unique;not null" json:"student_id"`
	BloodGroupID          *uint      `gorm:"index" json:"blood_group_id,omitempty"`
	Allergies             string     `json:"allergies"`
	ChronicConditions     string     `json:"chronic_conditions"`
	EmergencyContactName  string     `json:"emergency_contact_name"`
	EmergencyContactPhone string     `json:"emergency_contact_phone"`
	InsurancePolicyNo     string     `json:"insurance_policy_no"`
	ValidUntil            *time.Time `json:"valid_until,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (MedicalRecord) TableName() string { return "student.medical_records" }

type Grievance struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	StudentID   uint       `gorm:"not null;index" json:"student_id"`
	Category    string     `gorm:"index" json:"category"`
	Description string     `gorm:"not null" json:"description"`
	StatusID    *uint      `gorm:"index" json:"status_id,omitempty"`
	AssignedTo  *uint      `json:"assigned_to,omitempty"`
	Resolution  string     `json:"resolution"`
	CreatedAt   time.Time  `gorm:"index" json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

func (Grievance) TableName() string { return "student.grievances" }

// ─── Student Attendance ──────────────────────────────────────────────────────

type ClassSession struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	OfferingID uint       `gorm:"not null;index" json:"offering_id"`
	ClassDate  time.Time  `gorm:"not null;index" json:"class_date"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	RoomID     *uint      `gorm:"index" json:"room_id,omitempty"`
	FacultyID  uint       `gorm:"not null" json:"faculty_id"`
	StatusID   *uint      `gorm:"index" json:"status_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (ClassSession) TableName() string { return "student.class_sessions" }

type StudentAttendance struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	SessionID uint       `gorm:"not null;index" json:"session_id"`
	StudentID uint       `gorm:"not null;index" json:"student_id"`
	StatusID  *uint      `gorm:"index" json:"status_id,omitempty"`
	MarkedBy  *uint      `json:"marked_by,omitempty"`
	MarkedAt  *time.Time `json:"marked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (StudentAttendance) TableName() string { return "student.attendance" }

type StudentEnrollment struct {
	ID                   uint    `gorm:"primaryKey" json:"id"`
	StudentID            uint    `gorm:"not null;index" json:"student_id"`
	CourseRegistrationID uint    `gorm:"not null" json:"course_registration_id"`
	EnrollmentDate       string  `gorm:"default:CURRENT_DATE;index" json:"enrollment_date"`
	StatusID             *uint   `gorm:"index" json:"status_id,omitempty"`
	Grade                string  `json:"grade"`
	MarksObtained        float64 `json:"marks_obtained"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

func (StudentEnrollment) TableName() string { return "student.student_enrollments" }

type Alumni struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	StudentID       uint      `gorm:"unique;not null;index" json:"student_id"`
	GraduationYear  int       `gorm:"index" json:"graduation_year"`
	CurrentEmployer string    `json:"current_employer"`
	JobTitle        string    `json:"job_title"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	LinkedInURL     string    `json:"linkedin_url"`
	IsSubscribed    bool      `gorm:"default:true" json:"is_subscribed"`
	CreatedAt       time.Time `json:"created_at"`
}

func (Alumni) TableName() string { return "student.alumni" }
