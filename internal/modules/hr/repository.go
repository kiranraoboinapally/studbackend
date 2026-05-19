package hrmod

import (
	"fmt"
	"time"

	"university-erp-backend/internal/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Lookups
func (r *Repository) ListDesignations() ([]domain.Designation, error) {
	var list []domain.Designation
	return list, r.db.Where("is_active = true").Find(&list).Error
}
func (r *Repository) ListEmploymentTypes() ([]domain.EmploymentType, error) {
	var list []domain.EmploymentType
	return list, r.db.Find(&list).Error
}
func (r *Repository) ListLeaveTypes() ([]domain.LeaveType, error) {
	var list []domain.LeaveType
	return list, r.db.Where("is_active = true").Find(&list).Error
}
func (r *Repository) CreateDesignation(d *domain.Designation) error {
	return r.db.Create(d).Error
}
func (r *Repository) CreateLeaveType(lt *domain.LeaveType) error {
	return r.db.Create(lt).Error
}

// Employees
func (r *Repository) ListEmployees(departmentID uint, page, pageSize int) ([]domain.Employee, int64, error) {
	var list []domain.Employee
	var total int64
	q := r.db.Model(&domain.Employee{}).Where("is_active = true")
	if departmentID > 0 {
		q = q.Where("department_id = ?", departmentID)
	}
	q.Count(&total)
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	return list, total, q.Offset((page-1)*pageSize).Limit(pageSize).Order("first_name").Find(&list).Error
}
func (r *Repository) GetEmployee(id uint) (*domain.Employee, error) {
	var e domain.Employee
	return &e, r.db.First(&e, id).Error
}
func (r *Repository) GetEmployeeByUserID(userID uint) (*domain.Employee, error) {
	var e domain.Employee
	return &e, r.db.Where("user_id = ?", userID).First(&e).Error
}
func (r *Repository) CreateEmployee(e *domain.Employee) error {
	return r.db.Create(e).Error
}
func (r *Repository) UpdateEmployee(e *domain.Employee) error {
	return r.db.Save(e).Error
}
func (r *Repository) DeactivateEmployee(id uint) error {
	return r.db.Model(&domain.Employee{}).Where("id = ?", id).Update("is_active", false).Error
}
func (r *Repository) CountEmployees() (int64, error) {
	var c int64
	r.db.Model(&domain.Employee{}).Where("is_active = true").Count(&c)
	return c, nil
}
func (r *Repository) GenerateEmployeeCode() (string, error) {
	var c int64
	r.db.Model(&domain.Employee{}).Count(&c)
	return fmt.Sprintf("EMP%05d", c+1), nil
}

// Faculty
func (r *Repository) GetFaculty(employeeID uint) (*domain.Faculty, error) {
	var f domain.Faculty
	return &f, r.db.Where("employee_id = ?", employeeID).First(&f).Error
}
func (r *Repository) UpsertFaculty(f *domain.Faculty) error {
	return r.db.Save(f).Error
}
func (r *Repository) ListFaculty(departmentID uint) ([]FacultyDetail, error) {
	var list []FacultyDetail
	q := `SELECT e.*, f.specialization, f.qualification, f.research_area, f.office_hours, f.max_load_credits
		FROM hr.employees e JOIN hr.faculties f ON f.employee_id = e.id WHERE e.is_active = true`
	if departmentID > 0 {
		q += fmt.Sprintf(" AND e.department_id = %d", departmentID)
	}
	return list, r.db.Raw(q).Scan(&list).Error
}

// Staff
func (r *Repository) GetStaff(employeeID uint) (*domain.Staff, error) {
	var s domain.Staff
	return &s, r.db.Where("employee_id = ?", employeeID).First(&s).Error
}
func (r *Repository) UpsertStaff(s *domain.Staff) error {
	return r.db.Save(s).Error
}

// Department History
func (r *Repository) AddDeptHistory(h *domain.EmployeeDepartmentHistory) error {
	return r.db.Create(h).Error
}
func (r *Repository) GetDeptHistory(employeeID uint) ([]domain.EmployeeDepartmentHistory, error) {
	var list []domain.EmployeeDepartmentHistory
	return list, r.db.Where("employee_id = ?", employeeID).Order("effective_from DESC").Find(&list).Error
}

// Salary
func (r *Repository) GetCurrentSalary(employeeID uint) (*domain.Salary, error) {
	var s domain.Salary
	return &s, r.db.Where("employee_id = ? AND is_active = true", employeeID).First(&s).Error
}
func (r *Repository) CreateSalary(s *domain.Salary) error {
	// Deactivate old salary
	r.db.Model(&domain.Salary{}).Where("employee_id = ? AND is_active = true", s.EmployeeID).Update("is_active", false)
	s.IsActive = true
	return r.db.Create(s).Error
}
func (r *Repository) GetSalaryDetails(salaryID uint) ([]domain.SalaryDetail, error) {
	var list []domain.SalaryDetail
	return list, r.db.Where("salary_id = ?", salaryID).Find(&list).Error
}
func (r *Repository) ListSalaryComponents() ([]domain.SalaryComponent, error) {
	var list []domain.SalaryComponent
	return list, r.db.Where("is_active = true").Find(&list).Error
}
func (r *Repository) CreateSalaryComponent(sc *domain.SalaryComponent) error {
	return r.db.Create(sc).Error
}

// Payroll
func (r *Repository) CreatePayrollRun(p *domain.PayrollRun) error {
	return r.db.Create(p).Error
}
func (r *Repository) ListPayrollRuns(employeeID uint) ([]domain.PayrollRun, error) {
	var list []domain.PayrollRun
	q := r.db.Order("month DESC")
	if employeeID > 0 {
		q = q.Where("employee_id = ?", employeeID)
	}
	return list, q.Find(&list).Error
}
func (r *Repository) GetPayrollRun(id uint) (*domain.PayrollRun, error) {
	var p domain.PayrollRun
	return &p, r.db.First(&p, id).Error
}

// Leave
func (r *Repository) GetLeaveBalance(employeeID, leaveTypeID, year uint) (*domain.LeaveBalance, error) {
	var lb domain.LeaveBalance
	return &lb, r.db.Where("employee_id = ? AND leave_type_id = ? AND year = ?", employeeID, leaveTypeID, year).First(&lb).Error
}
func (r *Repository) ListLeaveBalances(employeeID, year uint) ([]domain.LeaveBalance, error) {
	var list []domain.LeaveBalance
	return list, r.db.Where("employee_id = ? AND year = ?", employeeID, year).Find(&list).Error
}
func (r *Repository) UpsertLeaveBalance(lb *domain.LeaveBalance) error {
	return r.db.Save(lb).Error
}
func (r *Repository) CreateLeaveRequest(req *domain.LeaveRequest) error {
	return r.db.Create(req).Error
}
func (r *Repository) GetLeaveRequest(id uint) (*domain.LeaveRequest, error) {
	var req domain.LeaveRequest
	return &req, r.db.First(&req, id).Error
}
func (r *Repository) ListLeaveRequests(employeeID uint, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	var list []domain.LeaveRequest
	var total int64
	q := r.db.Model(&domain.LeaveRequest{})
	if employeeID > 0 {
		q = q.Where("employee_id = ?", employeeID)
	}
	q.Count(&total)
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	return list, total, q.Offset((page-1)*pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error
}
func (r *Repository) UpdateLeaveRequest(req *domain.LeaveRequest) error {
	return r.db.Save(req).Error
}

// Attendance (HR)
func (r *Repository) MarkHRAttendance(a *domain.HRAttendance) error {
	return r.db.Save(a).Error
}
func (r *Repository) GetHRAttendance(employeeID uint, from, to time.Time) ([]domain.HRAttendance, error) {
	var list []domain.HRAttendance
	return list, r.db.Where("employee_id = ? AND attendance_date BETWEEN ? AND ?", employeeID, from, to).Order("attendance_date DESC").Find(&list).Error
}

// Recruitment
func (r *Repository) ListJobs() ([]domain.RecruitmentJob, error) {
	var list []domain.RecruitmentJob
	return list, r.db.Order("posted_date DESC").Find(&list).Error
}
func (r *Repository) GetJob(id uint) (*domain.RecruitmentJob, error) {
	var j domain.RecruitmentJob
	return &j, r.db.First(&j, id).Error
}
func (r *Repository) CreateJob(j *domain.RecruitmentJob) error {
	return r.db.Create(j).Error
}
func (r *Repository) ListJobApplications(jobID uint) ([]domain.JobApplication, error) {
	var list []domain.JobApplication
	return list, r.db.Where("job_id = ?", jobID).Order("applied_at DESC").Find(&list).Error
}
func (r *Repository) CreateJobApplication(ja *domain.JobApplication) error {
	return r.db.Create(ja).Error
}
func (r *Repository) UpdateJobApplicationStatus(id, statusID uint) error {
	return r.db.Model(&domain.JobApplication{}).Where("id = ?", id).Update("status_id", statusID).Error
}

type FacultyDetail struct {
	domain.Employee
	Specialization string `json:"specialization"`
	Qualification  string `json:"qualification"`
	ResearchArea   string `json:"research_area"`
	OfficeHours    string `json:"office_hours"`
	MaxLoadCredits int    `json:"max_load_credits"`
}
