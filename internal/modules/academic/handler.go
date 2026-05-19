package academicmod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(authMW)

	// Academic Terms
	api.HandleFunc("/academic/terms", h.ListTerms).Methods("GET")
	api.HandleFunc("/academic/terms/current", h.CurrentTerm).Methods("GET")
	api.HandleFunc("/academic/terms", h.CreateTerm).Methods("POST")
	api.HandleFunc("/academic/terms/{id:[0-9]+}", h.GetTerm).Methods("GET")
	api.HandleFunc("/academic/terms/{id:[0-9]+}", h.UpdateTerm).Methods("PUT")
	api.HandleFunc("/academic/terms/{id:[0-9]+}/set-current", h.SetCurrentTerm).Methods("POST")

	// Programs
	api.HandleFunc("/academic/programs", h.ListPrograms).Methods("GET")
	api.HandleFunc("/academic/programs", h.CreateProgram).Methods("POST")
	api.HandleFunc("/academic/programs/{id:[0-9]+}", h.GetProgram).Methods("GET")
	api.HandleFunc("/academic/programs/{id:[0-9]+}", h.UpdateProgram).Methods("PUT")
	api.HandleFunc("/academic/programs/{id:[0-9]+}/curriculum", h.GetCurriculum).Methods("GET")
	api.HandleFunc("/academic/programs/{id:[0-9]+}/subjects", h.AddSubjectToProgram).Methods("POST")
	api.HandleFunc("/academic/programs/{id:[0-9]+}/subjects/{subject_id:[0-9]+}", h.RemoveSubjectFromProgram).Methods("DELETE")

	// Subjects
	api.HandleFunc("/academic/subjects", h.ListSubjects).Methods("GET")
	api.HandleFunc("/academic/subjects", h.CreateSubject).Methods("POST")
	api.HandleFunc("/academic/subjects/{id:[0-9]+}", h.GetSubject).Methods("GET")
	api.HandleFunc("/academic/subjects/{id:[0-9]+}", h.UpdateSubject).Methods("PUT")

	// Batches & Sections
	api.HandleFunc("/academic/batches", h.ListBatches).Methods("GET")
	api.HandleFunc("/academic/batches", h.CreateBatch).Methods("POST")
	api.HandleFunc("/academic/batches/{id:[0-9]+}/sections", h.ListSections).Methods("GET")
	api.HandleFunc("/academic/batches/{id:[0-9]+}/sections", h.CreateSection).Methods("POST")

	// Course Offerings
	api.HandleFunc("/academic/offerings", h.ListOfferings).Methods("GET")
	api.HandleFunc("/academic/offerings", h.CreateOffering).Methods("POST")
	api.HandleFunc("/academic/offerings/{id:[0-9]+}", h.GetOffering).Methods("GET")
	api.HandleFunc("/academic/offerings/{id:[0-9]+}", h.UpdateOffering).Methods("PUT")
	api.HandleFunc("/academic/offerings/{id:[0-9]+}/timetable", h.GetOfferingTimetable).Methods("GET")
	api.HandleFunc("/academic/offerings/{id:[0-9]+}/timetable", h.CreateTimetableEntry).Methods("POST")

	// Registrations
	api.HandleFunc("/academic/term-registrations", h.RegisterForTerm).Methods("POST")
	api.HandleFunc("/academic/term-registrations", h.ListTermRegistrations).Methods("GET")
	api.HandleFunc("/academic/course-registrations", h.RegisterForCourse).Methods("POST")
	api.HandleFunc("/academic/course-registrations/{id:[0-9]+}/drop", h.DropCourse).Methods("POST")
	api.HandleFunc("/academic/students/{student_id:[0-9]+}/courses", h.StudentCourses).Methods("GET")
	api.HandleFunc("/academic/students/{student_id:[0-9]+}/timetable", h.StudentTimetable).Methods("GET")
	api.HandleFunc("/academic/students/me/timetable", h.MyTimetable).Methods("GET")
	api.HandleFunc("/academic/students/me/courses", h.MyCourses).Methods("GET")
	api.HandleFunc("/academic/students/me/attendance", h.MyAttendance).Methods("GET")

	// Attendance
	api.HandleFunc("/academic/attendance", h.MarkAttendance).Methods("POST")
	api.HandleFunc("/academic/students/{student_id:[0-9]+}/attendance", h.GetStudentAttendance).Methods("GET")
	api.HandleFunc("/academic/students/{student_id:[0-9]+}/attendance/summary", h.GetAttendanceSummary).Methods("GET")

	// Calendar
	api.HandleFunc("/academic/calendar", h.ListCalendar).Methods("GET")
	api.HandleFunc("/academic/calendar", h.CreateCalendarEvent).Methods("POST")
}

// Terms
func (h *Handler) ListTerms(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseUint(r.URL.Query().Get("campus_id"), 10, 64)
	data, err := h.service.ListTerms(r.Context(), uint(cid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CurrentTerm(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetCurrentTerm(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetTerm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetTerm(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateTerm(w http.ResponseWriter, r *http.Request) {
	var t domain.AcademicTerm
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateTerm(r.Context(), &t); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, t)
}
func (h *Handler) UpdateTerm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var t domain.AcademicTerm
	json.NewDecoder(r.Body).Decode(&t)
	if err := h.service.UpdateTerm(r.Context(), uint(id), &t); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, t)
}
func (h *Handler) SetCurrentTerm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.SetCurrentTerm(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "current term updated"})
}

// Programs
func (h *Handler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	did, _ := strconv.ParseUint(r.URL.Query().Get("department_id"), 10, 64)
	data, err := h.service.ListPrograms(r.Context(), uint(did))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetProgram(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	var p domain.Program
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateProgram(r.Context(), &p); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, p)
}
func (h *Handler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var p domain.Program
	json.NewDecoder(r.Body).Decode(&p)
	if err := h.service.UpdateProgram(r.Context(), uint(id), &p); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, p)
}
func (h *Handler) GetCurriculum(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetProgramWithCurriculum(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AddSubjectToProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var ps domain.ProgramSubject
	json.NewDecoder(r.Body).Decode(&ps)
	ps.ProgramID = uint(id)
	if err := h.service.AddSubjectToProgram(r.Context(), &ps); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ps)
}
func (h *Handler) RemoveSubjectFromProgram(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	sid, _ := strconv.ParseUint(mux.Vars(r)["subject_id"], 10, 64)
	if err := h.service.RemoveSubjectFromProgram(r.Context(), uint(id), uint(sid)); err != nil {
		response.Error(w, err)
		return
	}
	response.NoContent(w)
}

// Subjects
func (h *Handler) ListSubjects(w http.ResponseWriter, r *http.Request) {
	did, _ := strconv.ParseUint(r.URL.Query().Get("department_id"), 10, 64)
	data, err := h.service.ListSubjects(r.Context(), uint(did))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetSubject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetSubject(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	var s domain.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateSubject(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}
func (h *Handler) UpdateSubject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.Subject
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.UpdateSubject(r.Context(), uint(id), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, s)
}

// Batches & Sections
func (h *Handler) ListBatches(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.ParseUint(r.URL.Query().Get("program_id"), 10, 64)
	data, err := h.service.ListBatches(r.Context(), uint(pid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateBatch(w http.ResponseWriter, r *http.Request) {
	var b domain.Batch
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.CreateBatch(r.Context(), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, b)
}
func (h *Handler) ListSections(w http.ResponseWriter, r *http.Request) {
	bid, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.ListSections(r.Context(), uint(bid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateSection(w http.ResponseWriter, r *http.Request) {
	bid, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.Section
	json.NewDecoder(r.Body).Decode(&s)
	s.BatchID = uint(bid)
	if err := h.service.CreateSection(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}

// Offerings
func (h *Handler) ListOfferings(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	pid, _ := strconv.ParseUint(r.URL.Query().Get("program_id"), 10, 64)
	data, err := h.service.ListOfferings(r.Context(), uint(tid), uint(pid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetOffering(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetOffering(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateOffering(w http.ResponseWriter, r *http.Request) {
	var o domain.CourseOffering
	json.NewDecoder(r.Body).Decode(&o)
	if err := h.service.CreateOffering(r.Context(), &o); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, o)
}
func (h *Handler) UpdateOffering(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var o domain.CourseOffering
	json.NewDecoder(r.Body).Decode(&o)
	if err := h.service.UpdateOffering(r.Context(), uint(id), &o); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, o)
}
func (h *Handler) GetOfferingTimetable(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetOfferingTimetable(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateTimetableEntry(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var t domain.Timetable
	json.NewDecoder(r.Body).Decode(&t)
	t.OfferingID = uint(id)
	if err := h.service.CreateTimetableEntry(r.Context(), &t); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, t)
}

// Registrations
func (h *Handler) RegisterForTerm(w http.ResponseWriter, r *http.Request) {
	var tr domain.TermRegistration
	json.NewDecoder(r.Body).Decode(&tr)
	if err := h.service.RegisterStudentForTerm(r.Context(), &tr); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, tr)
}
func (h *Handler) ListTermRegistrations(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.ListTermRegistrations(r.Context(), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) RegisterForCourse(w http.ResponseWriter, r *http.Request) {
	var cr domain.CourseRegistration
	json.NewDecoder(r.Body).Decode(&cr)
	if err := h.service.RegisterStudentForCourse(r.Context(), &cr); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, cr)
}
func (h *Handler) DropCourse(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.DropCourse(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "course dropped"})
}
func (h *Handler) StudentCourses(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.ListStudentCourses(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) StudentTimetable(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetStudentTimetable(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyTimetable(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetStudentTimetable(r.Context(), userID, uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyCourses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.ListStudentCourses(r.Context(), userID, uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyAttendance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetAttendanceSummary(r.Context(), userID, uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Attendance
func (h *Handler) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	var a domain.StudentAttendance
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.service.MarkAttendance(r.Context(), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, a)
}
func (h *Handler) GetStudentAttendance(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	oid, _ := strconv.ParseUint(r.URL.Query().Get("offering_id"), 10, 64)
	data, err := h.service.GetStudentAttendance(r.Context(), uint(sid), uint(oid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetAttendanceSummary(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetAttendanceSummary(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Calendar
func (h *Handler) ListCalendar(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseUint(r.URL.Query().Get("campus_id"), 10, 64)
	data, err := h.service.ListCalendar(r.Context(), uint(cid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateCalendarEvent(w http.ResponseWriter, r *http.Request) {
	var e domain.AcademicCalendar
	json.NewDecoder(r.Body).Decode(&e)
	if err := h.service.CreateCalendarEvent(r.Context(), &e); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, e)
}
