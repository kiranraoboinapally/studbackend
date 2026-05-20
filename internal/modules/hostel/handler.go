package hostelmod

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
	api := r.PathPrefix("/api/v1/hostel").Subrouter()
	api.Use(authMW)

	// Stats
	api.HandleFunc("/stats", h.Stats).Methods("GET")

	// Hostels
	api.HandleFunc("/hostels", h.ListHostels).Methods("GET")
	api.HandleFunc("/hostels", h.CreateHostel).Methods("POST")
	api.HandleFunc("/hostels/{id:[0-9]+}", h.GetHostel).Methods("GET")
	api.HandleFunc("/hostels/{id:[0-9]+}", h.UpdateHostel).Methods("PUT")

	// Rooms
	api.HandleFunc("/rooms", h.ListRooms).Methods("GET")
	api.HandleFunc("/rooms", h.CreateRoom).Methods("POST")
	api.HandleFunc("/rooms/{id:[0-9]+}", h.GetRoom).Methods("GET")
	api.HandleFunc("/rooms/{id:[0-9]+}", h.UpdateRoom).Methods("PUT")

	// Beds
	api.HandleFunc("/beds", h.ListBeds).Methods("GET")
	api.HandleFunc("/beds", h.CreateBed).Methods("POST")

	// Allocations
	api.HandleFunc("/allocations", h.ListAllocations).Methods("GET")
	api.HandleFunc("/allocations/allocate", h.AllocateRoom).Methods("POST")
	api.HandleFunc("/allocations/{id:[0-9]+}/deallocate", h.DeallocateRoom).Methods("POST")

	// Mess Bills
	api.HandleFunc("/mess-bills", h.ListMessBills).Methods("GET")
	api.HandleFunc("/mess-bills", h.CreateMessBill).Methods("POST")
	api.HandleFunc("/mess-bills/{id:[0-9]+}/pay", h.PayMessBill).Methods("POST")

	// Maintenance Requests
	api.HandleFunc("/maintenance-requests", h.ListMaintenanceRequests).Methods("GET")
	api.HandleFunc("/maintenance-requests", h.CreateMaintenanceRequest).Methods("POST")
	api.HandleFunc("/maintenance-requests/{id:[0-9]+}/resolve", h.ResolveMaintenanceRequest).Methods("POST")

	// Visitor Logs
	api.HandleFunc("/visitor-logs", h.ListVisitorLogs).Methods("GET")
	api.HandleFunc("/visitor-logs", h.CreateVisitorLog).Methods("POST")
	api.HandleFunc("/visitor-logs/{id:[0-9]+}/exit", h.ExitVisitor).Methods("POST")
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetStats(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Hostels
func (h *Handler) ListHostels(w http.ResponseWriter, r *http.Request) {
	campusID, _ := strconv.ParseUint(r.URL.Query().Get("campus_id"), 10, 64)
	data, err := h.service.ListHostels(r.Context(), uint(campusID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetHostel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetHostel(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateHostel(w http.ResponseWriter, r *http.Request) {
	var hst domain.Hostel
	json.NewDecoder(r.Body).Decode(&hst)
	if err := h.service.CreateHostel(r.Context(), &hst); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, hst)
}
func (h *Handler) UpdateHostel(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var hst domain.Hostel
	json.NewDecoder(r.Body).Decode(&hst)
	if err := h.service.UpdateHostel(r.Context(), uint(id), &hst); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, hst)
}

// Rooms
func (h *Handler) ListRooms(w http.ResponseWriter, r *http.Request) {
	hostelID, _ := strconv.ParseUint(r.URL.Query().Get("hostel_id"), 10, 64)
	data, err := h.service.ListRooms(r.Context(), uint(hostelID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetRoom(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var rm domain.HostelRoom
	json.NewDecoder(r.Body).Decode(&rm)
	if err := h.service.CreateRoom(r.Context(), &rm); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, rm)
}
func (h *Handler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var rm domain.HostelRoom
	json.NewDecoder(r.Body).Decode(&rm)
	if err := h.service.UpdateRoom(r.Context(), uint(id), &rm); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, rm)
}

// Beds
func (h *Handler) ListBeds(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.ParseUint(r.URL.Query().Get("room_id"), 10, 64)
	data, err := h.service.ListBeds(r.Context(), uint(roomID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateBed(w http.ResponseWriter, r *http.Request) {
	var b domain.HostelBed
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.CreateBed(r.Context(), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, b)
}

// Allocations
func (h *Handler) ListAllocations(w http.ResponseWriter, r *http.Request) {
	hostelID, _ := strconv.ParseUint(r.URL.Query().Get("hostel_id"), 10, 64)
	studentID, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListAllocations(r.Context(), uint(hostelID), uint(studentID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AllocateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID uint  `json:"student_id"`
		RoomID    uint  `json:"room_id"`
		BedID     *uint `json:"bed_id,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	createdBy := middleware.GetUserID(r.Context())
	data, err := h.service.AllocateRoom(r.Context(), req.StudentID, req.RoomID, req.BedID, createdBy)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) DeallocateRoom(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.DeallocateRoom(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "room deallocated"})
}

// Mess Bills
func (h *Handler) ListMessBills(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListMessBills(r.Context(), uint(studentID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateMessBill(w http.ResponseWriter, r *http.Request) {
	var m domain.MessBill
	json.NewDecoder(r.Body).Decode(&m)
	if err := h.service.CreateMessBill(r.Context(), &m); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, m)
}
func (h *Handler) PayMessBill(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.PayMessBill(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "mess bill paid"})
}

// Maintenance Requests
func (h *Handler) ListMaintenanceRequests(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.ParseUint(r.URL.Query().Get("room_id"), 10, 64)
	studentID, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	data, err := h.service.ListMaintenanceRequests(r.Context(), uint(roomID), uint(studentID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateMaintenanceRequest(w http.ResponseWriter, r *http.Request) {
	var m domain.MaintenanceRequest
	json.NewDecoder(r.Body).Decode(&m)
	if err := h.service.CreateMaintenanceRequest(r.Context(), &m); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, m)
}
func (h *Handler) ResolveMaintenanceRequest(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	assignedTo := middleware.GetUserID(r.Context())
	if err := h.service.ResolveMaintenanceRequest(r.Context(), uint(id), assignedTo); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "maintenance request resolved"})
}

// Visitor Logs
func (h *Handler) ListVisitorLogs(w http.ResponseWriter, r *http.Request) {
	hostelID, _ := strconv.ParseUint(r.URL.Query().Get("hostel_id"), 10, 64)
	data, err := h.service.ListVisitorLogs(r.Context(), uint(hostelID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateVisitorLog(w http.ResponseWriter, r *http.Request) {
	var v domain.VisitorLog
	json.NewDecoder(r.Body).Decode(&v)
	if err := h.service.CreateVisitorLog(r.Context(), &v); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, v)
}
func (h *Handler) ExitVisitor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.ExitVisitor(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "visitor exit recorded"})
}
