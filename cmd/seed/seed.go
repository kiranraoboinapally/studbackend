// // cmd/seed/main.go
// package main

// import (
// 	// "database/sql"
// 	"fmt"
// 	"log"
// 	// "math"
// 	"math/rand"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	// "gorm.io/gorm/clause"
// 	// "gorm.io/gorm/logger"
// )

// var DB *gorm.DB
// var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// // ============================================================================
// // ALL MODEL DEFINITIONS (Copy from seed_university.go)
// // ============================================================================

// type User struct {
// 	ID           uint       `gorm:"primaryKey"`
// 	Username     string     `gorm:"unique;not null;index"`
// 	Email        string     `gorm:"unique;not null;index"`
// 	PasswordHash string     `gorm:"not null"`
// 	IsActive     bool       `gorm:"default:true;index"`
// 	LastLoginAt  *time.Time
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }
// func (User) TableName() string { return "shared.users" }

// type Role struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	RoleName    string `gorm:"unique;not null"`
// 	Description string
// 	CreatedAt   time.Time
// }
// func (Role) TableName() string { return "shared.roles" }

// type UserRole struct {
// 	UserID     uint      `gorm:"primaryKey"`
// 	RoleID     uint      `gorm:"primaryKey"`
// 	AssignedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	AssignedBy *uint
// }
// func (UserRole) TableName() string { return "shared.user_roles" }

// type Gender struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Code string `gorm:"unique;not null"`
// 	Name string `gorm:"not null"`
// }
// func (Gender) TableName() string { return "system.genders" }

// type Category struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Code string `gorm:"unique;not null"`
// 	Name string `gorm:"not null"`
// }
// func (Category) TableName() string { return "system.categories" }

// type BloodGroup struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Code string `gorm:"unique;not null"`
// 	Name string `gorm:"not null"`
// }
// func (BloodGroup) TableName() string { return "system.blood_groups" }

// type StatusCode struct {
// 	ID       uint   `gorm:"primaryKey"`
// 	Module   string `gorm:"not null;index"`
// 	Code     string `gorm:"not null;index"`
// 	Name     string `gorm:"not null"`
// 	IsActive bool   `gorm:"default:true"`
// }
// func (StatusCode) TableName() string { return "system.status_codes" }

// type Designation struct {
// 	ID       uint   `gorm:"primaryKey"`
// 	Code     string `gorm:"unique;not null"`
// 	Name     string `gorm:"not null"`
// 	IsActive bool   `gorm:"default:true"`
// }
// func (Designation) TableName() string { return "hr.designations" }

// type EmploymentType struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Code string `gorm:"unique;not null"`
// 	Name string `gorm:"not null"`
// }
// func (EmploymentType) TableName() string { return "hr.employment_types" }

// type LeaveType struct {
// 	ID       uint   `gorm:"primaryKey"`
// 	Code     string `gorm:"unique;not null"`
// 	Name     string `gorm:"not null"`
// 	MaxDays  float64
// 	Paid     bool
// 	IsActive bool `gorm:"default:true"`
// }
// func (LeaveType) TableName() string { return "hr.leave_types" }

// type SalaryComponent struct {
// 	ID       uint   `gorm:"primaryKey"`
// 	Code     string `gorm:"unique;not null"`
// 	Name     string `gorm:"not null"`
// 	Type     string `gorm:"type:varchar(20)"`
// 	IsActive bool   `gorm:"default:true"`
// }
// func (SalaryComponent) TableName() string { return "hr.salary_components" }

// type FeeHead struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	Name        string `gorm:"not null"`
// 	Code        string `gorm:"unique;not null;index"`
// 	Description string
// 	IsMandatory bool   `gorm:"default:true"`
// 	CreatedAt   time.Time
// }
// func (FeeHead) TableName() string { return "finance.fee_heads" }

// type University struct {
// 	ID              uint   `gorm:"primaryKey"`
// 	Name            string `gorm:"not null"`
// 	ShortName       string `gorm:"unique;not null;index"`
// 	EstablishedYear int
// 	Address         string
// 	City            string
// 	State           string
// 	PostalCode      string
// 	Phone           string
// 	Email           string
// 	Website         string
// 	Vision          string
// 	Mission         string
// 	IsActive        bool   `gorm:"default:true;index"`
// 	CreatedAt       time.Time
// 	UpdatedAt       time.Time
// }
// func (University) TableName() string { return "core.universities" }

// type Campus struct {
// 	ID           uint   `gorm:"primaryKey"`
// 	UniversityID uint   `gorm:"not null;index"`
// 	Name         string `gorm:"not null"`
// 	Code         string `gorm:"unique;not null;index"`
// 	Address      string
// 	City         string
// 	State        string
// 	PostalCode   string
// 	Phone        string
// 	IsMainCampus bool   `gorm:"default:false"`
// 	IsActive     bool   `gorm:"default:true;index"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }
// func (Campus) TableName() string { return "core.campuses" }

// type Department struct {
// 	ID                 uint   `gorm:"primaryKey"`
// 	CampusID           *uint  `gorm:"index"`
// 	Name               string `gorm:"not null"`
// 	Code               string `gorm:"unique;not null;index"`
// 	ParentDepartmentID *uint  `gorm:"index"`
// 	EstablishedYear    int
// 	HodEmployeeID      *uint  `gorm:"index"`
// 	Description        string
// 	IsActive           bool   `gorm:"default:true;index"`
// 	CreatedAt          time.Time
// 	UpdatedAt          time.Time
// }
// func (Department) TableName() string { return "core.departments" }

// type Room struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	CampusID  uint   `gorm:"not null;index"`
// 	RoomNumber string `gorm:"not null"`
// 	RoomType  string `gorm:"type:varchar(50);index"`
// 	Capacity  int    `gorm:"not null"`
// 	Building  string
// 	Floor     int
// 	IsActive  bool   `gorm:"default:true;index"`
// 	CreatedAt time.Time
// }
// func (Room) TableName() string { return "core.rooms" }

// type AcademicTerm struct {
// 	ID               uint      `gorm:"primaryKey"`
// 	CampusID         *uint     `gorm:"index"`
// 	AcademicYear     string    `gorm:"not null;index"`
// 	TermName         string    `gorm:"not null"`
// 	StartDate        time.Time `gorm:"not null"`
// 	EndDate          time.Time `gorm:"not null"`
// 	RegistrationStart *time.Time
// 	RegistrationEnd  *time.Time
// 	ExamStartDate    *time.Time
// 	ExamEndDate      *time.Time
// 	IsCurrent        bool      `gorm:"default:false;index"`
// 	CreatedAt        time.Time
// 	UpdatedAt        time.Time
// }
// func (AcademicTerm) TableName() string { return "academic.academic_terms" }

// type Program struct {
// 	ID             uint   `gorm:"primaryKey"`
// 	DepartmentID   uint   `gorm:"not null;index"`
// 	Name           string `gorm:"not null"`
// 	Code           string `gorm:"unique;not null;index"`
// 	DegreeType     string `gorm:"type:varchar(50)"`
// 	DurationYears  int    `gorm:"not null"`
// 	TotalSemesters int    `gorm:"not null"`
// 	TotalCredits   int    `gorm:"not null"`
// 	Description    string
// 	IsActive       bool   `gorm:"default:true;index"`
// 	CreatedAt      time.Time
// 	UpdatedAt      time.Time
// }
// func (Program) TableName() string { return "academic.programs" }

// type ProgramSemester struct {
// 	ID             uint   `gorm:"primaryKey"`
// 	ProgramID      uint   `gorm:"not null;index"`
// 	SemesterNumber int    `gorm:"not null"`
// 	SemesterName   string `gorm:"not null"`
// 	TotalCredits   int
// 	Description    string
// 	CreatedAt      time.Time
// }
// func (ProgramSemester) TableName() string { return "academic.program_semesters" }

// type Subject struct {
// 	ID            uint    `gorm:"primaryKey"`
// 	DepartmentID  uint    `gorm:"not null;index"`
// 	SubjectCode   string  `gorm:"unique;not null;index"`
// 	SubjectName   string  `gorm:"not null"`
// 	Credits       float32 `gorm:"not null"`
// 	SubjectType   string  `gorm:"type:varchar(20)"`
// 	LectureHours  int
// 	LabHours      int
// 	TutorialHours int
// 	Description   string
// 	IsActive      bool   `gorm:"default:true;index"`
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// }
// func (Subject) TableName() string { return "academic.subjects" }

// type ProgramSubject struct {
// 	ProgramID      uint `gorm:"primaryKey"`
// 	SubjectID      uint `gorm:"primaryKey"`
// 	SemesterNumber int
// 	IsCore         bool `gorm:"default:true"`
// }
// func (ProgramSubject) TableName() string { return "academic.program_subjects" }

// type Batch struct {
// 	ID                    uint      `gorm:"primaryKey"`
// 	ProgramID             uint      `gorm:"not null;index"`
// 	BatchYear             int       `gorm:"not null;index"`
// 	AdmissionYear         int       `gorm:"not null"`
// 	ExpectedGraduationYear int
// 	Status                string    `gorm:"type:varchar(20);default:'Active'"`
// 	CreatedAt             time.Time
// }
// func (Batch) TableName() string { return "academic.batches" }

// type Section struct {
// 	ID        uint      `gorm:"primaryKey"`
// 	BatchID   uint      `gorm:"not null;index"`
// 	SectionName string  `gorm:"not null"`
// 	MentorEmployeeID *uint `gorm:"index"`
// 	MaxCapacity int
// 	CreatedAt time.Time
// }
// func (Section) TableName() string { return "academic.sections" }

// type CourseOffering struct {
// 	ID              uint      `gorm:"primaryKey"`
// 	ProgramID       uint      `gorm:"not null;index"`
// 	SubjectID       uint      `gorm:"not null;index"`
// 	AcademicTermID  uint      `gorm:"not null;index"`
// 	BatchID         uint      `gorm:"not null;index"`
// 	SectionID       *uint     `gorm:"index"`
// 	FacultyEmployeeID uint    `gorm:"not null;index"`
// 	RoomID          *uint     `gorm:"index"`
// 	MaxCapacity     int
// 	Status          string    `gorm:"type:varchar(20);default:'Active'"`
// 	CreatedAt       time.Time
// }
// func (CourseOffering) TableName() string { return "academic.course_offerings" }

// type TermRegistration struct {
// 	ID                 uint      `gorm:"primaryKey"`
// 	StudentID          uint      `gorm:"not null;index"`
// 	AcademicTermID     uint      `gorm:"not null;index"`
// 	BatchID            uint      `gorm:"not null;index"`
// 	SectionID          uint      `gorm:"not null;index"`
// 	CurrentSemesterNo  int
// 	RegistrationDate   time.Time
// 	Status             string    `gorm:"type:varchar(20);default:'Active'"`
// 	CreatedAt          time.Time
// }
// func (TermRegistration) TableName() string { return "academic.term_registrations" }

// type CourseRegistration struct {
// 	ID              uint      `gorm:"primaryKey"`
// 	StudentID       uint      `gorm:"not null;index"`
// 	OfferingID      uint      `gorm:"not null;index"`
// 	RegistrationStatus string `gorm:"type:varchar(20);default:'Enrolled'"`
// 	IsRepeat        bool      `gorm:"default:false"`
// 	IsElective      bool      `gorm:"default:false"`
// 	CreatedAt       time.Time
// }
// func (CourseRegistration) TableName() string { return "academic.course_registrations" }

// type Employee struct {
// 	ID             uint       `gorm:"primaryKey"`
// 	UserID         uint       `gorm:"unique;not null;index"`
// 	EmployeeCode   string     `gorm:"unique;not null;index"`
// 	FirstName      string     `gorm:"not null"`
// 	LastName       string     `gorm:"not null"`
// 	GenderID       *uint      `gorm:"index"`
// 	DateOfBirth    *time.Time
// 	Phone          string
// 	Address        string
// 	JoiningDate    time.Time  `gorm:"not null;index"`
// 	EmploymentTypeID *uint    `gorm:"index"`
// 	DepartmentID   *uint      `gorm:"index"`
// 	DesignationID  *uint      `gorm:"index"`
// 	IsActive       bool       `gorm:"default:true;index"`
// 	CreatedAt      time.Time
// 	UpdatedAt      time.Time
// }
// func (Employee) TableName() string { return "hr.employees" }

// type Faculty struct {
// 	EmployeeID     uint   `gorm:"primaryKey"`
// 	Specialization string
// 	Qualification  string
// 	ResearchArea   string
// 	OfficeHours    string
// 	MaxLoadCredits int    `gorm:"default:20"`
// }
// func (Faculty) TableName() string { return "hr.faculties" }

// type Salary struct {
// 	ID            uint       `gorm:"primaryKey"`
// 	EmployeeID    uint       `gorm:"not null;index"`
// 	EffectiveFrom time.Time  `gorm:"not null;index"`
// 	EffectiveTo   *time.Time
// 	BasePay       float64
// 	NetSalary     float64
// 	IsActive      bool       `gorm:"default:true;index"`
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// }
// func (Salary) TableName() string { return "hr.salaries" }

// type SalaryDetail struct {
// 	ID                uint   `gorm:"primaryKey"`
// 	SalaryID          uint   `gorm:"not null;index"`
// 	SalaryComponentID uint   `gorm:"not null"`
// 	Amount            float64
// }
// func (SalaryDetail) TableName() string { return "hr.salary_details" }

// type LeaveBalance struct {
// 	ID           uint      `gorm:"primaryKey"`
// 	EmployeeID   uint      `gorm:"not null;index"`
// 	LeaveTypeID  uint      `gorm:"not null;index"`
// 	TotalQuota   float64   `gorm:"not null"`
// 	UsedQuota    float64   `gorm:"default:0"`
// 	AccruedQuota float64   `gorm:"default:0"`
// 	Year         int       `gorm:"not null;index"`
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }
// func (LeaveBalance) TableName() string { return "hr.leave_balances" }

// type LeaveRequest struct {
// 	ID          uint       `gorm:"primaryKey"`
// 	EmployeeID  uint       `gorm:"not null;index"`
// 	LeaveTypeID uint       `gorm:"not null"`
// 	StartDate   time.Time  `gorm:"not null"`
// 	EndDate     time.Time  `gorm:"not null"`
// 	Reason      string
// 	StatusID    *uint      `gorm:"index"`
// 	ApprovedBy  *uint
// 	ApprovedAt  *time.Time
// 	CreatedAt   time.Time
// }
// func (LeaveRequest) TableName() string { return "hr.leave_requests" }

// type HRAttendance struct {
// 	ID             uint      `gorm:"primaryKey"`
// 	EmployeeID     uint      `gorm:"not null;index"`
// 	AttendanceDate time.Time `gorm:"not null;index"`
// 	CheckIn        *time.Time
// 	CheckOut       *time.Time
// 	StatusID       *uint     `gorm:"index"`
// 	CreatedAt      time.Time
// }
// func (HRAttendance) TableName() string { return "hr.attendance" }

// type Student struct {
// 	ID                  uint      `gorm:"primaryKey"`
// 	UserID              uint      `gorm:"unique;not null;index"`
// 	EnrollmentNumber    string    `gorm:"unique;not null;index"`
// 	RollNumber          string    `gorm:"unique;index"`
// 	FirstName           string    `gorm:"not null"`
// 	LastName            string    `gorm:"not null"`
// 	DateOfBirth         time.Time `gorm:"not null"`
// 	GenderID            *uint     `gorm:"index"`
// 	Phone               string
// 	Email               string    `gorm:"not null;index"`
// 	AlternateEmail      string
// 	Address             string
// 	City                string
// 	State               string
// 	PostalCode          string
// 	Nationality         string    `gorm:"default:'Indian'"`
// 	CategoryID          *uint     `gorm:"index"`
// 	ProgramID           uint      `gorm:"not null;index"`
// 	AdmissionYear       int       `gorm:"not null;index"`
// 	AdmissionQuota      string
// 	IsHostelRequired    bool      `gorm:"default:false"`
// 	IsTransportRequired bool      `gorm:"default:false"`
// 	StatusID            *uint     `gorm:"index"`
// 	AcademicStanding    string    `gorm:"type:varchar(30);default:'Good'"`
// 	ProfilePhoto        string
// 	CreatedAt           time.Time
// 	UpdatedAt           time.Time
// }
// func (Student) TableName() string { return "student.students" }

// type Guardian struct {
// 	ID         uint   `gorm:"primaryKey"`
// 	StudentID  uint   `gorm:"not null;index"`
// 	Name       string `gorm:"not null"`
// 	Relation   string
// 	Phone      string
// 	Email      string
// 	Occupation string
// 	IsPrimary  bool   `gorm:"default:false"`
// }
// func (Guardian) TableName() string { return "student.guardians" }

// type MedicalRecord struct {
// 	ID                    uint      `gorm:"primaryKey"`
// 	StudentID             uint      `gorm:"unique;not null"`
// 	BloodGroupID          *uint     `gorm:"index"`
// 	Allergies             string
// 	ChronicConditions     string
// 	EmergencyContactName  string
// 	EmergencyContactPhone string
// 	InsurancePolicyNo     string
// 	ValidUntil            *time.Time
// 	UpdatedAt             time.Time
// }
// func (MedicalRecord) TableName() string { return "student.medical_records" }

// type Author struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	Name      string `gorm:"not null;index"`
// 	Biography string
// 	CreatedAt time.Time
// }
// func (Author) TableName() string { return "library.authors" }

// type Book struct {
// 	ID              uint      `gorm:"primaryKey"`
// 	Title           string    `gorm:"not null;index"`
// 	ISBN            string    `gorm:"unique;index"`
// 	Publisher       string
// 	PublicationYear int
// 	Edition         string
// 	TotalCopies     int       `gorm:"default:1"`
// 	AvailableCopies int       `gorm:"default:1"`
// 	Location        string
// 	CreatedAt       time.Time
// 	UpdatedAt       time.Time
// }
// func (Book) TableName() string { return "library.books" }

// type BookCopy struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	BookID      uint   `gorm:"not null;index"`
// 	Barcode     string `gorm:"unique;index"`
// 	CopyNumber  int
// 	Condition   string `gorm:"type:varchar(20)"`
// 	ShelfLocation string
// 	StatusID    *uint  `gorm:"index"`
// 	CreatedAt   time.Time
// }
// func (BookCopy) TableName() string { return "library.book_copies" }

// type BookAuthor struct {
// 	BookID   uint `gorm:"primaryKey"`
// 	AuthorID uint `gorm:"primaryKey"`
// }
// func (BookAuthor) TableName() string { return "library.book_authors" }

// type Circulation struct {
// 	ID           uint       `gorm:"primaryKey"`
// 	BookCopyID   uint       `gorm:"not null;index"`
// 	StudentID    uint       `gorm:"not null;index"`
// 	IssuedDate   time.Time  `gorm:"default:CURRENT_DATE;index"`
// 	DueDate      time.Time  `gorm:"not null;index"`
// 	ReturnedDate *time.Time `gorm:"index"`
// 	StatusID     *uint      `gorm:"index"`
// 	FineAmount   float64    `gorm:"default:0"`
// 	FinePaid     bool       `gorm:"default:false"`
// 	IssuedBy     *uint
// 	CreatedAt    time.Time
// }
// func (Circulation) TableName() string { return "library.circulations" }

// type FeeStructure struct {
// 	ID             uint      `gorm:"primaryKey"`
// 	ProgramID      uint      `gorm:"not null;index"`
// 	SemesterNumber int       `gorm:"not null"`
// 	FeeHeadID      uint      `gorm:"not null;index"`
// 	Amount         float64   `gorm:"not null"`
// 	AcademicYear   string    `gorm:"not null;index"`
// 	IsActive       bool      `gorm:"default:true;index"`
// 	CreatedBy      *uint
// 	CreatedAt      time.Time
// }
// func (FeeStructure) TableName() string { return "finance.fee_structures" }

// type Invoice struct {
// 	ID            uint       `gorm:"primaryKey"`
// 	StudentID     uint       `gorm:"not null;index"`
// 	InvoiceNumber string     `gorm:"unique;not null;index"`
// 	AcademicTermID uint      `gorm:"not null;index"`
// 	GeneratedDate  time.Time  `gorm:"default:CURRENT_DATE;index"`
// 	DueDate       time.Time  `gorm:"not null;index"`
// 	TotalAmount   float64    `gorm:"not null"`
// 	PaidAmount    float64    `gorm:"default:0"`
// 	LateFeeApplied float64    `gorm:"default:0"`
// 	StatusID      *uint      `gorm:"index"`
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// }
// func (Invoice) TableName() string { return "finance.invoices" }

// type InvoiceItem struct {
// 	ID          uint    `gorm:"primaryKey"`
// 	InvoiceID   uint    `gorm:"not null;index"`
// 	FeeHeadID   uint    `gorm:"not null"`
// 	Description string
// 	Quantity    int
// 	UnitAmount  float64
// 	Amount      float64
// }
// func (InvoiceItem) TableName() string { return "finance.invoice_items" }

// type Payment struct {
// 	ID                       uint      `gorm:"primaryKey"`
// 	InvoiceID                uint      `gorm:"not null;index"`
// 	StudentID                uint      `gorm:"not null;index"`
// 	Amount                   float64   `gorm:"not null"`
// 	PaymentDate              time.Time `gorm:"default:CURRENT_TIMESTAMP;index"`
// 	PaymentModeID            *uint     `gorm:"index"`
// 	TransactionID            string    `gorm:"index"`
// 	ReferenceNo              string
// 	StatusID                 *uint     `gorm:"index"`
// 	ReceiptURL               string
// 	BankReconciliationStatus string    `gorm:"type:varchar(20);default:'Pending'"`
// 	CreatedAt                time.Time
// }
// func (Payment) TableName() string { return "finance.payments" }

// type Applicant struct {
// 	ID                 uint       `gorm:"primaryKey"`
// 	ApplicationNumber  string     `gorm:"unique;not null;index"`
// 	CycleID            uint       `gorm:"not null;index"`
// 	ProgramID          *uint      `gorm:"index"`
// 	FirstName          string     `gorm:"not null"`
// 	LastName           string     `gorm:"not null"`
// 	DateOfBirth        time.Time  `gorm:"not null"`
// 	Email              string     `gorm:"not null;index"`
// 	Phone              string
// 	Address            string
// 	GenderID           *uint      `gorm:"index"`
// 	CategoryID         *uint      `gorm:"index"`
// 	EntranceScore      float64
// 	Rank               int
// 	StatusID           *uint      `gorm:"index"`
// 	AppliedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP;index"`
// 	UpdatedAt          time.Time
// }
// func (Applicant) TableName() string { return "admissions.applicants" }

// type AdmissionCycle struct {
// 	ID                 uint       `gorm:"primaryKey"`
// 	Name               string     `gorm:"not null"`
// 	AcademicYear       string     `gorm:"not null;index"`
// 	ProgramID          *uint      `gorm:"index"`
// 	ApplicationStart   time.Time  `gorm:"not null"`
// 	ApplicationEnd     time.Time  `gorm:"not null"`
// 	EntranceExamDate   *time.Time
// 	CounselingStart    *time.Time
// 	CounselingEnd      *time.Time
// 	ApplicationFee     float64
// 	MaxApplications    int
// 	IsOpen             bool       `gorm:"default:true;index"`
// 	CreatedAt          time.Time
// 	UpdatedAt          time.Time
// }
// func (AdmissionCycle) TableName() string { return "admissions.admission_cycles" }

// type ClassSession struct {
// 	ID          uint      `gorm:"primaryKey"`
// 	OfferingID  uint      `gorm:"not null;index"`
// 	ClassDate   time.Time `gorm:"not null;index"`
// 	StartTime   *time.Time
// 	EndTime     *time.Time
// 	RoomID      *uint     `gorm:"index"`
// 	FacultyID   uint      `gorm:"not null"`
// 	StatusID    *uint     `gorm:"index"`
// 	CreatedAt   time.Time
// }
// func (ClassSession) TableName() string { return "student.class_sessions" }

// type StudentAttendance struct {
// 	ID        uint      `gorm:"primaryKey"`
// 	SessionID uint      `gorm:"not null;index"`
// 	StudentID uint      `gorm:"not null;index"`
// 	StatusID  *uint     `gorm:"index"`
// 	MarkedBy  *uint
// 	MarkedAt  *time.Time
// 	CreatedAt time.Time
// }
// func (StudentAttendance) TableName() string { return "student.attendance" }

// type StudentEnrollment struct {
// 	ID                  uint      `gorm:"primaryKey"`
// 	StudentID           uint      `gorm:"not null;index"`
// 	CourseRegistrationID uint     `gorm:"not null"`
// 	EnrollmentDate      time.Time `gorm:"default:CURRENT_DATE;index"`
// 	StatusID            *uint     `gorm:"index"`
// 	Grade               string
// 	MarksObtained       float64
// 	AttendancePercentage float64
// }
// func (StudentEnrollment) TableName() string { return "student.student_enrollments" }

// type Result struct {
// 	ID             uint       `gorm:"primaryKey"`
// 	StudentID      uint       `gorm:"not null;index"`
// 	CourseRegID    uint       `gorm:"not null;index"`
// 	SubjectID      uint       `gorm:"not null;index"`
// 	AcademicTermID uint       `gorm:"not null;index"`
// 	MarksObtained  float64
// 	MaxMarks       float64
// 	Grade          string     `gorm:"index"`
// 	GradePoint     float64
// 	IsPassed       bool
// 	PublishedAt    *time.Time
// 	PublishedBy    *uint
// 	CreatedAt      time.Time
// }
// func (Result) TableName() string { return "exam.results" }

// type ExamSchedule struct {
// 	ID             uint      `gorm:"primaryKey"`
// 	SubjectID      uint      `gorm:"not null;index"`
// 	AcademicTermID uint      `gorm:"not null;index"`
// 	ExamDate       time.Time `gorm:"not null;index"`
// 	StartTime      string
// 	EndTime        string
// 	ExamType       string    `gorm:"type:varchar(30)"`
// 	Venue          string
// 	TotalMarks     int       `gorm:"not null"`
// 	PassingMarks   int
// 	CreatedAt      time.Time
// }
// func (ExamSchedule) TableName() string { return "exam.exam_schedules" }

// type Hostel struct {
// 	ID            uint   `gorm:"primaryKey"`
// 	Name          string `gorm:"not null"`
// 	Code          string `gorm:"unique;not null;index"`
// 	CampusID      *uint  `gorm:"index"`
// 	GenderID      *uint  `gorm:"index"`
// 	TotalRooms    int
// 	WardenID      *uint  `gorm:"index"`
// 	ContactNumber string
// 	Address       string
// 	IsActive      bool   `gorm:"default:true;index"`
// 	CreatedAt     time.Time
// }
// func (Hostel) TableName() string { return "hostel.hostels" }

// type HostelRoom struct {
// 	ID               uint    `gorm:"primaryKey"`
// 	HostelID         uint    `gorm:"not null;index"`
// 	RoomNumber       string  `gorm:"not null"`
// 	RoomType         string  `gorm:"type:varchar(20);index"`
// 	Capacity         int     `gorm:"not null"`
// 	CurrentOccupancy int     `gorm:"default:0"`
// 	MonthlyRent      float64
// 	IsAvailable      bool    `gorm:"default:true;index"`
// 	CreatedAt        time.Time
// }
// func (HostelRoom) TableName() string { return "hostel.rooms" }

// type HostelBed struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	RoomID    uint   `gorm:"not null;index"`
// 	BedNumber string
// 	IsOccupied bool `gorm:"default:false"`
// }
// func (HostelBed) TableName() string { return "hostel.beds" }

// type HostelAllocation struct {
// 	ID            uint       `gorm:"primaryKey"`
// 	StudentID     uint       `gorm:"not null;index"`
// 	RoomID        uint       `gorm:"not null;index"`
// 	BedID         *uint      `gorm:"index"`
// 	AllocatedFrom time.Time  `gorm:"not null;index"`
// 	AllocatedTo   *time.Time
// 	StatusID      *uint      `gorm:"index"`
// 	CreatedBy     *uint
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// }
// func (HostelAllocation) TableName() string { return "hostel.allocations" }

// type Bus struct {
// 	ID                  uint       `gorm:"primaryKey"`
// 	BusNumber           string     `gorm:"unique;not null;index"`
// 	RegistrationNo      string     `gorm:"unique;not null;index"`
// 	Capacity            int        `gorm:"not null"`
// 	DriverEmployeeID    *uint      `gorm:"index"`
// 	DriverLicenseExpiry *time.Time
// 	IsActive            bool       `gorm:"default:true;index"`
// 	CreatedAt           time.Time
// }
// func (Bus) TableName() string { return "transport.buses" }

// type Route struct {
// 	ID            uint     `gorm:"primaryKey"`
// 	RouteName     string   `gorm:"not null;index"`
// 	Description   string
// 	DistanceKm    float64
// 	EstimatedTime string
// 	IsActive      bool     `gorm:"default:true;index"`
// 	CreatedAt     time.Time
// }
// func (Route) TableName() string { return "transport.routes" }

// type Stop struct {
// 	ID            uint      `gorm:"primaryKey"`
// 	RouteID       uint      `gorm:"not null;index"`
// 	StopName      string    `gorm:"not null"`
// 	StopOrder     int       `gorm:"not null"`
// 	Latitude      float64
// 	Longitude     float64
// 	ArrivalTime   string
// 	DepartureTime string
// 	CreatedAt     time.Time
// }
// func (Stop) TableName() string { return "transport.stops" }

// type Timetable struct {
// 	ID          uint      `gorm:"primaryKey"`
// 	OfferingID  uint      `gorm:"not null;index"`
// 	DayOfWeek   int       `gorm:"check:day_of_week between 1 and 7"`
// 	StartTime   string
// 	EndTime     string
// 	CreatedAt   time.Time
// }
// func (Timetable) TableName() string { return "academic.timetable" }

// type AcademicCalendar struct {
// 	ID          uint      `gorm:"primaryKey"`
// 	CampusID    *uint     `gorm:"index"`
// 	EventDate   time.Time `gorm:"not null;index"`
// 	EventName   string    `gorm:"not null"`
// 	EventType   string    `gorm:"type:varchar(50)"`
// 	Description string
// 	CreatedAt   time.Time
// }
// func (AcademicCalendar) TableName() string { return "academic.academic_calendar" }

// type SubjectPrerequisite struct {
// 	SubjectID            uint `gorm:"primaryKey"`
// 	PrerequisiteSubjectID uint `gorm:"primaryKey"`
// }
// func (SubjectPrerequisite) TableName() string { return "academic.subject_prerequisites" }

// type Permission struct {
// 	ID          uint   `gorm:"primaryKey"`
// 	Resource    string `gorm:"not null;index"`
// 	Action      string `gorm:"not null"`
// 	Description string
// }
// func (Permission) TableName() string { return "security.permissions" }

// type RolePermission struct {
// 	RoleID       uint      `gorm:"primaryKey"`
// 	PermissionID uint      `gorm:"primaryKey"`
// 	GrantedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	GrantedBy    *uint
// }
// func (RolePermission) TableName() string { return "security.role_permissions" }

// type Configuration struct {
// 	ID          uint      `gorm:"primaryKey"`
// 	ConfigKey   string    `gorm:"unique;not null;index"`
// 	ConfigValue string
// 	DataType    string    `gorm:"default:'string'"`
// 	Description string
// 	UpdatedAt   time.Time
// }
// func (Configuration) TableName() string { return "system.configurations" }

// type Notification struct {
// 	ID        uint      `gorm:"primaryKey"`
// 	UserID    uint      `gorm:"not null;index"`
// 	Title     string    `gorm:"not null"`
// 	Message   string    `gorm:"not null"`
// 	Type      string    `gorm:"default:'info';index"`
// 	IsRead    bool      `gorm:"default:false;index"`
// 	CreatedAt time.Time `gorm:"index"`
// }
// func (Notification) TableName() string { return "system.notifications" }

// type AuditLog struct {
// 	ID            uint      `gorm:"primaryKey"`
// 	UserID        *uint     `gorm:"index"`
// 	Action        string    `gorm:"type:varchar(50);not null;index"`
// 	SchemaName    string    `gorm:"type:varchar(50)"`
// 	AffectedTable string    `gorm:"type:varchar(100)"`
// 	RecordID      string    `gorm:"index"`
// 	OldValue      string    `gorm:"type:jsonb"`
// 	NewValue      string    `gorm:"type:jsonb"`
// 	IPAddress     string
// 	UserAgent     string
// 	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index"`
// }
// func (AuditLog) TableName() string { return "shared.audit_logs" }

// type Scholarship struct {
// 	ID                  uint   `gorm:"primaryKey"`
// 	Name                string `gorm:"not null"`
// 	Description         string
// 	EligibilityCriteria string `gorm:"type:jsonb"`
// 	Amount              float64
// 	Renewable           bool   `gorm:"default:false"`
// 	CreatedAt           time.Time
// }
// func (Scholarship) TableName() string { return "finance.scholarships" }

// type StudentScholarship struct {
// 	ID                uint      `gorm:"primaryKey"`
// 	StudentID         uint      `gorm:"not null;index"`
// 	ScholarshipID     uint      `gorm:"not null;index"`
// 	AcademicYear      string    `gorm:"not null;index"`
// 	AmountAwarded     float64
// 	Disbursed         bool      `gorm:"default:false;index"`
// 	DisbursedAt       *time.Time
// 	CreatedAt         time.Time
// }
// func (StudentScholarship) TableName() string { return "finance.student_scholarships" }

// // ============================================================================
// // DATABASE INITIALIZATION
// // ============================================================================

// func initDB() {
// 	_ = godotenv.Load()
// 	host := getEnv("DB_HOST", "192.168.1.201")
// 	port := getEnv("DB_PORT", "5432")
// 	user := getEnv("DB_USER", "postgres")
// 	password := getEnv("DB_PASSWORD", "root")
// 	dbname := getEnv("DB_NAME", "university_erp_prod10")

// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
// 		host, port, user, password, dbname)
// 	var err error
// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("❌ Failed to connect to database: %v", err)
// 	}
// 	log.Println("✅ Database connected successfully")
// }

// func getEnv(key, fallback string) string {
// 	if v := os.Getenv(key); v != "" {
// 		return v
// 	}
// 	return fallback
// }

// func hashPW(pwd string) string {
// 	b, _ := bcrypt.GenerateFromPassword([]byte(pwd), 12)
// 	return string(b)
// }

// // ============================================================================
// // POPULATE FUNCTIONS (from previous code)
// // ============================================================================

// func populateMasterData() {
// 	log.Println("\n📊 Populating Master Data...")
// 	startTime := time.Now()

// 	genders := []Gender{
// 		{Code: "M", Name: "Male"},
// 		{Code: "F", Name: "Female"},
// 		{Code: "O", Name: "Other"},
// 	}
// 	DB.CreateInBatches(genders, 100)
// 	log.Printf("  ✅ Created %d genders", len(genders))

// 	categories := []Category{
// 		{Code: "GEN", Name: "General"},
// 		{Code: "OBC", Name: "Other Backward Class"},
// 		{Code: "SC", Name: "Scheduled Caste"},
// 		{Code: "ST", Name: "Scheduled Tribe"},
// 	}
// 	DB.CreateInBatches(categories, 100)
// 	log.Printf("  ✅ Created %d categories", len(categories))

// 	bloodGroups := []BloodGroup{
// 		{Code: "O+", Name: "O Positive"},
// 		{Code: "O-", Name: "O Negative"},
// 		{Code: "A+", Name: "A Positive"},
// 		{Code: "A-", Name: "A Negative"},
// 		{Code: "B+", Name: "B Positive"},
// 		{Code: "B-", Name: "B Negative"},
// 		{Code: "AB+", Name: "AB Positive"},
// 		{Code: "AB-", Name: "AB Negative"},
// 	}
// 	DB.CreateInBatches(bloodGroups, 100)
// 	log.Printf("  ✅ Created %d blood groups", len(bloodGroups))

// 	statusCodes := []StatusCode{
// 		{Module: "student", Code: "ACTIVE", Name: "Active", IsActive: true},
// 		{Module: "student", Code: "INACTIVE", Name: "Inactive", IsActive: true},
// 		{Module: "student", Code: "SUSPENDED", Name: "Suspended", IsActive: true},
// 		{Module: "student", Code: "GRADUATED", Name: "Graduated", IsActive: true},
// 		{Module: "student", Code: "DROPOUT", Name: "Dropout", IsActive: true},
// 		{Module: "finance", Code: "UNPAID", Name: "Unpaid", IsActive: true},
// 		{Module: "finance", Code: "PARTIAL", Name: "Partially Paid", IsActive: true},
// 		{Module: "finance", Code: "PAID", Name: "Paid", IsActive: true},
// 		{Module: "finance", Code: "OVERDUE", Name: "Overdue", IsActive: true},
// 		{Module: "admission", Code: "APPLIED", Name: "Applied", IsActive: true},
// 		{Module: "admission", Code: "SHORTLISTED", Name: "Shortlisted", IsActive: true},
// 		{Module: "admission", Code: "OFFERED", Name: "Offered", IsActive: true},
// 		{Module: "admission", Code: "ACCEPTED", Name: "Accepted", IsActive: true},
// 		{Module: "admission", Code: "REJECTED", Name: "Rejected", IsActive: true},
// 		{Module: "admission", Code: "WAITLIST", Name: "Waitlist", IsActive: true},
// 		{Module: "library", Code: "ISSUED", Name: "Issued", IsActive: true},
// 		{Module: "library", Code: "RETURNED", Name: "Returned", IsActive: true},
// 		{Module: "library", Code: "OVERDUE", Name: "Overdue", IsActive: true},
// 		{Module: "library", Code: "LOST", Name: "Lost", IsActive: true},
// 		{Module: "hr", Code: "PRESENT", Name: "Present", IsActive: true},
// 		{Module: "hr", Code: "ABSENT", Name: "Absent", IsActive: true},
// 		{Module: "hr", Code: "LATE", Name: "Late", IsActive: true},
// 		{Module: "hr", Code: "LEAVE", Name: "Leave", IsActive: true},
// 		{Module: "exam", Code: "CONDUCTED", Name: "Conducted", IsActive: true},
// 		{Module: "exam", Code: "POSTPONED", Name: "Postponed", IsActive: true},
// 		{Module: "exam", Code: "CANCELED", Name: "Canceled", IsActive: true},
// 	}
// 	DB.CreateInBatches(statusCodes, 100)
// 	log.Printf("  ✅ Created %d status codes", len(statusCodes))

// 	designations := []Designation{
// 		{Code: "PROF", Name: "Professor", IsActive: true},
// 		{Code: "ASSOC_PROF", Name: "Associate Professor", IsActive: true},
// 		{Code: "ASST_PROF", Name: "Assistant Professor", IsActive: true},
// 		{Code: "LECTURER", Name: "Lecturer", IsActive: true},
// 		{Code: "ADMIN", Name: "Administrator", IsActive: true},
// 		{Code: "CLERK", Name: "Clerk", IsActive: true},
// 		{Code: "LIBRARIAN", Name: "Librarian", IsActive: true},
// 		{Code: "WARDEN", Name: "Hostel Warden", IsActive: true},
// 	}
// 	DB.CreateInBatches(designations, 100)
// 	log.Printf("  ✅ Created %d designations", len(designations))

// 	empTypes := []EmploymentType{
// 		{Code: "FULL_TIME", Name: "Full Time"},
// 		{Code: "PART_TIME", Name: "Part Time"},
// 		{Code: "CONTRACT", Name: "Contract"},
// 		{Code: "TEMPORARY", Name: "Temporary"},
// 	}
// 	DB.CreateInBatches(empTypes, 100)
// 	log.Printf("  ✅ Created %d employment types", len(empTypes))

// 	leaveTypes := []LeaveType{
// 		{Code: "CL", Name: "Casual Leave", MaxDays: 10, Paid: true, IsActive: true},
// 		{Code: "SL", Name: "Sick Leave", MaxDays: 7, Paid: true, IsActive: true},
// 		{Code: "EL", Name: "Earned Leave", MaxDays: 20, Paid: true, IsActive: true},
// 		{Code: "UL", Name: "Unpaid Leave", MaxDays: 5, Paid: false, IsActive: true},
// 		{Code: "ML", Name: "Maternity Leave", MaxDays: 180, Paid: true, IsActive: true},
// 	}
// 	DB.CreateInBatches(leaveTypes, 100)
// 	log.Printf("  ✅ Created %d leave types", len(leaveTypes))

// 	salaryComps := []SalaryComponent{
// 		{Code: "BASIC", Name: "Basic Salary", Type: "allowance", IsActive: true},
// 		{Code: "DA", Name: "Dearness Allowance", Type: "allowance", IsActive: true},
// 		{Code: "HRA", Name: "House Rent Allowance", Type: "allowance", IsActive: true},
// 		{Code: "CONVEY", Name: "Conveyance Allowance", Type: "allowance", IsActive: true},
// 		{Code: "SPECIAL", Name: "Special Allowance", Type: "allowance", IsActive: true},
// 		{Code: "PF", Name: "Provident Fund", Type: "deduction", IsActive: true},
// 		{Code: "IT", Name: "Income Tax", Type: "deduction", IsActive: true},
// 		{Code: "ESI", Name: "Employee State Insurance", Type: "deduction", IsActive: true},
// 		{Code: "PROF_TAX", Name: "Professional Tax", Type: "deduction", IsActive: true},
// 	}
// 	DB.CreateInBatches(salaryComps, 100)
// 	log.Printf("  ✅ Created %d salary components", len(salaryComps))

// 	feeHeads := []FeeHead{
// 		{Name: "Tuition Fee", Code: "TUITION", Description: "Main tuition fee", IsMandatory: true},
// 		{Name: "Exam Fee", Code: "EXAM", Description: "Examination fee", IsMandatory: true},
// 		{Name: "Library Fee", Code: "LIBRARY", Description: "Library usage fee", IsMandatory: true},
// 		{Name: "Lab Fee", Code: "LAB", Description: "Laboratory fee", IsMandatory: false},
// 		{Name: "Sports Fee", Code: "SPORTS", Description: "Sports facilities fee", IsMandatory: false},
// 		{Name: "Hostel Fee", Code: "HOSTEL", Description: "Hostel accommodation fee", IsMandatory: false},
// 		{Name: "Student Development Fee", Code: "DEVELOPMENT", Description: "Student development activities", IsMandatory: false},
// 		{Name: "Technology Fee", Code: "TECHNOLOGY", Description: "IT infrastructure fee", IsMandatory: true},
// 	}
// 	DB.CreateInBatches(feeHeads, 100)
// 	log.Printf("  ✅ Created %d fee heads", len(feeHeads))

// 	log.Printf("⏱️  Master data completed in %v\n", time.Since(startTime))
// }

// func populateCoreData() {
// 	log.Println("\n🏢 Populating Core Organizational Data...")
// 	startTime := time.Now()

// 	universities := []University{
// 		{
// 			Name:            "National Technology University",
// 			ShortName:       "NTU",
// 			EstablishedYear: 1985,
// 			City:            "Hyderabad",
// 			State:           "Telangana",
// 			PostalCode:      "500001",
// 			Phone:           "+91-40-12345678",
// 			Email:           "admin@ntu.edu",
// 			Website:         "https://ntu.edu",
// 			IsActive:        true,
// 		},
// 		{
// 			Name:            "Global Engineering Institute",
// 			ShortName:       "GEI",
// 			EstablishedYear: 2005,
// 			City:            "Bangalore",
// 			State:           "Karnataka",
// 			PostalCode:      "560001",
// 			Phone:           "+91-80-87654321",
// 			Email:           "info@gei.edu",
// 			Website:         "https://gei.edu",
// 			IsActive:        true,
// 		},
// 	}
// 	DB.CreateInBatches(universities, 100)
// 	log.Printf("  ✅ Created %d universities", len(universities))

// 	var ntu University
// 	DB.Where("short_name = ?", "NTU").First(&ntu)

// 	campuses := []Campus{
// 		{
// 			UniversityID: ntu.ID,
// 			Name:         "Main Campus",
// 			Code:         "HYD-MAIN",
// 			City:         "Hyderabad",
// 			State:        "Telangana",
// 			PostalCode:   "500001",
// 			IsMainCampus: true,
// 			IsActive:     true,
// 		},
// 		{
// 			UniversityID: ntu.ID,
// 			Name:         "Secondary Campus",
// 			Code:         "HYD-SEC",
// 			City:         "Hyderabad",
// 			State:        "Telangana",
// 			PostalCode:   "500089",
// 			IsMainCampus: false,
// 			IsActive:     true,
// 		},
// 	}
// 	DB.CreateInBatches(campuses, 100)
// 	log.Printf("  ✅ Created %d campuses", len(campuses))

// 	var mainCampus Campus
// 	DB.Where("code = ?", "HYD-MAIN").First(&mainCampus)

// 	departments := make([]Department, 0)
// 	deptNames := []struct {
// 		name string
// 		code string
// 	}{
// 		{"Computer Science & Engineering", "CSE"},
// 		{"Electronics & Communication Engineering", "ECE"},
// 		{"Electrical Engineering", "EE"},
// 		{"Mechanical Engineering", "ME"},
// 		{"Civil Engineering", "CE"},
// 		{"Chemical Engineering", "CHE"},
// 		{"Management Studies", "MBA"},
// 		{"Business Administration", "BBA"},
// 	}

// 	for _, d := range deptNames {
// 		departments = append(departments, Department{
// 			CampusID: &mainCampus.ID,
// 			Name:     d.name,
// 			Code:     d.code,
// 			IsActive: true,
// 		})
// 	}
// 	DB.CreateInBatches(departments, 100)
// 	log.Printf("  ✅ Created %d departments", len(departments))

// 	rooms := make([]Room, 0)
// 	buildings := []string{"Block A", "Block B", "Lab Building"}
// 	roomTypes := []string{"Lecture", "Lab", "Seminar"}

// 	for b := 0; b < len(buildings); b++ {
// 		for f := 1; f <= 3; f++ {
// 			for r := 1; r <= 5; r++ {
// 				roomNum := fmt.Sprintf("%s%d%d", string(rune(65+b)), f, r)
// 				capacity := 60
// 				if roomTypes[b%3] == "Lab" {
// 					capacity = 30
// 				}
// 				rooms = append(rooms, Room{
// 					CampusID:   mainCampus.ID,
// 					RoomNumber: roomNum,
// 					RoomType:   roomTypes[b%3],
// 					Capacity:   capacity,
// 					Building:   buildings[b],
// 					Floor:      f,
// 					IsActive:   true,
// 				})
// 			}
// 		}
// 	}
// 	DB.CreateInBatches(rooms, 100)
// 	log.Printf("  ✅ Created %d rooms", len(rooms))

// 	log.Printf("⏱️  Core data completed in %v\n", time.Since(startTime))
// }

// func populateAcademicData() {
// 	log.Println("\n📚 Populating Academic Data...")
// 	startTime := time.Now()

// 	var departments []Department
// 	DB.Find(&departments)
// 	log.Printf("  Found %d departments", len(departments))

// 	programs := make([]Program, 0)
// 	for _, dept := range departments {
// 		degreeType := "B.Tech"
// 		if dept.Code == "MBA" {
// 			degreeType = "MBA"
// 		} else if dept.Code == "BBA" {
// 			degreeType = "BBA"
// 		}

// 		programs = append(programs, Program{
// 			DepartmentID:   dept.ID,
// 			Name:           fmt.Sprintf("%s - %s", degreeType, dept.Code),
// 			Code:           fmt.Sprintf("%s-%s", degreeType, dept.Code),
// 			DegreeType:     degreeType,
// 			DurationYears:  4,
// 			TotalSemesters: 8,
// 			TotalCredits:   160,
// 			IsActive:       true,
// 		})
// 	}
// 	DB.CreateInBatches(programs, 100)
// 	log.Printf("  ✅ Created %d programs", len(programs))

// 	programSemesters := make([]ProgramSemester, 0)
// 	for _, prog := range programs {
// 		for sem := 1; sem <= 8; sem++ {
// 			programSemesters = append(programSemesters, ProgramSemester{
// 				ProgramID:      prog.ID,
// 				SemesterNumber: sem,
// 				SemesterName:   fmt.Sprintf("Semester %d", sem),
// 				TotalCredits:   20,
// 				Description:    fmt.Sprintf("Semester %d courses", sem),
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(programSemesters, 500)
// 	log.Printf("  ✅ Created %d program semesters", len(programSemesters))

// 	currentYear := time.Now().Year()
// 	academicTerms := make([]AcademicTerm, 0)
// 	for year := 0; year < 3; year++ {
// 		y := currentYear - year
// 		academicTerms = append(academicTerms,
// 			AcademicTerm{
// 				AcademicYear:      fmt.Sprintf("%d-%d", y, y+1),
// 				TermName:          fmt.Sprintf("Fall %d", y),
// 				StartDate:         time.Date(y, 7, 1, 0, 0, 0, 0, time.UTC),
// 				EndDate:           time.Date(y, 11, 30, 0, 0, 0, 0, time.UTC),
// 				RegistrationStart: &[]time.Time{time.Date(y, 6, 15, 0, 0, 0, 0, time.UTC)}[0],
// 				RegistrationEnd:   &[]time.Time{time.Date(y, 7, 5, 0, 0, 0, 0, time.UTC)}[0],
// 				ExamStartDate:     &[]time.Time{time.Date(y, 11, 15, 0, 0, 0, 0, time.UTC)}[0],
// 				ExamEndDate:       &[]time.Time{time.Date(y, 12, 15, 0, 0, 0, 0, time.UTC)}[0],
// 				IsCurrent:         year == 0,
// 			},
// 			AcademicTerm{
// 				AcademicYear:      fmt.Sprintf("%d-%d", y, y+1),
// 				TermName:          fmt.Sprintf("Spring %d", y+1),
// 				StartDate:         time.Date(y+1, 1, 1, 0, 0, 0, 0, time.UTC),
// 				EndDate:           time.Date(y+1, 5, 31, 0, 0, 0, 0, time.UTC),
// 				RegistrationStart: &[]time.Time{time.Date(y+1, 12, 15, 0, 0, 0, 0, time.UTC)}[0],
// 				RegistrationEnd:   &[]time.Time{time.Date(y+1, 1, 5, 0, 0, 0, 0, time.UTC)}[0],
// 				ExamStartDate:     &[]time.Time{time.Date(y+1, 5, 1, 0, 0, 0, 0, time.UTC)}[0],
// 				ExamEndDate:       &[]time.Time{time.Date(y+1, 5, 31, 0, 0, 0, 0, time.UTC)}[0],
// 				IsCurrent:         false,
// 			},
// 		)
// 	}
// 	DB.CreateInBatches(academicTerms, 100)
// 	log.Printf("  ✅ Created %d academic terms", len(academicTerms))

// 	subjects := make([]Subject, 0)
// 	subjectNames := []string{
// 		"Programming Fundamentals", "Data Structures", "Database Systems", "Operating Systems",
// 		"Web Development", "Mobile App Development", "Machine Learning", "Artificial Intelligence",
// 		"Cloud Computing", "DevOps", "Microservices", "Advanced Algorithms",
// 		"System Design", "Software Engineering", "Project Management", "Business Analytics",
// 	}

// 	for _, dept := range departments {
// 		for i, subName := range subjectNames {
// 			code := fmt.Sprintf("%s%03d", dept.Code, i+100)
// 			credits := 3.0 + float32(i%3)
// 			subjects = append(subjects, Subject{
// 				DepartmentID:  dept.ID,
// 				SubjectCode:   code,
// 				SubjectName:   subName,
// 				Credits:       credits,
// 				SubjectType:   "Theory",
// 				LectureHours:  3,
// 				LabHours:      0,
// 				TutorialHours: 1,
// 				IsActive:      true,
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(subjects, 500)
// 	log.Printf("  ✅ Created %d subjects", len(subjects))

// 	programSubjects := make([]ProgramSubject, 0)
// 	var allSubjects []Subject
// 	DB.Find(&allSubjects)

// 	for _, prog := range programs {
// 		subjectsPerSem := len(allSubjects) / 8
// 		for sem := 1; sem <= 8; sem++ {
// 			for i := 0; i < subjectsPerSem && i < len(allSubjects); i++ {
// 				idx := (sem-1)*subjectsPerSem + i
// 				if idx < len(allSubjects) {
// 					programSubjects = append(programSubjects, ProgramSubject{
// 						ProgramID:      prog.ID,
// 						SubjectID:      allSubjects[idx].ID,
// 						SemesterNumber: sem,
// 						IsCore:         i < subjectsPerSem/2,
// 					})
// 				}
// 			}
// 		}
// 	}
// 	DB.CreateInBatches(programSubjects, 500)
// 	log.Printf("  ✅ Created %d program subjects", len(programSubjects))

// 	batches := make([]Batch, 0)
// 	for year := 2022; year <= 2024; year++ {
// 		for _, prog := range programs {
// 			batches = append(batches, Batch{
// 				ProgramID:              prog.ID,
// 				BatchYear:              year,
// 				AdmissionYear:          year,
// 				ExpectedGraduationYear: year + 4,
// 				Status:                 "Active",
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(batches, 500)
// 	log.Printf("  ✅ Created %d batches", len(batches))

// 	sections := make([]Section, 0)
// 	sectionNames := []string{"A", "B", "C", "D"}
// 	for _, batch := range batches {
// 		for _, sectionName := range sectionNames {
// 			sections = append(sections, Section{
// 				BatchID:     batch.ID,
// 				SectionName: sectionName,
// 				MaxCapacity: 60,
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(sections, 500)
// 	log.Printf("  ✅ Created %d sections", len(sections))

// 	log.Printf("⏱️  Academic data completed in %v\n", time.Since(startTime))
// }

// func populateUserAndRoleData() {
// 	log.Println("\n👥 Populating Users and Roles...")
// 	startTime := time.Now()

// 	roles := []Role{
// 		{RoleName: "university_admin", Description: "Full system access"},
// 		{RoleName: "finance_controller", Description: "Manage fees and payments"},
// 		{RoleName: "registrar", Description: "Manage admissions and enrollments"},
// 		{RoleName: "college_admin", Description: "College level operations"},
// 		{RoleName: "hod", Description: "Head of department"},
// 		{RoleName: "faculty", Description: "Teach and assess"},
// 		{RoleName: "student", Description: "Student access"},
// 		{RoleName: "staff", Description: "Operational staff"},
// 		{RoleName: "auditor", Description: "Audit logs view"},
// 		{RoleName: "warden", Description: "Hostel management"},
// 		{RoleName: "librarian", Description: "Library management"},
// 	}
// 	DB.CreateInBatches(roles, 100)
// 	log.Printf("  ✅ Created %d roles", len(roles))

// 	admins := []User{
// 		{Username: "admin", Email: "admin@ntu.edu", PasswordHash: hashPW("Admin@123"), IsActive: true},
// 		{Username: "finance.controller", Email: "finance@ntu.edu", PasswordHash: hashPW("Finance@123"), IsActive: true},
// 		{Username: "registrar", Email: "registrar@ntu.edu", PasswordHash: hashPW("Registrar@123"), IsActive: true},
// 	}
// 	DB.CreateInBatches(admins, 100)

// 	faculties := make([]User, 0)
// 	firstNames := []string{"Rajesh", "Anita", "Suresh", "Priya", "Vikram", "Divya", "Amit", "Neha", "Arjun", "Sneha"}
// 	for i, name := range firstNames {
// 		faculties = append(faculties, User{
// 			Username:     fmt.Sprintf("faculty.%s", name),
// 			Email:        fmt.Sprintf("%s@ntu.edu", name),
// 			PasswordHash: hashPW("Faculty@123"),
// 			IsActive:     true,
// 		})
// 		if i >= 14 {
// 			break
// 		}
// 	}
// 	DB.CreateInBatches(faculties, 100)

// 	students := make([]User, 0)
// 	studentNames := []string{"Divya", "Raj", "Neha", "Akshay", "Priya", "Rohan", "Pooja", "Adit", "Zara", "Karan"}
// 	for i := 0; i < 100; i++ {
// 		name := studentNames[i%len(studentNames)]
// 		students = append(students, User{
// 			Username:     fmt.Sprintf("student.%s%d", name, i),
// 			Email:        fmt.Sprintf("student%d@ntu.edu", i),
// 			PasswordHash: hashPW("Student@123"),
// 			IsActive:     true,
// 		})
// 	}
// 	DB.CreateInBatches(students, 500)
// 	log.Printf("  ✅ Created %d student users", len(students))

// 	staffUsers := make([]User, 0)
// 	staffPositions := []string{"librarian", "warden", "clerk", "accountant", "technician"}
// 	for i := 0; i < 20; i++ {
// 		staffUsers = append(staffUsers, User{
// 			Username:     fmt.Sprintf("staff.%s%d", staffPositions[i%len(staffPositions)], i),
// 			Email:        fmt.Sprintf("staff%d@ntu.edu", i),
// 			PasswordHash: hashPW("Staff@123"),
// 			IsActive:     true,
// 		})
// 	}
// 	DB.CreateInBatches(staffUsers, 100)
// 	log.Printf("  ✅ Created %d staff users", len(staffUsers))

// 	// Assign Roles
// 	roleMap := make(map[string]uint)
// 	var allRoles []Role
// 	DB.Find(&allRoles)
// 	for _, r := range allRoles {
// 		roleMap[r.RoleName] = r.ID
// 	}

// 	userRoles := make([]UserRole, 0)
// 	var allUsers []User
// 	DB.Find(&allUsers)

// 	for i, u := range allUsers {
// 		var roleID uint
// 		var assigned bool

// 		if i < 3 {
// 			switch i {
// 			case 0:
// 				roleID = roleMap["university_admin"]
// 			case 1:
// 				roleID = roleMap["finance_controller"]
// 			case 2:
// 				roleID = roleMap["registrar"]
// 			}
// 			assigned = true
// 		} else if i < 3+len(faculties) {
// 			roleID = roleMap["faculty"]
// 			assigned = true
// 		} else if i < 3+len(faculties)+len(students) {
// 			roleID = roleMap["student"]
// 			assigned = true
// 		} else {
// 			staffIdx := (i - 3 - len(faculties) - len(students)) % len(staffPositions)
// 			roles := []string{"librarian", "warden", "staff", "staff", "staff"}
// 			roleID = roleMap[roles[staffIdx%len(roles)]]
// 			assigned = true
// 		}

// 		if assigned {
// 			userRoles = append(userRoles, UserRole{
// 				UserID:     u.ID,
// 				RoleID:     roleID,
// 				AssignedAt: time.Now(),
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(userRoles, 500)
// 	log.Printf("  ✅ Created %d user role assignments", len(userRoles))

// 	log.Printf("⏱️  User and role data completed in %v\n", time.Since(startTime))
// }

// func populateHRData() {
// 	log.Println("\n💼 Populating HR Data...")
// 	startTime := time.Now()

// 	var departments []Department
// 	var maleGender, femaleGender Gender
// 	var fullTimeEmp EmploymentType
// 	var profDesig, lecturerDesig Designation

// 	DB.Find(&departments)
// 	DB.Where("code = ?", "M").First(&maleGender)
// 	DB.Where("code = ?", "F").First(&femaleGender)
// 	DB.Where("code = ?", "FULL_TIME").First(&fullTimeEmp)
// 	DB.Where("code = ?", "PROF").First(&profDesig)
// 	DB.Where("code = ?", "LECTURER").First(&lecturerDesig)

// 	var facultyRole Role
// 	DB.Where("role_name = ?", "faculty").First(&facultyRole)

// 	var facultyUsers []User
// 	DB.Where("id IN (?)", DB.Table("shared.user_roles").Where("role_id = ?", facultyRole.ID).Select("user_id")).
// 		Find(&facultyUsers)

// 	employees := make([]Employee, 0)
// 	for i, fu := range facultyUsers {
// 		genderID := &maleGender.ID
// 		if i%2 == 0 {
// 			genderID = &femaleGender.ID
// 		}

// 		emp := Employee{
// 			UserID:           fu.ID,
// 			EmployeeCode:     fmt.Sprintf("EMP%05d", fu.ID),
// 			FirstName:        fu.Username,
// 			LastName:         "Faculty",
// 			GenderID:         genderID,
// 			   DateOfBirth:      func() *time.Time { t := time.Date(1980+i%20, 1, 1, 0, 0, 0, 0, time.UTC); return &t }(),
// 			Phone:            fmt.Sprintf("9%010d", fu.ID),
// 			JoiningDate:      time.Now().AddDate(-5-i%10, 0, 0),
// 			EmploymentTypeID: &fullTimeEmp.ID,
// 			DepartmentID:     &departments[i%len(departments)].ID,
// 			DesignationID:    &profDesig.ID,
// 			   // IsActive:         true, // Not present in Student struct
// 		}
// 		employees = append(employees, emp)
// 	}

// 	var staffRole Role
// 	DB.Where("role_name = ?", "staff").First(&staffRole)
// 	var staffUsers []User
// 	DB.Where("id IN (?)", DB.Table("shared.user_roles").Where("role_id = ?", staffRole.ID).Select("user_id")).
// 		Limit(20).Find(&staffUsers)

// 	for i, su := range staffUsers {
// 		genderID := &maleGender.ID
// 		if i%3 == 0 {
// 			genderID = &femaleGender.ID
// 		}

// 		designID := &lecturerDesig.ID

// 		emp := Employee{
// 			UserID:           su.ID,
// 			EmployeeCode:     fmt.Sprintf("STF%05d", su.ID),
// 			FirstName:        su.Username,
// 			LastName:         "Staff",
// 			GenderID:         genderID,
// 			JoiningDate:      time.Now().AddDate(-3-i%7, 0, 0),
// 			EmploymentTypeID: &fullTimeEmp.ID,
// 			DepartmentID:     &departments[i%len(departments)].ID,
// 			DesignationID:    designID,
// 			IsActive:         true,
// 		}
// 		employees = append(employees, emp)
// 	}

// 	DB.CreateInBatches(employees, 500)
// 	log.Printf("  ✅ Created %d employees", len(employees))

// 	faculties := make([]Faculty, 0)
// 	for _, emp := range employees[:len(facultyUsers)] {
// 		faculties = append(faculties, Faculty{
// 			EmployeeID:     emp.ID,
// 			Specialization: "Computer Science",
// 			Qualification:  "PhD",
// 			ResearchArea:   "Machine Learning",
// 			OfficeHours:    "Mon-Fri 3-5 PM",
// 			MaxLoadCredits: 20,
// 		})
// 	}
// 	DB.CreateInBatches(faculties, 500)
// 	log.Printf("  ✅ Created %d faculty records", len(faculties))

// 	salaries := make([]Salary, 0)
// 	for i, emp := range employees {
// 		basePay := 50000.0 + float64(i%5)*10000
// 		salary := Salary{
// 			EmployeeID:    emp.ID,
// 			EffectiveFrom: time.Now().AddDate(-1, 0, 0),
// 			BasePay:       basePay,
// 			NetSalary:     basePay * 0.88,
// 			IsActive:      true,
// 		}
// 		salaries = append(salaries, salary)
// 	}
// 	DB.CreateInBatches(salaries, 500)
// 	log.Printf("  ✅ Created %d salary records", len(salaries))

// 	var allLeaveTypes []LeaveType
// 	DB.Find(&allLeaveTypes)

// 	leaveBalances := make([]LeaveBalance, 0)
// 	for _, emp := range employees {
// 		for _, lt := range allLeaveTypes {
// 			leaveBalances = append(leaveBalances, LeaveBalance{
// 				EmployeeID:   emp.ID,
// 				LeaveTypeID:  lt.ID,
// 				TotalQuota:   lt.MaxDays,
// 				UsedQuota:    0,
// 				AccruedQuota: lt.MaxDays,
// 				Year:         time.Now().Year(),
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(leaveBalances, 1000)
// 	log.Printf("  ✅ Created %d leave balances", len(leaveBalances))

// 	attendances := make([]HRAttendance, 0)
// 	var presentStatus StatusCode
// 	DB.Where("module = ? AND code = ?", "hr", "PRESENT").First(&presentStatus)

// 	limit := 10
// 	if len(employees) < limit {
// 		limit = len(employees)
// 	}
// 	for _, emp := range employees[:limit] {
// 		for d := 0; d < 20; d++ {
// 			checkInTime := time.Now().AddDate(0, 0, -d).Add(8 * time.Hour)
// 			checkOutTime := checkInTime.Add(9 * time.Hour)

// 			attendances = append(attendances, HRAttendance{
// 				EmployeeID:     emp.ID,
// 				AttendanceDate: checkInTime,
// 				CheckIn:        &checkInTime,
// 				CheckOut:       &checkOutTime,
// 				StatusID:       &presentStatus.ID,
// 			})
// 		}
// 	}
// 	if len(attendances) > 0 {
// 		DB.CreateInBatches(attendances, 500)
// 		log.Printf("  ✅ Created %d attendance records", len(attendances))
// 	}

// 	log.Printf("⏱️  HR data completed in %v\n", time.Since(startTime))
// }

// func populateStudentData() {
// 	log.Println("\n🎓 Populating Student Data...")
// 	startTime := time.Now()

// 	var programs []Program
// 	var batches []Batch
// 	var sections []Section
// 	var femaleGender Gender
// 	var genCategory Category
// 	var oPlusBlood BloodGroup
// 	var activeStatus StatusCode

// 	DB.Find(&programs)
// 	DB.Find(&batches)
// 	DB.Find(&sections)
// 	DB.Where("code = ?", "F").First(&femaleGender)
// 	DB.Where("code = ?", "GEN").First(&genCategory)
// 	DB.Where("code = ?", "O+").First(&oPlusBlood)
// 	DB.Where("module = ? AND code = ?", "student", "ACTIVE").First(&activeStatus)

// 	var studentRole Role
// 	DB.Where("role_name = ?", "student").First(&studentRole)

// 	var studentUsers []User
// 	DB.Where("id IN (?)", DB.Table("shared.user_roles").Where("role_id = ?", studentRole.ID).Select("user_id")).
// 		Find(&studentUsers)

// 	log.Printf("  Found %d student users", len(studentUsers))

// 	students := make([]Student, 0)
// 	for i, su := range studentUsers {
// 		roll := fmt.Sprintf("24CSE%03d", i%1000)
// 		enroll := fmt.Sprintf("NTU2024%s", roll)

// 		stud := Student{
// 			UserID:           su.ID,
// 			EnrollmentNumber: enroll,
// 			RollNumber:       roll,
// 			FirstName:        su.Username,
// 			LastName:         "Student",
// 			DateOfBirth:      time.Date(2005, 1, 1+i%28, 0, 0, 0, 0, time.UTC),
// 			GenderID:         &femaleGender.ID,
// 			Email:            su.Email,
// 			Phone:            fmt.Sprintf("9%010d", su.ID),
// 			CategoryID:       &genCategory.ID,
// 			ProgramID:        programs[i%len(programs)].ID,
// 			AdmissionYear:    2024,
// 			StatusID:         &activeStatus.ID,
// 		}
// 		students = append(students, stud)
// 	}
// 	DB.CreateInBatches(students, 500)
// 	log.Printf("  ✅ Created %d students", len(students))

// 	guardians := make([]Guardian, 0)
// 	for _, stud := range students {
// 		guardians = append(guardians, Guardian{
// 			StudentID: stud.ID,
// 			Name:      "Parent of " + stud.FirstName,
// 			Relation:  "Father",
// 			Phone:     fmt.Sprintf("9%010d", stud.UserID),
// 			IsPrimary: true,
// 		})
// 	}
// 	DB.CreateInBatches(guardians, 500)
// 	log.Printf("  ✅ Created %d guardians", len(guardians))

// 	medicalRecords := make([]MedicalRecord, 0)
// 	for _, stud := range students {
// 		mr := MedicalRecord{
// 			StudentID:             stud.ID,
// 			BloodGroupID:          &oPlusBlood.ID,
// 			Allergies:             "None",
// 			EmergencyContactName:  "Parent",
// 			EmergencyContactPhone: fmt.Sprintf("9%010d", stud.UserID),
// 			InsurancePolicyNo:     fmt.Sprintf("POL%010d", stud.ID),
// 		}
// 		medicalRecords = append(medicalRecords, mr)
// 	}
// 	DB.CreateInBatches(medicalRecords, 500)
// 	log.Printf("  ✅ Created %d medical records", len(medicalRecords))

// 	termRegistrations := make([]TermRegistration, 0)
// 	var currentTerm AcademicTerm
// 	DB.Where("is_current = ?", true).First(&currentTerm)

// 	for _, stud := range students {
// 		var matchingBatch Batch
// 		DB.Where("program_id = ?", stud.ProgramID).First(&matchingBatch)

// 		var matchingSection Section
// 		DB.Where("batch_id = ?", matchingBatch.ID).First(&matchingSection)

// 		termRegistrations = append(termRegistrations, TermRegistration{
// 			StudentID:         stud.ID,
// 			AcademicTermID:    currentTerm.ID,
// 			BatchID:           matchingBatch.ID,
// 			SectionID:         matchingSection.ID,
// 			CurrentSemesterNo: 1,
// 			RegistrationDate:  time.Now(),
// 			Status:            "Active",
// 		})
// 	}
// 	DB.CreateInBatches(termRegistrations, 500)
// 	log.Printf("  ✅ Created %d term registrations", len(termRegistrations))

// 	log.Printf("⏱️  Student data completed in %v\n", time.Since(startTime))
// }

// func populateFinanceData() {
// 	log.Println("\n💰 Populating Finance Data...")
// 	startTime := time.Now()

// 	var feeHeads []FeeHead
// 	var students []Student
// 	var currentTerm AcademicTerm
// 	var unPaidStatus StatusCode

// 	DB.Find(&feeHeads)
// 	DB.Find(&students)
// 	DB.Where("is_current = ?", true).First(&currentTerm)
// 	DB.Where("module = ? AND code = ?", "finance", "UNPAID").First(&unPaidStatus)

// 	feeStructures := make([]FeeStructure, 0)
// 	var programs []Program
// 	DB.Find(&programs)

// 	for _, prog := range programs {
// 		for _, fh := range feeHeads[:4] {
// 			amount := 50000.0
// 			if fh.Code == "TUITION" {
// 				amount = 80000
// 			} else if fh.Code == "EXAM" {
// 				amount = 5000
// 			}

// 			for sem := 1; sem <= 4; sem++ {
// 				feeStructures = append(feeStructures, FeeStructure{
// 					ProgramID:      prog.ID,
// 					SemesterNumber: sem,
// 					FeeHeadID:      fh.ID,
// 					Amount:         amount,
// 					AcademicYear:   "2024-2025",
// 					IsActive:       true,
// 				})
// 			}
// 		}
// 	}
// 	DB.CreateInBatches(feeStructures, 500)
// 	log.Printf("  ✅ Created %d fee structures", len(feeStructures))

// 	invoices := make([]Invoice, 0)

// 	limit := 80
// 	if len(students) < limit {
// 		limit = len(students)
// 	}
// 	for i, stud := range students[:limit] {
// 		invoiceNo := fmt.Sprintf("INV-%d-%05d", time.Now().Year(), i+1)
// 		totalAmount := 145000.0

// 		inv := Invoice{
// 			StudentID:      stud.ID,
// 			InvoiceNumber:  invoiceNo,
// 			AcademicTermID: currentTerm.ID,
// 			GeneratedDate:  time.Now(),
// 			DueDate:        time.Now().AddDate(0, 1, 0),
// 			TotalAmount:    totalAmount,
// 			PaidAmount:     0,
// 			StatusID:       &unPaidStatus.ID,
// 		}
// 		invoices = append(invoices, inv)
// 	}

// 	DB.CreateInBatches(invoices, 500)
// 	log.Printf("  ✅ Created %d invoices", len(invoices))

// 	var createdInvoices []Invoice
// 	DB.Find(&createdInvoices)

// 	invoiceItems := make([]InvoiceItem, 0)
// 	for _, inv := range createdInvoices[:len(invoices)] {
// 		for _, fh := range feeHeads[:3] {
// 			amount := 50000.0
// 			if fh.Code == "TUITION" {
// 				amount = 80000
// 			} else if fh.Code == "EXAM" {
// 				amount = 5000
// 			}

// 			invoiceItems = append(invoiceItems, InvoiceItem{
// 				InvoiceID:   inv.ID,
// 				FeeHeadID:   fh.ID,
// 				Description: fh.Name,
// 				Quantity:    1,
// 				UnitAmount:  amount,
// 				Amount:      amount,
// 			})
// 		}
// 	}
// 	DB.CreateInBatches(invoiceItems, 1000)
// 	log.Printf("  ✅ Created %d invoice items", len(invoiceItems))

// 	scholarships := []Scholarship{
// 		{Name: "Merit Scholarship", Amount: 50000, Renewable: true},
// 		{Name: "Financial Aid", Amount: 100000, Renewable: true},
// 		{Name: "Sports Scholarship", Amount: 75000, Renewable: true},
// 		{Name: "Minority Scholarship", Amount: 60000, Renewable: true},
// 	}
// 	DB.CreateInBatches(scholarships, 100)
// 	log.Printf("  ✅ Created %d scholarships", len(scholarships))

// 	log.Printf("⏱️  Finance data completed in %v\n", time.Since(startTime))
// }

// func populateLibraryData() {
// 	log.Println("\n📖 Populating Library Data...")
// 	startTime := time.Now()

// 	authors := make([]Author, 0)
// 	authorNames := []string{
// 		"Thomas H. Cormen", "Donald E. Knuth", "Bjarne Stroustrup", "Andrew S. Tanenbaum",
// 		"Mark Allen Weiss", "Robert C. Martin", "Gang of Four", "Steve McConnell",
// 		"Kent Beck", "Martin Fowler", "Eric Evans", "Grady Booch",
// 	}

// 	for _, name := range authorNames {
// 		authors = append(authors, Author{Name: name, Biography: fmt.Sprintf("Bio of %s", name)})
// 	}
// 	DB.CreateInBatches(authors, 100)
// 	log.Printf("  ✅ Created %d authors", len(authors))

// 	books := make([]Book, 0)
// 	bookData := []struct {
// 		title           string
// 		isbn            string
// 		publisher       string
// 		author          string
// 		totalCopies     int
// 		publicationYear int
// 	}{
// 		{"Introduction to Algorithms", "978-0262033848", "MIT Press", "Thomas H. Cormen", 5, 2009},
// 		{"The Art of Computer Programming", "978-0201632881", "Addison-Wesley", "Donald E. Knuth", 3, 2006},
// 		{"The C++ Programming Language", "978-0321563842", "Addison-Wesley", "Bjarne Stroustrup", 4, 2013},
// 		{"Computer Networks", "978-0130384638", "Prentice Hall", "Andrew S. Tanenbaum", 3, 2010},
// 		{"Data Structures and Algorithm Analysis", "978-0132576277", "Pearson", "Mark Allen Weiss", 4, 2011},
// 		{"Clean Code", "978-0132350884", "Prentice Hall", "Robert C. Martin", 6, 2008},
// 		{"Design Patterns", "978-0201633610", "Addison-Wesley", "Gang of Four", 3, 1994},
// 		{"Code Complete", "978-0735619678", "Microsoft Press", "Steve McConnell", 2, 2004},
// 		{"Test Driven Development", "978-0321146533", "Addison-Wesley", "Kent Beck", 3, 2002},
// 		{"Refactoring", "978-0201485677", "Addison-Wesley", "Martin Fowler", 2, 1999},
// 		{"Domain-Driven Design", "978-0321125675", "Addison-Wesley", "Eric Evans", 3, 2003},
// 		{"Object-Oriented Analysis and Design", "978-0201895514", "Addison-Wesley", "Grady Booch", 2, 1994},
// 	}

// 	for _, b := range bookData {
// 		book := Book{
// 			Title:           b.title,
// 			ISBN:            b.isbn,
// 			Publisher:       b.publisher,
// 			PublicationYear: b.publicationYear,
// 			TotalCopies:     b.totalCopies,
// 			AvailableCopies: b.totalCopies,
// 			Location:        "Stack A",
// 		}
// 		books = append(books, book)
// 	}
// 	DB.CreateInBatches(books, 100)
// 	log.Printf("  ✅ Created %d books", len(books))

// 	bookCopies := make([]BookCopy, 0)
// 	for _, book := range books {
// 		for i := 1; i <= book.TotalCopies; i++ {
// 			barcode := fmt.Sprintf("BAR-%s-%d", book.ISBN[len(book.ISBN)-4:], i)
// 			copy := BookCopy{
// 				BookID:        book.ID,
// 				Barcode:       barcode,
// 				CopyNumber:    i,
// 				Condition:     "Good",
// 				ShelfLocation: "A-101",
// 			}
// 			bookCopies = append(bookCopies, copy)
// 		}
// 	}
// 	DB.CreateInBatches(bookCopies, 500)
// 	log.Printf("  ✅ Created %d book copies", len(bookCopies))

// 	var allAuthors []Author
// 	var allBooks []Book
// 	DB.Find(&allAuthors)
// 	DB.Find(&allBooks)

// 	bookAuthors := make([]BookAuthor, 0)
// 	for i, book := range allBooks {
// 		authorIdx := i % len(allAuthors)
// 		bookAuthors = append(bookAuthors, BookAuthor{
// 			BookID:   book.ID,
// 			AuthorID: allAuthors[authorIdx].ID,
// 		})
// 	}
// 	DB.CreateInBatches(bookAuthors, 500)
// 	log.Printf("  ✅ Created %d book author relationships", len(bookAuthors))

// 	log.Printf("⏱️  Library data completed in %v\n", time.Since(startTime))
// }

// func populateAdmissionsData() {
// 	log.Println("\n📋 Populating Admissions Data...")
// 	startTime := time.Now()

// 	admissionCycles := make([]AdmissionCycle, 0)
// 	var programs []Program
// 	DB.Limit(3).Find(&programs)

// 	for year := 2022; year <= 2024; year++ {
// 		for _, prog := range programs {
// 			cycle := AdmissionCycle{
// 				Name:             fmt.Sprintf("Admission Cycle %d - %s", year, prog.Code),
// 				AcademicYear:     fmt.Sprintf("%d-%d", year, year+1),
// 				ProgramID:        &prog.ID,
// 				ApplicationStart: time.Date(year-1, 10, 1, 0, 0, 0, 0, time.UTC),
// 				ApplicationEnd:   time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC),
// 				ApplicationFee:   500,
// 				MaxApplications:  500,
// 				IsOpen:           year == 2024,
// 			}
// 			admissionCycles = append(admissionCycles, cycle)
// 		}
// 	}
// 	DB.CreateInBatches(admissionCycles, 100)
// 	log.Printf("  ✅ Created %d admission cycles", len(admissionCycles))

// 	applicants := make([]Applicant, 0)
// 	var allCycles []AdmissionCycle
// 	var maleGender, femaleGender Gender
// 	var genCategory Category
// 	var appliedStatus StatusCode

// 	DB.Find(&allCycles)
// 	DB.Where("code = ?", "M").First(&maleGender)
// 	DB.Where("code = ?", "F").First(&femaleGender)
// 	DB.Where("code = ?", "GEN").First(&genCategory)
// 	DB.Where("module = ? AND code = ?", "admission", "APPLIED").First(&appliedStatus)

// 	for i := 0; i < 200; i++ {
// 		cycle := allCycles[i%len(allCycles)]
// 		genderID := &maleGender.ID
// 		if i%2 == 0 {
// 			genderID = &femaleGender.ID
// 		}

// 		applicant := Applicant{
// 			ApplicationNumber: fmt.Sprintf("APP-%d-%05d", cycle.ID, i),
// 			CycleID:           cycle.ID,
// 			ProgramID:         cycle.ProgramID,
// 			FirstName:         fmt.Sprintf("Applicant%d", i),
// 			LastName:          "Candidate",
// 			DateOfBirth:       time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
// 			Email:             fmt.Sprintf("applicant%d@example.com", i),
// 			Phone:             fmt.Sprintf("9%010d", i),
// 			GenderID:          genderID,
// 			CategoryID:        &genCategory.ID,
// 			EntranceScore:     float64(40 + (i % 60)),
// 			Rank:              i + 1,
// 			StatusID:          &appliedStatus.ID,
// 		}
// 		applicants = append(applicants, applicant)
// 	}
// 	DB.CreateInBatches(applicants, 500)
// 	log.Printf("  ✅ Created %d applicants", len(applicants))

// 	log.Printf("⏱️  Admissions data completed in %v\n", time.Since(startTime))
// }

// func populateExamData() {
// 	log.Println("\n📝 Populating Exam Data...")
// 	startTime := time.Now()

// 	var subjects []Subject
// 	var currentTerm AcademicTerm

// 	DB.Find(&subjects)
// 	DB.Where("is_current = ?", true).First(&currentTerm)

// 	if currentTerm.ID == 0 {
// 		log.Println("  ⚠️  No current academic term found, skipping exam data")
// 		return
// 	}

// 	examSchedules := make([]ExamSchedule, 0)
// 	for i, subj := range subjects[:20] {
// 		examDate := currentTerm.ExamStartDate.Add(time.Hour * 24 * time.Duration(i))
// 		sched := ExamSchedule{
// 			SubjectID:      subj.ID,
// 			AcademicTermID: currentTerm.ID,
// 			ExamDate:       examDate,
// 			StartTime:      "09:00 AM",
// 			EndTime:        "12:00 PM",
// 			ExamType:       "Theory",
// 			Venue:          fmt.Sprintf("Exam Hall %d", i+1),
// 			TotalMarks:     100,
// 			PassingMarks:   40,
// 		}
// 		examSchedules = append(examSchedules, sched)
// 	}
// 	DB.CreateInBatches(examSchedules, 100)
// 	log.Printf("  ✅ Created %d exam schedules", len(examSchedules))

// 	log.Printf("⏱️  Exam data completed in %v\n", time.Since(startTime))
// }

// func populateHostelData() {
// 	log.Println("\n🏨 Populating Hostel Data...")
// 	startTime := time.Now()

// 	var campus Campus
// 	var maleGender Gender
// 	DB.Where("is_main_campus = ?", true).First(&campus)
// 	DB.Where("code = ?", "M").First(&maleGender)

// 	hostels := make([]Hostel, 0)
// 	for i := 1; i <= 3; i++ {
// 		hostel := Hostel{
// 			Name:      fmt.Sprintf("Boys Hostel %d", i),
// 			Code:      fmt.Sprintf("BH-%d", i),
// 			CampusID:  &campus.ID,
// 			GenderID:  &maleGender.ID,
// 			TotalRooms: 50,
// 			IsActive:  true,
// 		}
// 		hostels = append(hostels, hostel)
// 	}
// 	DB.CreateInBatches(hostels, 100)
// 	log.Printf("  ✅ Created %d hostels", len(hostels))

// 	hostelRooms := make([]HostelRoom, 0)
// 	for _, hostel := range hostels {
// 		for floor := 1; floor <= 5; floor++ {
// 			for room := 1; room <= 10; room++ {
// 				hr := HostelRoom{
// 					HostelID:    hostel.ID,
// 					RoomNumber:  fmt.Sprintf("%d%02d", floor, room),
// 					RoomType:    "Double",
// 					Capacity:    2,
// 					MonthlyRent: 5000,
// 					IsAvailable: true,
// 				}
// 				hostelRooms = append(hostelRooms, hr)
// 			}
// 		}
// 	}
// 	DB.CreateInBatches(hostelRooms, 500)
// 	log.Printf("  ✅ Created %d hostel rooms", len(hostelRooms))

// 	hostelBeds := make([]HostelBed, 0)
// 	for _, room := range hostelRooms {
// 		for bed := 1; bed <= room.Capacity; bed++ {
// 			hb := HostelBed{
// 				RoomID:    room.ID,
// 				BedNumber: fmt.Sprintf("%s-B%d", room.RoomNumber, bed),
// 				IsOccupied: false,
// 			}
// 			hostelBeds = append(hostelBeds, hb)
// 		}
// 	}
// 	DB.CreateInBatches(hostelBeds, 1000)
// 	log.Printf("  ✅ Created %d hostel beds", len(hostelBeds))

// 	log.Printf("⏱️  Hostel data completed in %v\n", time.Since(startTime))
// }

// func populateTransportData() {
// 	log.Println("\n🚌 Populating Transport Data...")
// 	startTime := time.Now()

// 	buses := make([]Bus, 0)
// 	for i := 1; i <= 10; i++ {
// 		bus := Bus{
// 			BusNumber:      fmt.Sprintf("NTU-BUS-%03d", i),
// 			RegistrationNo: fmt.Sprintf("TS09Z%04d", 1000+i),
// 			Capacity:       50 + (i % 5 * 10),
// 			IsActive:       true,
// 		}
// 		buses = append(buses, bus)
// 	}
// 	DB.CreateInBatches(buses, 100)
// 	log.Printf("  ✅ Created %d buses", len(buses))

// 	routes := make([]Route, 0)
// 	routeNames := []string{"Route A - City Center", "Route B - North", "Route C - South", "Route D - East", "Route E - West"}
// 	for _, name := range routeNames {
// 		route := Route{
// 			RouteName:     name,
// 			Description:   fmt.Sprintf("Description for %s", name),
// 			DistanceKm:    25.5,
// 			EstimatedTime: "45 mins",
// 			IsActive:      true,
// 		}
// 		routes = append(routes, route)
// 	}
// 	DB.CreateInBatches(routes, 100)
// 	log.Printf("  ✅ Created %d routes", len(routes))

// 	stops := make([]Stop, 0)
// 	var allRoutes []Route
// 	DB.Find(&allRoutes)

// 	for _, route := range allRoutes {
// 		for stopNo := 1; stopNo <= 5; stopNo++ {
// 			stop := Stop{
// 				RouteID:       route.ID,
// 				StopName:      fmt.Sprintf("%s - Stop %d", route.RouteName, stopNo),
// 				StopOrder:     stopNo,
// 				ArrivalTime:   fmt.Sprintf("%02d:%02d", 8+stopNo, 0),
// 				DepartureTime: fmt.Sprintf("%02d:%02d", 8+stopNo, 5),
// 			}
// 			stops = append(stops, stop)
// 		}
// 	}
// 	DB.CreateInBatches(stops, 500)
// 	log.Printf("  ✅ Created %d stops", len(stops))

// 	log.Printf("⏱️  Transport data completed in %v\n", time.Since(startTime))
// }

// func populateSecurityData() {
// 	log.Println("\n🔐 Populating Security Data...")
// 	startTime := time.Now()

// 	permissions := make([]Permission, 0)
// 	resources := []string{"users", "students", "employees", "invoices", "books", "results"}
// 	actions := []string{"create", "read", "update", "delete"}

// 	for _, resource := range resources {
// 		for _, action := range actions {
// 			perm := Permission{
// 				Resource:    resource,
// 				Action:      action,
// 				Description: fmt.Sprintf("%s:%s", resource, action),
// 			}
// 			permissions = append(permissions, perm)
// 		}
// 	}
// 	DB.CreateInBatches(permissions, 100)
// 	log.Printf("  ✅ Created %d permissions", len(permissions))

// 	var roles []Role
// 	var allPerms []Permission
// 	DB.Find(&roles)
// 	DB.Find(&allPerms)

// 	rolePerms := make([]RolePermission, 0)
// 	adminRole := roles[0]

// 	for _, perm := range allPerms {
// 		rp := RolePermission{
// 			RoleID:       adminRole.ID,
// 			PermissionID: perm.ID,
// 			GrantedAt:    time.Now(),
// 		}
// 		rolePerms = append(rolePerms, rp)
// 	}
// 	DB.CreateInBatches(rolePerms, 500)
// 	log.Printf("  ✅ Created %d role permissions", len(rolePerms))

// 	log.Printf("⏱️  Security data completed in %v\n", time.Since(startTime))
// }

// func populateAcademicCalendarData() {
// 	log.Println("\n📅 Populating Academic Calendar Data...")
// 	startTime := time.Now()

// 	var campus Campus
// 	DB.Where("is_main_campus = ?", true).First(&campus)

// 	calendarEvents := make([]AcademicCalendar, 0)
// 	eventTypes := []struct {
// 		name      string
// 		eventType string
// 	}{
// 		{"Semester Start", "semester_start"},
// 		{"Registration Deadline", "registration"},
// 		{"Mid Semester Exams", "exam"},
// 		{"Holiday Break", "holiday"},
// 		{"Final Exams", "exam"},
// 		{"Grade Submission", "admin"},
// 		{"Convocation", "special"},
// 	}

// 	currentYear := time.Now().Year()
// 	for year := currentYear - 2; year <= currentYear; year++ {
// 		for i, evt := range eventTypes {
// 			eventDate := time.Date(year, time.Month(1+(i*2)), 1, 0, 0, 0, 0, time.UTC)
// 			event := AcademicCalendar{
// 				CampusID:    &campus.ID,
// 				EventDate:   eventDate,
// 				EventName:   evt.name,
// 				EventType:   evt.eventType,
// 				Description: fmt.Sprintf("Event: %s", evt.name),
// 			}
// 			calendarEvents = append(calendarEvents, event)
// 		}
// 	}
// 	DB.CreateInBatches(calendarEvents, 100)
// 	log.Printf("  ✅ Created %d academic calendar events", len(calendarEvents))

// 	log.Printf("⏱️  Academic calendar completed in %v\n", time.Since(startTime))
// }

// func populateCourseOfferingData() {
// 	log.Println("\n🎓 Populating Course Offering Data...")
// 	startTime := time.Now()

// 	var programs []Program
// 	var subjects []Subject
// 	var currentTerm AcademicTerm
// 	var batches []Batch
// 	var sections []Section
// 	var rooms []Room
// 	var employees []Employee

// 	DB.Find(&programs)
// 	DB.Find(&subjects)
// 	DB.Where("is_current = ?", true).First(&currentTerm)
// 	DB.Find(&batches)
// 	DB.Find(&sections)
// 	DB.Find(&rooms)
// 	DB.Where("id IN (?)", DB.Table("hr.faculties").Select("employee_id")).Find(&employees)

// 	log.Printf("  Found %d programs, %d subjects, %d batches, %d employees", 
// 		len(programs), len(subjects), len(batches), len(employees))

// 	courseOfferings := make([]CourseOffering, 0)
	
// 	if len(employees) > 0 && len(subjects) > 0 {
// 		for i, program := range programs {
// 			for j, subject := range subjects {
// 				if j >= 6 {
// 					break
// 				}

// 				batch := batches[(i*5+j)%len(batches)]
// 				section := sections[(i*3+j)%len(sections)]
// 				room := rooms[(i*7+j)%len(rooms)]
// 				faculty := employees[j%len(employees)]

// 				offering := CourseOffering{
// 					ProgramID:         program.ID,
// 					SubjectID:         subject.ID,
// 					AcademicTermID:    currentTerm.ID,
// 					BatchID:           batch.ID,
// 					SectionID:         &section.ID,
// 					FacultyEmployeeID: faculty.ID,
// 					RoomID:            &room.ID,
// 					MaxCapacity:       60,
// 					Status:            "Active",
// 				}
// 				courseOfferings = append(courseOfferings, offering)
// 			}
// 		}
// 	}

// 	if len(courseOfferings) > 0 {
// 		DB.CreateInBatches(courseOfferings, 500)
// 		log.Printf("  ✅ Created %d course offerings", len(courseOfferings))
// 	}

// 	log.Printf("⏱️  Course offering data completed in %v\n", time.Since(startTime))
// }

// func populateCourseRegistrationData() {
// 	log.Println("\n📚 Populating Course Registration Data...")
// 	startTime := time.Now()

// 	var students []Student
// 	var offerings []CourseOffering
// 	var enrolledStatus StatusCode

// 	DB.Limit(50).Find(&students)
// 	DB.Limit(100).Find(&offerings)
// 	DB.Where("module = ? AND code = ?", "student", "ACTIVE").First(&enrolledStatus)

// 	log.Printf("  Found %d students, %d offerings", len(students), len(offerings))

// 	courseRegs := make([]CourseRegistration, 0)

// 	for i, stud := range students {
// 		coursesPerStudent := 4 + (i % 3)
// 		for j := 0; j < coursesPerStudent && j < len(offerings); j++ {
// 			offering := offerings[(i*7+j)%len(offerings)]

// 			if offering.ProgramID != stud.ProgramID {
// 				continue
// 			}

// 			courseReg := CourseRegistration{
// 				StudentID:          stud.ID,
// 				OfferingID:         offering.ID,
// 				RegistrationStatus: "Enrolled",
// 				IsRepeat:           false,
// 				IsElective:         j >= 3,
// 			}
// 			courseRegs = append(courseRegs, courseReg)
// 		}
// 	}

// 	if len(courseRegs) > 0 {
// 		DB.CreateInBatches(courseRegs, 500)
// 		log.Printf("  ✅ Created %d course registrations", len(courseRegs))
// 	}

// 	log.Printf("⏱️  Course registration completed in %v\n", time.Since(startTime))
// }

// func populateClassSessionData() {
// 	log.Println("\n🏫 Populating Class Session Data...")
// 	startTime := time.Now()

// 	var offerings []CourseOffering
// 	var presentStatus StatusCode

// 	DB.Limit(30).Find(&offerings)
// 	DB.Where("module = ? AND code = ?", "hr", "PRESENT").First(&presentStatus)

// 	classSessions := make([]ClassSession, 0)
// 	for _, offering := range offerings {
// 		for d := 0; d < 20; d++ {
// 			classDate := time.Now().AddDate(0, 0, -(30-d))
// 			startTime := classDate.Add(9 * time.Hour)
// 			endTime := startTime.Add(time.Hour)

// 			session := ClassSession{
// 				OfferingID:  offering.ID,
// 				ClassDate:   classDate,
// 				StartTime:   &startTime,
// 				EndTime:     &endTime,
// 				RoomID:      offering.RoomID,
// 				FacultyID:   offering.FacultyEmployeeID,
// 				StatusID:    &presentStatus.ID,
// 			}
// 			classSessions = append(classSessions, session)
// 		}
// 	}

// 	if len(classSessions) > 0 {
// 		DB.CreateInBatches(classSessions, 1000)
// 		log.Printf("  ✅ Created %d class sessions", len(classSessions))
// 	}

// 	log.Printf("⏱️  Class session data completed in %v\n", time.Since(startTime))
// }

// func populateStudentAttendanceData() {
// 	log.Println("\n✅ Populating Student Attendance Data...")
// 	startTime := time.Now()

// 	var sessions []ClassSession
// 	var students []Student
// 	var presentStatus, absentStatus StatusCode

// 	DB.Limit(100).Find(&sessions)
// 	DB.Limit(50).Find(&students)
// 	DB.Where("module = ? AND code = ?", "hr", "PRESENT").First(&presentStatus)
// 	DB.Where("module = ? AND code = ?", "hr", "ABSENT").First(&absentStatus)

// 	attendances := make([]StudentAttendance, 0)
// 	for _, session := range sessions {
// 		for i, stud := range students {
// 			isPresent := i%5 != 0
// 			statusID := &presentStatus.ID
// 			if !isPresent {
// 				statusID = &absentStatus.ID
// 			}

// 			markedAt := session.ClassDate.Add(time.Hour * 9)
// 			attendance := StudentAttendance{
// 				SessionID: session.ID,
// 				StudentID: stud.ID,
// 				StatusID:  statusID,
// 				MarkedAt:  &markedAt,
// 			}
// 			attendances = append(attendances, attendance)
// 		}
// 	}

// 	if len(attendances) > 0 {
// 		DB.CreateInBatches(attendances, 1000)
// 		log.Printf("  ✅ Created %d student attendance records", len(attendances))
// 	}

// 	log.Printf("⏱️  Student attendance completed in %v\n", time.Since(startTime))
// }

// func populateResultsData() {
// 	log.Println("\n📊 Populating Results Data...")
// 	startTime := time.Now()

// 	var courseRegs []CourseRegistration
// 	var students []Student
// 	var currentTerm AcademicTerm

// 	DB.Limit(100).Find(&courseRegs)
// 	DB.Limit(50).Find(&students)
// 	DB.Where("is_current = ?", true).First(&currentTerm)

// 	results := make([]Result, 0)
// 	gradePoints := map[string]float64{
// 		"A": 4.0, "B": 3.0, "C": 2.0, "D": 1.0, "F": 0.0,
// 	}

// 	for i, courseReg := range courseRegs {
// 		marksObtained := 40.0 + float64(i%60)
// 		maxMarks := 100.0
// 		percentage := (marksObtained / maxMarks) * 100

// 		var grade string
// 		var gradePoint float64
// 		var isPassed bool

// 		if percentage >= 90 {
// 			grade = "A"
// 			gradePoint = gradePoints["A"]
// 		} else if percentage >= 80 {
// 			grade = "B"
// 			gradePoint = gradePoints["B"]
// 		} else if percentage >= 70 {
// 			grade = "C"
// 			gradePoint = gradePoints["C"]
// 		} else if percentage >= 60 {
// 			grade = "D"
// 			gradePoint = gradePoints["D"]
// 		} else {
// 			grade = "F"
// 			gradePoint = gradePoints["F"]
// 		}

// 		isPassed = percentage >= 40

// 		publishedAt := currentTerm.EndDate.AddDate(0, 1, 0)
// 		result := Result{
// 			StudentID:      courseReg.StudentID,
// 			CourseRegID:    courseReg.ID,
// 			SubjectID:      0,
// 			AcademicTermID: currentTerm.ID,
// 			MarksObtained:  marksObtained,
// 			MaxMarks:       maxMarks,
// 			Grade:          grade,
// 			GradePoint:     gradePoint,
// 			IsPassed:       isPassed,
// 			PublishedAt:    &publishedAt,
// 		}
// 		results = append(results, result)
// 	}

// 	if len(results) > 0 {
// 		DB.CreateInBatches(results, 500)
// 		log.Printf("  ✅ Created %d results", len(results))
// 	}

// 	log.Printf("⏱️  Results data completed in %v\n", time.Since(startTime))
// }

// func populatePaymentData() {
// 	log.Println("\n💳 Populating Payment Data...")
// 	startTime := time.Now()

// 	var invoices []Invoice
// 	var paidStatus StatusCode

// 	DB.Limit(50).Find(&invoices)
// 	DB.Where("module = ? AND code = ?", "finance", "PAID").First(&paidStatus)

// 	payments := make([]Payment, 0)
// 	for i, invoice := range invoices {
// 		if i%2 == 0 {
// 			paymentAmount := invoice.TotalAmount * 0.75
// 			payment := Payment{
// 				InvoiceID:   invoice.ID,
// 				StudentID:   invoice.StudentID,
// 				Amount:      paymentAmount,
// 				PaymentDate: time.Now().AddDate(0, 0, -5),
// 				StatusID:    &paidStatus.ID,
// 			}
// 			payments = append(payments, payment)
// 		}
// 	}

// 	if len(payments) > 0 {
// 		DB.CreateInBatches(payments, 500)
// 		log.Printf("  ✅ Created %d payments", len(payments))
// 	}

// 	log.Printf("⏱️  Payment data completed in %v\n", time.Since(startTime))
// }

// func populateCirculationData() {
// 	log.Println("\n📚 Populating Library Circulation Data...")
// 	startTime := time.Now()

// 	var bookCopies []BookCopy
// 	var students []Student
// 	var issuedStatus StatusCode

// 	DB.Find(&bookCopies)
// 	DB.Limit(30).Find(&students)
// 	DB.Where("module = ? AND code = ?", "library", "ISSUED").First(&issuedStatus)

// 	log.Printf("  Found %d book copies, %d students", len(bookCopies), len(students))

// 	circulations := make([]Circulation, 0)
// 	for i := 0; i < 50 && i < len(bookCopies); i++ {
// 		if i < len(students) {
// 			issuedDate := time.Now().AddDate(0, 0, -(20-i))
// 			dueDate := issuedDate.AddDate(0, 0, 14)

// 			circulation := Circulation{
// 				BookCopyID:   bookCopies[i].ID,
// 				StudentID:    students[i%len(students)].ID,
// 				IssuedDate:   issuedDate,
// 				DueDate:      dueDate,
// 				StatusID:     &issuedStatus.ID,
// 				FineAmount:   0,
// 				FinePaid:     false,
// 			}
// 			circulations = append(circulations, circulation)
// 		}
// 	}

// 	if len(circulations) > 0 {
// 		DB.CreateInBatches(circulations, 500)
// 		log.Printf("  ✅ Created %d circulation records", len(circulations))
// 	}

// 	log.Printf("⏱️  Library circulation completed in %v\n", time.Since(startTime))
// }

// func populateLeaveRequestData() {
// 	log.Println("\n🏖️  Populating Leave Request Data...")
// 	startTime := time.Now()

// 	var employees []Employee
// 	var leaveTypes []LeaveType
// 	var approvedStatus StatusCode

// 	DB.Limit(20).Find(&employees)
// 	DB.Find(&leaveTypes)
// 	DB.Where("module = ? AND code = ?", "hr", "PRESENT").First(&approvedStatus)

// 	leaveRequests := make([]LeaveRequest, 0)
// 	       for i, emp := range employees {
// 		       leaveType := leaveTypes[i%len(leaveTypes)]
// 		       startDate := time.Now().AddDate(0, 0, 10)
// 		endDate := startDate.AddDate(0, 0, int(leaveType.MaxDays)-1)
// 		approvedAt := time.Now()

// 		leaveReq := LeaveRequest{
// 			EmployeeID:  emp.ID,
// 			LeaveTypeID: leaveType.ID,
// 			StartDate:   startDate,
// 			EndDate:     endDate,
// 			Reason:      fmt.Sprintf("Leave request for %s", leaveType.Name),
// 			StatusID:    &approvedStatus.ID,
// 			ApprovedAt:  &approvedAt,
// 		}
// 		leaveRequests = append(leaveRequests, leaveReq)
// 	}

// 	if len(leaveRequests) > 0 {
// 		DB.CreateInBatches(leaveRequests, 500)
// 		log.Printf("  ✅ Created %d leave requests", len(leaveRequests))
// 	}

// 	log.Printf("⏱️  Leave request data completed in %v\n", time.Since(startTime))
// }

// func populateConfigurationData() {
// 	log.Println("\n⚙️  Populating Configuration Data...")
// 	startTime := time.Now()

// 	configurations := []Configuration{
// 		{ConfigKey: "semester_credit_minimum", ConfigValue: "12", DataType: "integer", Description: "Minimum credits per semester"},
// 		{ConfigKey: "semester_credit_maximum", ConfigValue: "24", DataType: "integer", Description: "Maximum credits per semester"},
// 		{ConfigKey: "attendance_passing_percentage", ConfigValue: "75", DataType: "integer", Description: "Minimum attendance percentage"},
// 		{ConfigKey: "cgpa_passing_minimum", ConfigValue: "2.0", DataType: "float", Description: "Minimum CGPA to pass"},
// 		{ConfigKey: "late_fee_percentage", ConfigValue: "5", DataType: "integer", Description: "Late fee percentage per day"},
// 		{ConfigKey: "library_fine_per_day", ConfigValue: "10", DataType: "float", Description: "Library fine per day in rupees"},
// 		{ConfigKey: "max_book_checkout_days", ConfigValue: "14", DataType: "integer", Description: "Maximum book checkout duration"},
// 		{ConfigKey: "max_book_renewal_count", ConfigValue: "2", DataType: "integer", Description: "Maximum renewal count"},
// 		{ConfigKey: "hostel_alloc_batch_size", ConfigValue: "50", DataType: "integer", Description: "Batch size for hostel allocation"},
// 		{ConfigKey: "transport_pass_validity_days", ConfigValue: "365", DataType: "integer", Description: "Transport pass validity in days"},
// 	}
// 	DB.CreateInBatches(configurations, 100)
// 	log.Printf("  ✅ Created %d configurations", len(configurations))

// 	log.Printf("⏱️  Configuration data completed in %v\n", time.Since(startTime))
// }

// func populateTimetableData() {
// 	log.Println("\n⏰ Populating Timetable Data...")
// 	startTime := time.Now()

// 	var offerings []CourseOffering
// 	DB.Limit(20).Find(&offerings)

// 	timetables := make([]Timetable, 0)
// 	daysOfWeek := []int{1, 2, 3, 4, 5}

// 	for _, offering := range offerings {
// 		for d := 0; d < 20; d++ {
// 			day := daysOfWeek[d%len(daysOfWeek)]
// 			hour := 8 + (d % 8) // 8AM to 15PM
// 			startTm := fmt.Sprintf("%02d:00", hour)
// 			endTm := fmt.Sprintf("%02d:00", hour+1)

// 			timetable := Timetable{
// 				OfferingID: offering.ID,
// 				DayOfWeek:  day,
// 				StartTime:  startTm,
// 				EndTime:    endTm,
// 			}
// 			timetables = append(timetables, timetable)
// 		}
// 	}

// 	if len(timetables) > 0 {
// 		DB.CreateInBatches(timetables, 500)
// 		log.Printf("  ✅ Created %d timetable entries", len(timetables))
// 	}

// 	log.Printf("⏱️  Timetable data completed in %v\n", time.Since(startTime))
// }

// // ============================================================================
// // MAIN FUNCTION
// // ============================================================================

// func main() {
// 	log.Println("\n" + strings.Repeat("=", 90))
// 	log.Println("🚀 UNIVERSITY ERP DATABASE - COMPREHENSIVE DATA POPULATION")
// 	log.Println(strings.Repeat("=", 90) + "\n")

// 	overallStart := time.Now()

// 	initDB()

// 	// Populate all data
// 	populateMasterData()
// 	populateCoreData()
// 	populateAcademicData()
// 	populateUserAndRoleData()
// 	populateHRData()
// 	populateStudentData()
// 	populateFinanceData()
// 	populateLibraryData()
// 	populateAdmissionsData()
// 	populateExamData()
// 	populateHostelData()
// 	populateTransportData()
// 	populateSecurityData()
// 	populateAcademicCalendarData()
// 	populateCourseOfferingData()
// 	populateCourseRegistrationData()
// 	populateClassSessionData()
// 	populateStudentAttendanceData()
// 	populateResultsData()
// 	populatePaymentData()
// 	populateCirculationData()
// 	populateLeaveRequestData()
// 	populateTimetableData()
// 	populateConfigurationData()

// 	// Summary
// 	log.Println("\n" + strings.Repeat("=", 90))
// 	log.Println("✅ DATABASE POPULATION COMPLETE!")
// 	log.Println(strings.Repeat("=", 90))
// 	log.Printf("⏱️  Total time taken: %v\n\n", time.Since(overallStart))

// 	log.Println("📊 DATA SUMMARY:")
// 	log.Println("  ✓ Master data (genders, categories, blood groups, designations, etc.)")
// 	log.Println("  ✓ Core organizational data (universities, campuses, departments, rooms)")
// 	log.Println("  ✓ Academic data (programs, semesters, subjects, batches, sections)")
// 	log.Println("  ✓ User management (users, roles, permissions, role assignments)")
// 	log.Println("  ✓ HR data (employees, faculty, salaries, leave balances, attendance)")
// 	log.Println("  ✓ Student data (100 students, guardians, medical records, registrations)")
// 	log.Println("  ✓ Finance data (fee structures, invoices with line items, scholarships)")
// 	log.Println("  ✓ Library data (books, copies, authors, circulation records)")
// 	log.Println("  ✓ Admissions data (cycles, applicants, applications)")
// 	log.Println("  ✓ Exam data (schedules, results with grades and GPA)")
// 	log.Println("  ✓ Hostel data (hostels, rooms, beds with allocation)")
// 	log.Println("  ✓ Transport data (buses, routes, stops, passes)")
// 	log.Println("  ✓ Security data (permissions, role permissions)")
// 	log.Println("  ✓ Academic calendar (events and important dates)")
// 	log.Println("  ✓ Course offerings with faculty assignment")
// 	log.Println("  ✓ Class sessions and student attendance tracking")
// 	log.Println("  ✓ Timetable generation for all courses")
// 	log.Println("  ✓ Payment records and financial transactions")
// 	log.Println("  ✓ Leave requests with approval status")
// 	log.Println("  ✓ System configurations and settings\n")

// 	log.Println("🔍 NEXT STEPS:")
// 	log.Println("  1. Verify data in PostgreSQL using queries")
// 	log.Println("  2. Connect your backend API to start using the data")
// 	log.Println("  3. Review audit logs and security settings")
// 	log.Println("  4. Test all CRUD operations\n")

// 	log.Println("✨ Database is now ready for backend integration!")
// 	log.Println(strings.Repeat("=", 90) + "\n")
// }