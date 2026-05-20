package hrmod

import (
	"context"
	"fmt"
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
	return &Service{repo: repo, bus: bus, outbox: ob, db: db}
}

// Lookups
func (s *Service) ListDesignations(ctx context.Context) ([]domain.Designation, error) {
	return s.repo.ListDesignations()
}
func (s *Service) ListEmploymentTypes(ctx context.Context) ([]domain.EmploymentType, error) {
	return s.repo.ListEmploymentTypes()
}
func (s *Service) ListLeaveTypes(ctx context.Context) ([]domain.LeaveType, error) {
	return s.repo.ListLeaveTypes()
}
func (s *Service) CreateDesignation(ctx context.Context, d *domain.Designation) error {
	return s.repo.CreateDesignation(d)
}
func (s *Service) CreateLeaveType(ctx context.Context, lt *domain.LeaveType) error {
	return s.repo.CreateLeaveType(lt)
}

// Employees - with outbox event on creation
func (s *Service) ListEmployees(ctx context.Context, departmentID uint, page, pageSize int) ([]domain.Employee, int64, error) {
	return s.repo.ListEmployees(departmentID, page, pageSize)
}
func (s *Service) GetEmployee(ctx context.Context, id uint) (*domain.Employee, error) {
	e, err := s.repo.GetEmployee(id)
	if err != nil {
		return nil, apperrors.NotFound("employee not found")
	}
	return e, nil
}
func (s *Service) GetMyProfile(ctx context.Context, userID uint) (*domain.Employee, error) {
	e, err := s.repo.GetEmployeeByUserID(userID)
	if err != nil {
		return nil, apperrors.NotFound("employee profile not found")
	}
	return e, nil
}
func (s *Service) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	if e.FirstName == "" || e.LastName == "" {
		return apperrors.BadRequest("first and last name are required")
	}
	code, err := s.repo.GenerateEmployeeCode()
	if err != nil {
		return err
	}
	e.EmployeeCode = code
	if e.JoiningDate.IsZero() {
		e.JoiningDate = time.Now()
	}
	e.IsActive = true

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		// Create department history if department assigned
		if e.DepartmentID != nil {
			hist := &domain.EmployeeDepartmentHistory{
				EmployeeID:    e.ID,
				DepartmentID:  *e.DepartmentID,
				EffectiveFrom: e.JoiningDate,
			}
			if err := tx.Create(hist).Error; err != nil {
				return err
			}
		}
		designation := ""
		if e.DesignationID != nil {
			var d domain.Designation
			tx.First(&d, *e.DesignationID)
			designation = d.Name
		}
		return s.outbox.WriteEvent(tx, "Employee", fmt.Sprintf("%d", e.ID),
			eventbus.EventEmployeeOnboarded,
			eventbus.EmployeeOnboardedPayload{
				EmployeeID:   e.ID,
				UserID:       e.UserID,
				DepartmentID: e.DepartmentID,
				Designation:  designation,
			},
		)
	})
	return txErr
}
func (s *Service) UpdateEmployee(ctx context.Context, id uint, e *domain.Employee) error {
	existing, err := s.repo.GetEmployee(id)
	if err != nil {
		return apperrors.NotFound("employee not found")
	}
	e.ID = existing.ID
	e.EmployeeCode = existing.EmployeeCode
	return s.repo.UpdateEmployee(e)
}
func (s *Service) DeactivateEmployee(ctx context.Context, id uint) error {
	return s.repo.DeactivateEmployee(id)
}

// Faculty
func (s *Service) ListFaculty(ctx context.Context, departmentID uint) ([]FacultyDetail, error) {
	return s.repo.ListFaculty(departmentID)
}
func (s *Service) GetFacultyProfile(ctx context.Context, employeeID uint) (*domain.Faculty, error) {
	return s.repo.GetFaculty(employeeID)
}
func (s *Service) UpsertFacultyProfile(ctx context.Context, f *domain.Faculty) error {
	return s.repo.UpsertFaculty(f)
}

// Staff
func (s *Service) GetStaffProfile(ctx context.Context, employeeID uint) (*domain.Staff, error) {
	return s.repo.GetStaff(employeeID)
}
func (s *Service) UpsertStaffProfile(ctx context.Context, st *domain.Staff) error {
	return s.repo.UpsertStaff(st)
}

// Department History
func (s *Service) TransferDepartment(ctx context.Context, employeeID, deptID uint) error {
	s.db.Exec(`UPDATE hr.employee_department_history SET effective_to = ? WHERE employee_id = ? AND effective_to IS NULL`, time.Now(), employeeID)
	s.db.Model(&domain.Employee{}).Where("id = ?", employeeID).Update("department_id", deptID)
	hist := &domain.EmployeeDepartmentHistory{
		EmployeeID:    employeeID,
		DepartmentID:  deptID,
		EffectiveFrom: time.Now(),
	}
	return s.repo.AddDeptHistory(hist)
}
func (s *Service) GetDeptHistory(ctx context.Context, employeeID uint) ([]domain.EmployeeDepartmentHistory, error) {
	return s.repo.GetDeptHistory(employeeID)
}

// Salary
func (s *Service) GetCurrentSalary(ctx context.Context, employeeID uint) (*domain.Salary, error) {
	return s.repo.GetCurrentSalary(employeeID)
}
func (s *Service) AssignSalary(ctx context.Context, s2 *domain.Salary) error {
	if s2.EmployeeID == 0 {
		return apperrors.BadRequest("employee_id is required")
	}
	s2.EffectiveFrom = time.Now()
	return s.repo.CreateSalary(s2)
}
func (s *Service) ListSalaryComponents(ctx context.Context) ([]domain.SalaryComponent, error) {
	return s.repo.ListSalaryComponents()
}
func (s *Service) CreateSalaryComponent(ctx context.Context, sc *domain.SalaryComponent) error {
	return s.repo.CreateSalaryComponent(sc)
}

// Payroll - with outbox event for finance integration
func (s *Service) RunPayroll(ctx context.Context, employeeID uint, month time.Time, processedBy uint) (*domain.PayrollRun, error) {
	sal, err := s.repo.GetCurrentSalary(employeeID)
	if err != nil {
		return nil, apperrors.BadRequest("no active salary found for employee")
	}

	// Calculate deductions
	var totalDeductions float64
	details, _ := s.repo.GetSalaryDetails(sal.ID)
	for _, d := range details {
		var comp domain.SalaryComponent
		s.db.First(&comp, d.SalaryComponentID)
		if comp.Type == "deduction" {
			totalDeductions += d.Amount
		}
	}

	pr := &domain.PayrollRun{
		EmployeeID:      employeeID,
		Month:            month,
		GrossPay:         sal.BasePay,
		TotalDeductions:  totalDeductions,
		NetPay:           sal.NetSalary,
		ProcessedAt:      time.Now(),
		ProcessedBy:      &processedBy,
	}

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pr).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "PayrollRun", fmt.Sprintf("%d", pr.ID),
			eventbus.EventPayrollProcessed,
			eventbus.PayrollProcessedPayload{
				PayrollID:  pr.ID,
				EmployeeID: pr.EmployeeID,
				Month:      pr.Month.Format("2006-01"),
				GrossPay:   pr.GrossPay,
				NetPay:     pr.NetPay,
			},
		)
	})
	if txErr != nil {
		return nil, txErr
	}
	return pr, nil
}
func (s *Service) ListPayrollRuns(ctx context.Context, employeeID uint) ([]domain.PayrollRun, error) {
	return s.repo.ListPayrollRuns(employeeID)
}

// Leave
func (s *Service) GetLeaveBalances(ctx context.Context, employeeID uint, year int) ([]domain.LeaveBalance, error) {
	return s.repo.ListLeaveBalances(employeeID, uint(year))
}
func (s *Service) RequestLeave(ctx context.Context, req *domain.LeaveRequest) error {
	if req.EmployeeID == 0 || req.LeaveTypeID == 0 {
		return apperrors.BadRequest("employee_id and leave_type_id are required")
	}
	req.CreatedAt = time.Now()

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}
		return s.outbox.WriteEvent(tx, "LeaveRequest", fmt.Sprintf("%d", req.ID),
			eventbus.EventLeaveRequested,
			map[string]interface{}{
				"request_id":   req.ID,
				"employee_id":  req.EmployeeID,
				"leave_type_id": req.LeaveTypeID,
				"start_date":   req.StartDate,
				"end_date":     req.EndDate,
			},
		)
	})
	return txErr
}
func (s *Service) ListLeaveRequests(ctx context.Context, employeeID uint, page, pageSize int) ([]domain.LeaveRequest, int64, error) {
	return s.repo.ListLeaveRequests(employeeID, page, pageSize)
}
func (s *Service) ApproveLeave(ctx context.Context, id, approvedBy uint) error {
	req, err := s.repo.GetLeaveRequest(id)
	if err != nil {
		return apperrors.NotFound("leave request not found")
	}
	now := time.Now()
	req.ApprovedBy = &approvedBy
	req.ApprovedAt = &now

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(req).Error; err != nil {
			return err
		}
		// Update leave balance
		var lb domain.LeaveBalance
		if err := tx.Where("employee_id = ? AND leave_type_id = ? AND year = ?",
			req.EmployeeID, req.LeaveTypeID, now.Year()).First(&lb).Error; err == nil {
			days := req.EndDate.Sub(req.StartDate).Hours() / 24
			lb.UsedQuota += days
			tx.Save(&lb)
		}
		return s.outbox.WriteEvent(tx, "LeaveRequest", fmt.Sprintf("%d", req.ID),
			eventbus.EventLeaveApproved,
			map[string]interface{}{
				"request_id":   req.ID,
				"employee_id":  req.EmployeeID,
				"approved_by":  approvedBy,
			},
		)
	})
	return txErr
}
func (s *Service) RejectLeave(ctx context.Context, id uint) error {
	req, err := s.repo.GetLeaveRequest(id)
	if err != nil {
		return apperrors.NotFound("leave request not found")
	}
	var rejected uint = 3
	req.StatusID = &rejected
	return s.repo.UpdateLeaveRequest(req)
}

// Attendance
func (s *Service) MarkAttendance(ctx context.Context, a *domain.HRAttendance) error {
	if a.EmployeeID == 0 {
		return apperrors.BadRequest("employee_id is required")
	}
	return s.repo.MarkHRAttendance(a)
}
func (s *Service) GetAttendance(ctx context.Context, employeeID uint, from, to time.Time) ([]domain.HRAttendance, error) {
	return s.repo.GetHRAttendance(employeeID, from, to)
}

// Recruitment
func (s *Service) ListJobs(ctx context.Context) ([]domain.RecruitmentJob, error) {
	return s.repo.ListJobs()
}
func (s *Service) GetJob(ctx context.Context, id uint) (*domain.RecruitmentJob, error) {
	return s.repo.GetJob(id)
}
func (s *Service) PostJob(ctx context.Context, j *domain.RecruitmentJob) error {
	j.PostedDate = time.Now()
	return s.repo.CreateJob(j)
}
func (s *Service) ListJobApplications(ctx context.Context, jobID uint) ([]domain.JobApplication, error) {
	return s.repo.ListJobApplications(jobID)
}
func (s *Service) ApplyForJob(ctx context.Context, ja *domain.JobApplication) error {
	if ja.ApplicantName == "" || ja.Email == "" {
		return apperrors.BadRequest("applicant name and email are required")
	}
	ja.AppliedAt = time.Now()
	return s.repo.CreateJobApplication(ja)
}
func (s *Service) UpdateJobApplicationStatus(ctx context.Context, id, statusID uint) error {
	return s.repo.UpdateJobApplicationStatus(id, statusID)
}

// Stats dashboard
func (s *Service) GetHRStats(ctx context.Context) (map[string]interface{}, error) {
	total, _ := s.repo.CountEmployees()
	faculty, _ := s.repo.ListFaculty(0)
	jobs, _ := s.repo.ListJobs()
	return map[string]interface{}{
		"total_employees": total,
		"total_faculty":   len(faculty),
		"open_jobs":       len(jobs),
	}, nil
}
