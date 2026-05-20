package transportmod

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
	api := r.PathPrefix("/api/v1/transport").Subrouter()
	api.Use(authMW)

	// Stats
	api.HandleFunc("/stats", h.Stats).Methods("GET")

	// Buses
	api.HandleFunc("/buses", h.ListBuses).Methods("GET")
	api.HandleFunc("/buses", h.CreateBus).Methods("POST")
	api.HandleFunc("/buses/{id:[0-9]+}", h.GetBus).Methods("GET")
	api.HandleFunc("/buses/{id:[0-9]+}", h.UpdateBus).Methods("PUT")

	// Routes
	api.HandleFunc("/routes", h.ListRoutes).Methods("GET")
	api.HandleFunc("/routes", h.CreateRoute).Methods("POST")
	api.HandleFunc("/routes/{id:[0-9]+}", h.GetRoute).Methods("GET")
	api.HandleFunc("/routes/{id:[0-9]+}", h.UpdateRoute).Methods("PUT")

	// Stops
	api.HandleFunc("/stops", h.ListStops).Methods("GET")
	api.HandleFunc("/stops", h.CreateStop).Methods("POST")
	api.HandleFunc("/stops/{id:[0-9]+}", h.UpdateStop).Methods("PUT")

	// Bus Assignments
	api.HandleFunc("/bus-assignments", h.ListBusAssignments).Methods("GET")
	api.HandleFunc("/bus-assignments", h.AssignBusToRoute).Methods("POST")
	api.HandleFunc("/bus-assignments/{id:[0-9]+}/end", h.EndBusAssignment).Methods("POST")

	// Student Passes
	api.HandleFunc("/passes", h.ListStudentPasses).Methods("GET")
	api.HandleFunc("/passes/issue", h.IssuePass).Methods("POST")
	api.HandleFunc("/passes/{id:[0-9]+}/renew", h.RenewPass).Methods("POST")

	// Vehicle Maintenance
	api.HandleFunc("/maintenance", h.ListVehicleMaintenance).Methods("GET")
	api.HandleFunc("/maintenance", h.CreateVehicleMaintenance).Methods("POST")
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetStats(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}

// Buses
func (h *Handler) ListBuses(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"
	data, err := h.service.ListBuses(r.Context(), activeOnly)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetBus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetBus(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateBus(w http.ResponseWriter, r *http.Request) {
	var b domain.Bus
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.CreateBus(r.Context(), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, b)
}
func (h *Handler) UpdateBus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var b domain.Bus
	json.NewDecoder(r.Body).Decode(&b)
	if err := h.service.UpdateBus(r.Context(), uint(id), &b); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, b)
}

// Routes
func (h *Handler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	activeOnly := r.URL.Query().Get("active") == "true"
	data, err := h.service.ListRoutes(r.Context(), activeOnly)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	data, err := h.service.GetRoute(r.Context(), uint(id))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var rt domain.Route
	json.NewDecoder(r.Body).Decode(&rt)
	if err := h.service.CreateRoute(r.Context(), &rt); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, rt)
}
func (h *Handler) UpdateRoute(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var rt domain.Route
	json.NewDecoder(r.Body).Decode(&rt)
	if err := h.service.UpdateRoute(r.Context(), uint(id), &rt); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, rt)
}

// Stops
func (h *Handler) ListStops(w http.ResponseWriter, r *http.Request) {
	routeID, _ := strconv.ParseUint(r.URL.Query().Get("route_id"), 10, 64)
	data, err := h.service.ListStops(r.Context(), uint(routeID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateStop(w http.ResponseWriter, r *http.Request) {
	var st domain.Stop
	json.NewDecoder(r.Body).Decode(&st)
	if err := h.service.CreateStop(r.Context(), &st); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, st)
}
func (h *Handler) UpdateStop(w http.ResponseWriter, r *http.Request) {
	var st domain.Stop
	json.NewDecoder(r.Body).Decode(&st)
	if err := h.service.UpdateStop(r.Context(), &st); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, st)
}

// Bus Assignments
func (h *Handler) ListBusAssignments(w http.ResponseWriter, r *http.Request) {
	busID, _ := strconv.ParseUint(r.URL.Query().Get("bus_id"), 10, 64)
	routeID, _ := strconv.ParseUint(r.URL.Query().Get("route_id"), 10, 64)
	data, err := h.service.ListBusAssignments(r.Context(), uint(busID), uint(routeID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) AssignBusToRoute(w http.ResponseWriter, r *http.Request) {
	var ba domain.BusAssignment
	json.NewDecoder(r.Body).Decode(&ba)
	if err := h.service.AssignBusToRoute(r.Context(), &ba); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, ba)
}
func (h *Handler) EndBusAssignment(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err := h.service.EndBusAssignment(r.Context(), uint(id)); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "bus assignment ended"})
}

// Student Passes
func (h *Handler) ListStudentPasses(w http.ResponseWriter, r *http.Request) {
	studentID, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	routeID, _ := strconv.ParseUint(r.URL.Query().Get("route_id"), 10, 64)
	data, err := h.service.ListStudentPasses(r.Context(), uint(studentID), uint(routeID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) IssuePass(w http.ResponseWriter, r *http.Request) {
	var sp domain.StudentPass
	json.NewDecoder(r.Body).Decode(&sp)
	createdBy := middleware.GetUserID(r.Context())
	data, err := h.service.IssuePass(r.Context(), &sp, createdBy)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, data)
}
func (h *Handler) RenewPass(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	var req struct {
		ValidTo  string  `json:"valid_to"`
		FeePaid  float64 `json:"fee_paid"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	validTo, _ := time.Parse("2006-01-02", req.ValidTo)
	if err := h.service.RenewPass(r.Context(), uint(id), validTo, req.FeePaid); err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "pass renewed"})
}

// Vehicle Maintenance
func (h *Handler) ListVehicleMaintenance(w http.ResponseWriter, r *http.Request) {
	busID, _ := strconv.ParseUint(r.URL.Query().Get("bus_id"), 10, 64)
	data, err := h.service.ListVehicleMaintenance(r.Context(), uint(busID))
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, data)
}
func (h *Handler) CreateVehicleMaintenance(w http.ResponseWriter, r *http.Request) {
	var vm domain.VehicleMaintenance
	json.NewDecoder(r.Body).Decode(&vm)
	if err := h.service.CreateVehicleMaintenance(r.Context(), &vm); err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, vm)
}
