package admissionsmod

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
	// Public - open cycles visible without auth
	r.HandleFunc("/api/v1/admissions/cycles/open", h.OpenCycles).Methods("GET")

	api := r.PathPrefix("/api/v1/admissions").Subrouter()
	api.Use(authMW)

	// Cycles
	api.HandleFunc("/cycles", h.ListCycles).Methods("GET")
	api.HandleFunc("/cycles", h.CreateCycle).Methods("POST")
	api.HandleFunc("/cycles/{id:[0-9]+}", h.GetCycle).Methods("GET")
	api.HandleFunc("/cycles/{id:[0-9]+}", h.UpdateCycle).Methods("PUT")
	api.HandleFunc("/cycles/{id:[0-9]+}/close", h.CloseCycle).Methods("POST")
	api.HandleFunc("/cycles/{id:[0-9]+}/stats", h.CycleStats).Methods("GET")
	api.HandleFunc("/cycles/{id:[0-9]+}/seat-allocations", h.ListSeatAllocations).Methods("GET")
	api.HandleFunc("/cycles/{id:[0-9]+}/waitlist", h.GetWaitlist).Methods("GET")

	// Applicants
	api.HandleFunc("/applicants", h.ListApplicants).Methods("GET")
	api.HandleFunc("/applicants", h.Submit).Methods("POST")
	api.HandleFunc("/applicants/{id:[0-9]+}", h.GetApplicant).Methods("GET")
	api.HandleFunc("/applicants/{id:[0-9]+}", h.UpdateApplicant).Methods("PUT")
	api.HandleFunc("/applicants/{id:[0-9]+}/status", h.UpdateStatus).Methods("PUT")
	api.HandleFunc("/applicants/{id:[0-9]+}/status-history", h.StatusHistory).Methods("GET")

	// Documents
	api.HandleFunc("/applicants/{id:[0-9]+}/documents", h.GetDocuments).Methods("GET")
	api.HandleFunc("/applicants/{id:[0-9]+}/documents", h.UploadDocument).Methods("POST")
	api.HandleFunc("/documents/{id:[0-9]+}/verify", h.VerifyDocument).Methods("POST")

	// Seat Allocation & Waitlist
	api.HandleFunc("/seat-allocations", h.AllocateSeat).Methods("POST")
	api.HandleFunc("/applicants/{id:[0-9]+}/seat-allocation", h.GetSeatAllocation).Methods("GET")
	api.HandleFunc("/waitlist", h.AddToWaitlist).Methods("POST")

	// Conversion
	api.HandleFunc("/applicants/{id:[0-9]+}/convert-to-student", h.ConvertToStudent).Methods("POST")
}

func (h *Handler) OpenCycles(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetOpenCycles(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListCycles(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListCycles(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetCycle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetCycle(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateCycle(w http.ResponseWriter, r *http.Request) {
	var c domain.AdmissionCycle
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateCycle(r.Context(), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, c)
}
func (h *Handler) UpdateCycle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var c domain.AdmissionCycle
	json.NewDecoder(r.Body).Decode(&c)
	if err := h.service.UpdateCycle(r.Context(), uint(id), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, c)
}
func (h *Handler) CloseCycle(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.CloseCycle(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "cycle closed"})
}
func (h *Handler) CycleStats(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetCycleStats(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListApplicants(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseUint(r.URL.Query().Get("cycle_id"), 10, 64)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	data, total, err := h.service.ListApplicants(r.Context(), uint(cid), page, pageSize)
	if err != nil {
		response.Error(w, err)
		return
	}
	if page < 1 { page = 1 }
	if pageSize < 1 { pageSize = 20 }
	response.List(w, data, total, page, pageSize)
}
func (h *Handler) GetApplicant(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetApplicant(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	var a domain.Applicant
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.Submit(r.Context(), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, a)
}
func (h *Handler) UpdateApplicant(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var a domain.Applicant
	json.NewDecoder(r.Body).Decode(&a)
	if err := h.service.UpdateApplicant(r.Context(), uint(id), &a); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, a)
}
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		StatusID uint `json:"status_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	userID := middleware.GetUserID(r.Context())
	if err := h.service.UpdateStatus(r.Context(), uint(id), req.StatusID, userID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}
func (h *Handler) StatusHistory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetStatusHistory(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetDocuments(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var d domain.Document
	json.NewDecoder(r.Body).Decode(&d)
	d.ApplicantID = uint(id)
	if err := h.service.UploadDocument(r.Context(), &d); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, d)
}
func (h *Handler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	verifiedBy := middleware.GetUserID(r.Context())
	if err := h.service.VerifyDocument(r.Context(), uint(id), verifiedBy); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "document verified"})
}
func (h *Handler) AllocateSeat(w http.ResponseWriter, r *http.Request) {
	var sa domain.SeatAllocation
	json.NewDecoder(r.Body).Decode(&sa)
	if err := h.service.AllocateSeat(r.Context(), &sa); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, sa)
}
func (h *Handler) GetSeatAllocation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetSeatAllocation(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ListSeatAllocations(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.ListSeatAllocations(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AddToWaitlist(w http.ResponseWriter, r *http.Request) {
	var wl domain.Waitlist
	json.NewDecoder(r.Body).Decode(&wl)
	if err := h.service.AddToWaitlist(r.Context(), &wl); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, wl)
}
func (h *Handler) GetWaitlist(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetWaitlist(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) ConvertToStudent(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		StudentID uint `json:"student_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.service.ConvertToStudent(r.Context(), uint(id), req.StudentID); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "applicant converted to student"})
}
