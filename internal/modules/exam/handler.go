package exammod

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
	api := r.PathPrefix("/api/v1/exam").Subrouter()
	api.Use(authMW)

	// Components
	api.HandleFunc("/components", h.ListComponents).Methods("GET")
	api.HandleFunc("/components", h.CreateComponent).Methods("POST")
	api.HandleFunc("/components/{id:[0-9]+}", h.UpdateComponent).Methods("PUT")

	// Schedules
	api.HandleFunc("/schedules", h.ListSchedules).Methods("GET")
	api.HandleFunc("/schedules", h.CreateSchedule).Methods("POST")
	api.HandleFunc("/schedules/{id:[0-9]+}", h.GetSchedule).Methods("GET")
	api.HandleFunc("/schedules/{id:[0-9]+}", h.UpdateSchedule).Methods("PUT")

	// Results
	api.HandleFunc("/results", h.EnterResult).Methods("POST")
	api.HandleFunc("/results/bulk", h.BulkEnterResults).Methods("POST")
	api.HandleFunc("/results/{id:[0-9]+}", h.UpdateResult).Methods("PUT")
	api.HandleFunc("/results/{id:[0-9]+}/component-marks", h.GetComponentMarks).Methods("GET")
	api.HandleFunc("/results/{id:[0-9]+}/component-marks", h.EnterComponentMarks).Methods("POST")
	api.HandleFunc("/results/publish", h.PublishResults).Methods("POST")

	// Student Results
	api.HandleFunc("/students/{student_id:[0-9]+}/results", h.StudentResults).Methods("GET")
	api.HandleFunc("/students/{student_id:[0-9]+}/transcript", h.Transcript).Methods("GET")
	api.HandleFunc("/students/{student_id:[0-9]+}/sgpa", h.SGPA).Methods("GET")
	api.HandleFunc("/students/{student_id:[0-9]+}/cgpa", h.CGPA).Methods("GET")
	api.HandleFunc("/students/me/results", h.MyResults).Methods("GET")
	api.HandleFunc("/students/me/transcript", h.MyTranscript).Methods("GET")
	api.HandleFunc("/students/me/cgpa", h.MyCGPA).Methods("GET")

	// Revaluation
	api.HandleFunc("/revaluations", h.RequestRevaluation).Methods("POST")
	api.HandleFunc("/revaluations", h.ListRevaluations).Methods("GET")
	api.HandleFunc("/revaluations/{id:[0-9]+}/process", h.ProcessRevaluation).Methods("POST")

	// Supplementary
	api.HandleFunc("/supplementary", h.ListSupplementary).Methods("GET")
	api.HandleFunc("/supplementary", h.CreateSupplementary).Methods("POST")
}

func (h *Handler) ListComponents(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("subject_id"), 10, 64)
	data, err := h.service.ListComponents(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateComponent(w http.ResponseWriter, r *http.Request) {
	var c domain.ExamComponent
	json.NewDecoder(r.Body).Decode(&c)
	if err := h.service.CreateComponent(r.Context(), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, c)
}
func (h *Handler) UpdateComponent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var c domain.ExamComponent
	json.NewDecoder(r.Body).Decode(&c)
	if err := h.service.UpdateComponent(r.Context(), uint(id), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, c)
}
func (h *Handler) ListSchedules(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	sid, _ := strconv.ParseUint(r.URL.Query().Get("subject_id"), 10, 64)
	data, err := h.service.ListSchedules(r.Context(), uint(tid), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetSchedule(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var s domain.ExamSchedule
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.CreateSchedule(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}
func (h *Handler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var s domain.ExamSchedule
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.UpdateSchedule(r.Context(), uint(id), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, s)
}
func (h *Handler) EnterResult(w http.ResponseWriter, r *http.Request) {
	var res domain.Result
	json.NewDecoder(r.Body).Decode(&res)
	if err := h.service.EnterResult(r.Context(), &res); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, res)
}
func (h *Handler) BulkEnterResults(w http.ResponseWriter, r *http.Request) {
	var results []domain.Result
	json.NewDecoder(r.Body).Decode(&results)
	if err := h.service.BulkEnterResults(r.Context(), results); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]int{"count": len(results)})
}
func (h *Handler) UpdateResult(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var res domain.Result
	json.NewDecoder(r.Body).Decode(&res)
	if err := h.service.UpdateResult(r.Context(), uint(id), &res); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, res)
}
func (h *Handler) PublishResults(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TermID uint `json:"term_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.PublishResults(r.Context(), req.TermID, userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "results published"})
}
func (h *Handler) GetComponentMarks(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetComponentMarks(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) EnterComponentMarks(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var cm domain.ComponentMarks
	json.NewDecoder(r.Body).Decode(&cm)
	cm.ResultID = uint(id)
	if err := h.service.EnterComponentMarks(r.Context(), &cm); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, cm)
}
func (h *Handler) StudentResults(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetStudentResults(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) Transcript(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	data, err := h.service.GetTranscript(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) SGPA(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetSGPA(r.Context(), uint(sid), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CGPA(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(mux.Vars(r)["student_id"], 10, 64)
	data, err := h.service.GetCGPA(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyResults(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.GetStudentResults(r.Context(), userID, uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyTranscript(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	data, err := h.service.GetTranscript(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) MyCGPA(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	data, err := h.service.GetCGPA(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) RequestRevaluation(w http.ResponseWriter, r *http.Request) {
	var req domain.RevaluationRequest
	json.NewDecoder(r.Body).Decode(&req)
	if req.StudentID == 0 {
		req.StudentID = middleware.GetUserID(r.Context())
	}
	if err := h.service.RequestRevaluation(r.Context(), &req); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, req)
}
func (h *Handler) ListRevaluations(w http.ResponseWriter, r *http.Request) {
	sid, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListRevaluations(r.Context(), uint(sid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ProcessRevaluation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		ReviewedMarks float64 `json:"reviewed_marks"`
		Grade         string  `json:"grade"`
		Remarks       string  `json:"remarks"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.ProcessRevaluation(r.Context(), uint(id), req.ReviewedMarks, req.Grade, req.Remarks, userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "revaluation processed"})
}
func (h *Handler) ListSupplementary(w http.ResponseWriter, r *http.Request) {
	tid, _ := strconv.ParseUint(r.URL.Query().Get("term_id"), 10, 64)
	data, err := h.service.ListSupplementary(r.Context(), uint(tid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateSupplementary(w http.ResponseWriter, r *http.Request) {
	var s domain.SupplementaryExam
	json.NewDecoder(r.Body).Decode(&s)
	if err := h.service.CreateSupplementary(r.Context(), &s); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, s)
}
