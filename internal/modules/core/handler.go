package coremod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	// Public lookups
	pub := r.PathPrefix("/api/v1").Subrouter()
	pub.HandleFunc("/lookups/genders", h.Genders).Methods("GET")
	pub.HandleFunc("/lookups/categories", h.Categories).Methods("GET")
	pub.HandleFunc("/lookups/blood-groups", h.BloodGroups).Methods("GET")
	pub.HandleFunc("/lookups/status-codes", h.StatusCodes).Methods("GET")

	// Protected core routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(authMW)

	// Universities
	api.HandleFunc("/universities", h.ListUniversities).Methods("GET")
	api.HandleFunc("/universities", h.CreateUniversity).Methods("POST")
	api.HandleFunc("/universities/{id:[0-9]+}", h.GetUniversity).Methods("GET")
	api.HandleFunc("/universities/{id:[0-9]+}", h.UpdateUniversity).Methods("PUT")

	// Campuses
	api.HandleFunc("/campuses", h.ListCampuses).Methods("GET")
	api.HandleFunc("/campuses", h.CreateCampus).Methods("POST")
	api.HandleFunc("/campuses/{id:[0-9]+}", h.GetCampus).Methods("GET")
	api.HandleFunc("/campuses/{id:[0-9]+}", h.UpdateCampus).Methods("PUT")

	// Departments
	api.HandleFunc("/departments", h.ListDepartments).Methods("GET")
	api.HandleFunc("/departments", h.CreateDepartment).Methods("POST")
	api.HandleFunc("/departments/{id:[0-9]+}", h.GetDepartment).Methods("GET")
	api.HandleFunc("/departments/{id:[0-9]+}", h.UpdateDepartment).Methods("PUT")

	// Rooms
	api.HandleFunc("/rooms", h.ListRooms).Methods("GET")
	api.HandleFunc("/rooms", h.CreateRoom).Methods("POST")
}

func (h *Handler) Genders(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListGenders(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListCategories(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) BloodGroups(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListBloodGroups(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) StatusCodes(w http.ResponseWriter, r *http.Request) {
	module := r.URL.Query().Get("module")
	data, err := h.service.ListStatusCodes(r.Context(), module)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

func (h *Handler) ListUniversities(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.ListUniversities(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetUniversity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetUniversity(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateUniversity(w http.ResponseWriter, r *http.Request) {
	var u domain.University
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateUniversity(r.Context(), &u); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, u)
}
func (h *Handler) UpdateUniversity(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var u domain.University
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.UpdateUniversity(r.Context(), uint(id), &u); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, u)
}

func (h *Handler) ListCampuses(w http.ResponseWriter, r *http.Request) {
	uid, _ := strconv.ParseUint(r.URL.Query().Get("university_id"), 10, 64)
	data, err := h.service.ListCampuses(r.Context(), uint(uid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetCampus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetCampus(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateCampus(w http.ResponseWriter, r *http.Request) {
	var c domain.Campus
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateCampus(r.Context(), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, c)
}
func (h *Handler) UpdateCampus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var c domain.Campus
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.UpdateCampus(r.Context(), uint(id), &c); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, c)
}

func (h *Handler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseUint(r.URL.Query().Get("campus_id"), 10, 64)
	data, err := h.service.ListDepartments(r.Context(), uint(cid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetDepartment(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var d domain.Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateDepartment(r.Context(), &d); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, d)
}
func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var d domain.Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.UpdateDepartment(r.Context(), uint(id), &d); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, d)
}

func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseUint(r.URL.Query().Get("campus_id"), 10, 64)
	data, err := h.service.ListRooms(r.Context(), uint(cid))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var rm domain.Room
	if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
		response.Error(w, err)
		return
	}
	if err := h.service.CreateRoom(r.Context(), &rm); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, rm)
}
