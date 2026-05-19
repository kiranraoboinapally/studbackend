package domain

import "time"

// ─── HR: Designations & Employment Types ─────────────────────────────────────

type Designation struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Code     string `gorm:"unique;not null" json:"code"`
	Name     string `gorm:"not null" json:"name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (Designation) TableName() string { return "hr.designations" }

type EmploymentType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}

func (EmploymentType) TableName() string { return "hr.employment_types" }

type LeaveType struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Code     string  `gorm:"unique;not null" json:"code"`
	Name     string  `gorm:"not null" json:"name"`
	MaxDays  float64 `json:"max_days"`
	Paid     bool    `json:"paid"`
	IsActive bool    `gorm:"default:true" json:"is_active"`
}

func (LeaveType) TableName() string { return "hr.leave_types" }

// ─── HR: Employees ───────────────────────────────────────────────────────────

type Employee struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"unique;not null;index" json:"user_id"`
	EmployeeCode     string     `gorm:"unique;not null;index" json:"employee_code"`
	FirstName        string     `gorm:"not null" json:"first_name"`
	LastName         string     `gorm:"not null" json:"last_name"`
	GenderID         *uint      `gorm:"index" json:"gender_id,omitempty"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Phone            string     `json:"phone"`
	Address          string     `json:"address"`
	JoiningDate      time.Time  `gorm:"not null;index" json:"joining_date"`
	EmploymentTypeID *uint      `gorm:"index" json:"employment_type_id,omitempty"`
	DepartmentID     *uint      `gorm:"index" json:"department_id,omitempty"`
	DesignationID    *uint      `gorm:"index" json:"designation_id,omitempty"`
	IsActive         bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (Employee) TableName() string { return "hr.employees" }

type EmployeeDepartmentHistory struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	EmployeeID    uint       `gorm:"not null;index" json:"employee_id"`
	DepartmentID  uint       `gorm:"not null" json:"department_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (EmployeeDepartmentHistory) TableName() string { return "hr.employee_department_history" }

type EmployeeDesignationHistory struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	EmployeeID    uint       `gorm:"not null;index" json:"employee_id"`
	DesignationID uint       `gorm:"not null" json:"designation_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (EmployeeDesignationHistory) TableName() string { return "hr.employee_designation_history" }

// ─── HR: Faculty & Staff ─────────────────────────────────────────────────────

type Faculty struct {
	EmployeeID     uint   `gorm:"primaryKey" json:"employee_id"`
	Specialization string `json:"specialization"`
	Qualification  string `json:"qualification"`
	ResearchArea   string `json:"research_area"`
	OfficeHours    string `json:"office_hours"`
	MaxLoadCredits int    `gorm:"default:20" json:"max_load_credits"`
}

func (Faculty) TableName() string { return "hr.faculties" }

type Staff struct {
	EmployeeID   uint   `gorm:"primaryKey" json:"employee_id"`
	JobTitle     string `json:"job_title"`
	SupervisorID *uint  `gorm:"index" json:"supervisor_id,omitempty"`
	WorkLocation string `json:"work_location"`
}

func (Staff) TableName() string { return "hr.staffs" }

// ─── HR: Salary & Payroll ────────────────────────────────────────────────────

type SalaryComponent struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Code     string `gorm:"unique;not null" json:"code"`
	Name     string `gorm:"not null" json:"name"`
	Type     string `gorm:"type:varchar(20)" json:"type"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (SalaryComponent) TableName() string { return "hr.salary_components" }

type Salary struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	EmployeeID    uint       `gorm:"not null;index" json:"employee_id"`
	EffectiveFrom time.Time  `gorm:"not null;index" json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	BasePay       float64    `json:"base_pay"`
	NetSalary     float64    `json:"net_salary"`
	IsActive      bool       `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Salary) TableName() string { return "hr.salaries" }

type SalaryDetail struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	SalaryID          uint    `gorm:"not null;index" json:"salary_id"`
	SalaryComponentID uint    `gorm:"not null" json:"salary_component_id"`
	Amount            float64 `json:"amount"`
}

func (SalaryDetail) TableName() string { return "hr.salary_details" }

type PayrollRun struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	EmployeeID      uint      `gorm:"not null;index" json:"employee_id"`
	Month           time.Time `gorm:"not null;index" json:"month"`
	GrossPay        float64   `json:"gross_pay"`
	TotalDeductions float64   `json:"total_deductions"`
	NetPay          float64   `json:"net_pay"`
	ProcessedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"processed_at"`
	ProcessedBy     *uint     `json:"processed_by,omitempty"`
}

func (PayrollRun) TableName() string { return "hr.payroll_runs" }

// ─── HR: Leave Management ────────────────────────────────────────────────────

type LeaveBalance struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   uint      `gorm:"not null;index" json:"employee_id"`
	LeaveTypeID  uint      `gorm:"not null;index" json:"leave_type_id"`
	TotalQuota   float64   `gorm:"not null" json:"total_quota"`
	UsedQuota    float64   `gorm:"default:0" json:"used_quota"`
	AccruedQuota float64   `gorm:"default:0" json:"accrued_quota"`
	Year         int       `gorm:"not null;index" json:"year"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (LeaveBalance) TableName() string { return "hr.leave_balances" }

type LeaveRequest struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EmployeeID  uint       `gorm:"not null;index" json:"employee_id"`
	LeaveTypeID uint       `gorm:"not null" json:"leave_type_id"`
	StartDate   time.Time  `gorm:"not null;index" json:"start_date"`
	EndDate     time.Time  `gorm:"not null" json:"end_date"`
	Reason      string     `json:"reason"`
	StatusID    *uint      `gorm:"index" json:"status_id,omitempty"`
	ApprovedBy  *uint      `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (LeaveRequest) TableName() string { return "hr.leave_requests" }

type HRAttendance struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	EmployeeID     uint       `gorm:"not null;index" json:"employee_id"`
	AttendanceDate time.Time  `gorm:"not null;index" json:"attendance_date"`
	CheckIn        *time.Time `json:"check_in,omitempty"`
	CheckOut       *time.Time `json:"check_out,omitempty"`
	StatusID       *uint      `gorm:"index" json:"status_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (HRAttendance) TableName() string { return "hr.attendance" }

// ─── HR: Recruitment ─────────────────────────────────────────────────────────

type RecruitmentJob struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Title            string     `gorm:"not null" json:"title"`
	DepartmentID     *uint      `gorm:"index" json:"department_id,omitempty"`
	EmploymentTypeID *uint      `gorm:"index" json:"employment_type_id,omitempty"`
	Vacancies        int        `json:"vacancies"`
	PostedDate       time.Time  `gorm:"default:CURRENT_DATE;index" json:"posted_date"`
	ClosingDate      *time.Time `json:"closing_date,omitempty"`
	Description      string     `json:"description"`
	StatusID         *uint      `gorm:"index" json:"status_id,omitempty"`
}

func (RecruitmentJob) TableName() string { return "hr.recruitment_jobs" }

type JobApplication struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	JobID         uint      `gorm:"not null;index" json:"job_id"`
	ApplicantName string    `gorm:"not null" json:"applicant_name"`
	Email         string    `gorm:"not null;index" json:"email"`
	Phone         string    `json:"phone"`
	ResumePath    string    `json:"resume_path"`
	AppliedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"applied_at"`
	StatusID      *uint     `gorm:"index" json:"status_id,omitempty"`
}

func (JobApplication) TableName() string { return "hr.job_applications" }
