package studentmod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

// Handler holds HTTP handlers for the student module.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes mounts student routes on the given router.
func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	sub := r.PathPrefix("/api/v1/students").Subrouter()
	sub.Use(authMW)

	sub.HandleFunc("", h.List).Methods("GET")
	sub.HandleFunc("/enroll", h.Enroll).Methods("POST")
	sub.HandleFunc("/me", h.MyProfile).Methods("GET")
	sub.HandleFunc("/me/dashboard", h.MyDashboard).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}", h.GetByID).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}/guardians", h.GetGuardians).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}/guardians", h.AddGuardian).Methods("POST")
	sub.HandleFunc("/{id:[0-9]+}/grievances", h.GetGrievances).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}/grievances", h.FileGrievance).Methods("POST")
}

func (h *Handler) Enroll(w http.ResponseWriter, r *http.Request) {
	var req EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}
	student, err := h.service.EnrollStudent(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, student)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	students, total, err := h.service.List(r.Context(), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	response.List(w, students, total, page, pageSize)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	student, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, student)
}

func (h *Handler) MyProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	student, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, student)
}

func (h *Handler) MyDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	student, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	dashboard, err := h.service.GetDashboard(r.Context(), student.ID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, dashboard)
}

func (h *Handler) GetGuardians(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	guardians, err := h.service.GetGuardians(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, guardians)
}

func (h *Handler) AddGuardian(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var g domain.Guardian
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		response.Error(w, err)
		return
	}
	g.StudentID = uint(id)
	if err := h.service.AddGuardian(r.Context(), &g); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, g)
}

func (h *Handler) GetGrievances(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	grievances, err := h.service.GetGrievances(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, grievances)
}

func (h *Handler) FileGrievance(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var g domain.Grievance
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		response.Error(w, err)
		return
	}
	g.StudentID = uint(id)
	if err := h.service.FileGrievance(r.Context(), &g); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, g)
}
