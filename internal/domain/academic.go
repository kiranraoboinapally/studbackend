package domain

import "time"

// ─── Academic Terms & Calendar ───────────────────────────────────────────────

type AcademicTerm struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	CampusID          *uint      `gorm:"index" json:"campus_id,omitempty"`
	AcademicYear      string     `gorm:"not null;index" json:"academic_year"`
	TermName          string     `gorm:"not null" json:"term_name"`
	StartDate         time.Time  `gorm:"not null" json:"start_date"`
	EndDate           time.Time  `gorm:"not null" json:"end_date"`
	RegistrationStart *time.Time `json:"registration_start,omitempty"`
	RegistrationEnd   *time.Time `json:"registration_end,omitempty"`
	ExamStartDate     *time.Time `json:"exam_start_date,omitempty"`
	ExamEndDate       *time.Time `json:"exam_end_date,omitempty"`
	IsCurrent         bool       `gorm:"default:false;index" json:"is_current"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func (AcademicTerm) TableName() string { return "academic.academic_terms" }

type AcademicCalendar struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CampusID    *uint     `gorm:"index" json:"campus_id,omitempty"`
	EventDate   time.Time `gorm:"not null;index" json:"event_date"`
	EventName   string    `gorm:"not null" json:"event_name"`
	EventType   string    `gorm:"type:varchar(50)" json:"event_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (AcademicCalendar) TableName() string { return "academic.academic_calendar" }

// ─── Programs & Curriculum ───────────────────────────────────────────────────

type Program struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DepartmentID   uint      `gorm:"not null;index" json:"department_id"`
	Name           string    `gorm:"not null" json:"name"`
	Code           string    `gorm:"unique;not null;index" json:"code"`
	DegreeType     string    `gorm:"type:varchar(50)" json:"degree_type"`
	DurationYears  int       `gorm:"not null" json:"duration_years"`
	TotalSemesters int       `gorm:"not null" json:"total_semesters"`
	TotalCredits   int       `gorm:"not null" json:"total_credits"`
	Description    string    `json:"description"`
	IsActive       bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Program) TableName() string { return "academic.programs" }

type ProgramSemester struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProgramID      uint      `gorm:"not null;index" json:"program_id"`
	SemesterNumber int       `gorm:"not null" json:"semester_number"`
	SemesterName   string    `gorm:"not null" json:"semester_name"`
	TotalCredits   int       `json:"total_credits"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

func (ProgramSemester) TableName() string { return "academic.program_semesters" }

type Subject struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	DepartmentID  uint    `gorm:"not null;index" json:"department_id"`
	SubjectCode   string  `gorm:"unique;not null;index" json:"subject_code"`
	SubjectName   string  `gorm:"not null" json:"subject_name"`
	Credits       float32 `gorm:"not null" json:"credits"`
	SubjectType   string  `gorm:"type:varchar(20)" json:"subject_type"`
	LectureHours  int     `json:"lecture_hours"`
	LabHours      int     `json:"lab_hours"`
	TutorialHours int     `json:"tutorial_hours"`
	Description   string  `json:"description"`
	IsActive      bool    `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Subject) TableName() string { return "academic.subjects" }

type ProgramSubject struct {
	ProgramID      uint `gorm:"primaryKey" json:"program_id"`
	SubjectID      uint `gorm:"primaryKey" json:"subject_id"`
	SemesterNumber int  `json:"semester_number"`
	IsCore         bool `gorm:"default:true" json:"is_core"`
}

func (ProgramSubject) TableName() string { return "academic.program_subjects" }

type SubjectPrerequisite struct {
	SubjectID             uint `gorm:"primaryKey" json:"subject_id"`
	PrerequisiteSubjectID uint `gorm:"primaryKey" json:"prerequisite_subject_id"`
}

func (SubjectPrerequisite) TableName() string { return "academic.subject_prerequisites" }

// ─── Batches, Sections & Registrations ───────────────────────────────────────

type Batch struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	ProgramID              uint      `gorm:"not null;index" json:"program_id"`
	BatchYear              int       `gorm:"not null;index" json:"batch_year"`
	AdmissionYear          int       `gorm:"not null" json:"admission_year"`
	ExpectedGraduationYear int       `json:"expected_graduation_year"`
	Status                 string    `gorm:"type:varchar(20);default:'Active'" json:"status"`
	CreatedAt              time.Time `json:"created_at"`
}

func (Batch) TableName() string { return "academic.batches" }

type Section struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	BatchID          uint      `gorm:"not null;index" json:"batch_id"`
	SectionName      string    `gorm:"not null" json:"section_name"`
	MentorEmployeeID *uint     `gorm:"index" json:"mentor_employee_id,omitempty"`
	MaxCapacity      int       `json:"max_capacity"`
	CreatedAt        time.Time `json:"created_at"`
}

func (Section) TableName() string { return "academic.sections" }

type CourseOffering struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProgramID         uint      `gorm:"not null;index" json:"program_id"`
	SubjectID         uint      `gorm:"not null;index" json:"subject_id"`
	AcademicTermID    uint      `gorm:"not null;index" json:"academic_term_id"`
	BatchID           uint      `gorm:"not null;index" json:"batch_id"`
	SectionID         *uint     `gorm:"index" json:"section_id,omitempty"`
	FacultyEmployeeID uint      `gorm:"not null;index" json:"faculty_employee_id"`
	RoomID            *uint     `gorm:"index" json:"room_id,omitempty"`
	MaxCapacity       int       `json:"max_capacity"`
	Status            string    `gorm:"type:varchar(20);default:'Active'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

func (CourseOffering) TableName() string { return "academic.course_offerings" }

type TermRegistration struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	StudentID         uint      `gorm:"not null;index" json:"student_id"`
	AcademicTermID    uint      `gorm:"not null;index" json:"academic_term_id"`
	BatchID           uint      `gorm:"not null;index" json:"batch_id"`
	SectionID         uint      `gorm:"not null;index" json:"section_id"`
	CurrentSemesterNo int       `json:"current_semester_no"`
	RegistrationDate  time.Time `json:"registration_date"`
	Status            string    `gorm:"type:varchar(20);default:'Active'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

func (TermRegistration) TableName() string { return "academic.term_registrations" }

type CourseRegistration struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	StudentID          uint      `gorm:"not null;index" json:"student_id"`
	OfferingID         uint      `gorm:"not null;index" json:"offering_id"`
	RegistrationStatus string    `gorm:"type:varchar(20);default:'Enrolled'" json:"registration_status"`
	IsRepeat           bool      `gorm:"default:false" json:"is_repeat"`
	IsElective         bool      `gorm:"default:false" json:"is_elective"`
	CreatedAt          time.Time `json:"created_at"`
}

func (CourseRegistration) TableName() string { return "academic.course_registrations" }

type Timetable struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OfferingID uint      `gorm:"not null;index" json:"offering_id"`
	DayOfWeek  int       `json:"day_of_week"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Timetable) TableName() string { return "academic.timetable" }
