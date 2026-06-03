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
	r.HandleFunc("/api/v1/auth/verify-otp", h.VerifyOTP).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/auth/resend-otp", h.ResendOTP).Methods("POST", "OPTIONS")

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

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MobileNumber string `json:"mobile_number"`
		OTP          string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	resp, err := h.service.VerifyOTP(r.Context(), req.MobileNumber, req.OTP)
	if err != nil {
		response.Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

func (h *Handler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MobileNumber string `json:"mobile_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	resp, err := h.service.ResendOTP(r.Context(), req.MobileNumber)
	if err != nil {
		response.Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"otp":     resp.OTP,
		"data":    resp,
	})
}
