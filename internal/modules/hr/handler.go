package hrmod

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	api := r.PathPrefix("/api/v1/hr").Subrouter()
	api.Use(authMW)

	// Lookups
	api.HandleFunc("/designations", h.ListDesignations).Methods("GET")
	api.HandleFunc("/designations", h.CreateDesignation).Methods("POST")
	api.HandleFunc("/employment-types", h.ListEmploymentTypes).Methods("GET")
	api.HandleFunc("/leave-types", h.ListLeaveTypes).Methods("GET")
	api.HandleFunc("/leave-types", h.CreateLeaveType).Methods("POST")
	api.HandleFunc("/salary-components", h.ListSalaryComponents).Methods("GET")
	api.HandleFunc("/salary-components", h.CreateSalaryComponent).Methods("POST")

	// Stats
	api.HandleFunc("/stats", h.Stats).Methods("GET")

	// Employees
	api.HandleFunc("/employees", h.ListEmployees).Methods("GET")
	api.HandleFunc("/employees", h.CreateEmployee).Methods("POST")
	api.HandleFunc("/employees/me", h.MyProfile).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}", h.GetEmployee).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}", h.UpdateEmployee).Methods("PUT")
	api.HandleFunc("/employees/{id:[0-9]+}/deactivate", h.DeactivateEmployee).Methods("POST")
	api.HandleFunc("/employees/{id:[0-9]+}/transfer", h.TransferDepartment).Methods("POST")
	api.HandleFunc("/employees/{id:[0-9]+}/dept-history", h.DeptHistory).Methods("GET")

	// Faculty
	api.HandleFunc("/faculty", h.ListFaculty).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}/faculty-profile", h.GetFacultyProfile).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}/faculty-profile", h.UpsertFacultyProfile).Methods("PUT")

	// Staff
	api.HandleFunc("/employees/{id:[0-9]+}/staff-profile", h.GetStaffProfile).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}/staff-profile", h.UpsertStaffProfile).Methods("PUT")

	// Salary
	api.HandleFunc("/employees/{id:[0-9]+}/salary", h.GetCurrentSalary).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}/salary", h.AssignSalary).Methods("POST")

	// Payroll
	api.HandleFunc("/payroll/run", h.RunPayroll).Methods("POST")
	api.HandleFunc("/payroll", h.ListPayrollRuns).Methods("GET")
	api.HandleFunc("/employees/{id:[0-9]+}/payslips", h.EmployeePayslips).Methods("GET")

	// Leave
	api.HandleFunc("/employees/{id:[0-9]+}/leave-balances", h.LeaveBalances).Methods("GET")
	api.HandleFunc("/leave-requests", h.RequestLeave).Methods("POST")
	api.HandleFunc("/leave-requests", h.ListLeaveRequests).Methods("GET")
	api.HandleFunc("/leave-requests/{id:[0-9]+}/approve", h.ApproveLeave).Methods("POST")
	api.HandleFunc("/leave-requests/{id:[0-9]+}/reject", h.RejectLeave).Methods("POST")

	// Attendance
	api.HandleFunc("/attendance", h.MarkAttendance).Methods("POST")
	api.HandleFunc("/employees/{id:[0-9]+}/attendance", h.GetAttendance).Methods("GET")

	// Recruitment
	api.HandleFunc("/jobs", h.ListJobs).Methods("GET")
	api.HandleFunc("/jobs", h.PostJob).Methods("POST")
	api.HandleFunc("/jobs/{id:[0-9]+}", h.GetJob).Methods("GET")
	api.HandleFunc("/jobs/{id:[0-9]+}/applications", h.ListJobApplications).Methods("GET")
	api.HandleFunc("/jobs/{id:[0-9]+}/applications", h.ApplyForJob).Methods("POST")
	api.HandleFunc("/job-applications/{id:[0-9]+}/status", h.UpdateJobApplicationStatus).Methods("PUT")
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetHRStats(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListDesignations(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListDesignations(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateDesignation(w http.ResponseWriter, r *http.Request) {
	var d domain.Designation
	json.NewDecoder(r.Body).Decode(&d)
	if err := h.service.CreateDesignation(r.Context(), &d); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, d)
}
func (h *Handler) ListEmploymentTypes(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListEmploymentTypes(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListLeaveTypes(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListLeaveTypes(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateLeaveType(w http.ResponseWriter, r *http.Request) {
	var lt domain.LeaveType
	json.NewDecoder(r.Body).Decode(&lt)
	if err := h.service.CreateLeaveType(r.Context(), &lt); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, lt)
}
func (h *Handler) ListSalaryComponents(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListSalaryComponents(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateSalaryComponent(w http.ResponseWriter, r *http.Request) {
	var sc domain.SalaryComponent
	json.NewDecoder(r.Body).Decode(&sc)
	if err := h.service.CreateSalaryComponent(r.Context(), &sc); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, sc)
}
func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	did, _ := strconv.ParseUint(r.URL.Query().Get("department_id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	data, total, err := h.service.ListEmployees(r.Context(), uint(did), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	response.List(w, data, total, page, pageSize)
}
func (h *Handler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetEmployee(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	data, err := h.service.GetMyProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var e domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateEmployee(r.Context(), &e); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, e)
}
func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var e domain.Employee
	json.NewDecoder(r.Body).Decode(&e)
	if err := h.service.UpdateEmployee(r.Context(), uint(id), &e); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, e)
}
func (h *Handler) DeactivateEmployee(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.DeactivateEmployee(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "employee deactivated"})
}
func (h *Handler) TransferDepartment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		DepartmentID uint `json:"department_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.service.TransferDepartment(r.Context(), uint(id), req.DepartmentID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "department transfer complete"})
}
func (h *Handler) DeptHistory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetDeptHistory(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListFaculty(w http.ResponseWriter, r *http.Request) {
	did, _ := strconv.ParseUint(r.URL.Query().Get("department_id"), 10, 64)
	data, err := h.service.ListFaculty(r.Context(), uint(did))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetFacultyProfile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetFacultyProfile(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) UpsertFacultyProfile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var f domain.Faculty
	json.NewDecoder(r.Body).Decode(&f)
	f.EmployeeID = uint(id)
	if err := h.service.UpsertFacultyProfile(r.Context(), &f); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, f)
}
func (h *Handler) GetStaffProfile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetStaffProfile(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) UpsertStaffProfile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.Staff
	json.NewDecoder(r.Body).Decode(&s)
	s.EmployeeID = uint(id)
	if err := h.service.UpsertStaffProfile(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, s)
}
func (h *Handler) GetCurrentSalary(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetCurrentSalary(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AssignSalary(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.Salary
	json.NewDecoder(r.Body).Decode(&s)
	s.EmployeeID = uint(id)
	if err := h.service.AssignSalary(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}
func (h *Handler) RunPayroll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID uint   `json:"employee_id"`
		Month      string `json:"month"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	month, _ := time.Parse("2006-01", req.Month)
	userID := middleware.GetUserID(r.Context())
	data, err := h.service.RunPayroll(r.Context(), req.EmployeeID, month, userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) ListPayrollRuns(w http.ResponseWriter, r *http.Request) {
	eid, _ := strconv.ParseUint(r.URL.Query().Get("employee_id"), 10, 64)
	data, err := h.service.ListPayrollRuns(r.Context(), uint(eid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) EmployeePayslips(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.ListPayrollRuns(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) LeaveBalances(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if year == 0 {
		year = time.Now().Year()
	}
	data, err := h.service.GetLeaveBalances(r.Context(), uint(id), year)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) RequestLeave(w http.ResponseWriter, r *http.Request) {
	var req domain.LeaveRequest
	json.NewDecoder(r.Body).Decode(&req)
	if req.EmployeeID == 0 {
		req.EmployeeID = middleware.GetUserID(r.Context())
	}
	if err := h.service.RequestLeave(r.Context(), &req); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, req)
}
func (h *Handler) ListLeaveRequests(w http.ResponseWriter, r *http.Request) {
	eid, _ := strconv.ParseUint(r.URL.Query().Get("employee_id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	data, total, err := h.service.ListLeaveRequests(r.Context(), uint(eid), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	response.List(w, data, total, page, pageSize)
}
func (h *Handler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.ApproveLeave(r.Context(), uint(id), userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "leave approved"})
}
func (h *Handler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.RejectLeave(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "leave rejected"})
}
func (h *Handler) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	var a domain.HRAttendance
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.service.MarkAttendance(r.Context(), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, a)
}
func (h *Handler) GetAttendance(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	from, _ := time.Parse("2006-01-02", r.URL.Query().Get("from"))
	to, _ := time.Parse("2006-01-02", r.URL.Query().Get("to"))
	if from.IsZero() {
		from = time.Now().AddDate(0, -1, 0)
	}
	if to.IsZero() {
		to = time.Now()
	}
	data, err := h.service.GetAttendance(r.Context(), uint(id), from, to)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListJobs(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetJob(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {
	var j domain.RecruitmentJob
	json.NewDecoder(r.Body).Decode(&j)
	if err := h.service.PostJob(r.Context(), &j); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, j)
}
func (h *Handler) ListJobApplications(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.ListJobApplications(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ApplyForJob(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var ja domain.JobApplication
	json.NewDecoder(r.Body).Decode(&ja)
	ja.JobID = uint(id)
	if err := h.service.ApplyForJob(r.Context(), &ja); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ja)
}
func (h *Handler) UpdateJobApplicationStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		StatusID uint `json:"status_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.service.UpdateJobApplicationStatus(r.Context(), uint(id), req.StatusID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}
