package authmod

import (
	"encoding/json"
	"net/http"

	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

// Handler holds HTTP handlers for the auth module.
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes mounts auth routes on the given router.
func (h *Handler) RegisterRoutes(r *mux.Router, authMW mux.MiddlewareFunc) {
	// Public routes
	r.HandleFunc("/api/v1/auth/login", h.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/auth/register", h.Register).Methods("POST", "OPTIONS")

	// Protected routes
	protected := r.PathPrefix("/api/v1/auth").Subrouter()
	protected.Use(authMW)
	protected.HandleFunc("/profile", h.Profile).Methods("GET")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	resp, err := h.service.Login(r.Context(), req, r.RemoteAddr, r.UserAgent())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, resp)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	resp, err := h.service.Register(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, resp)
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	resp, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, resp)
}
