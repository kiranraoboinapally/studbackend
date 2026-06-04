package authmod

import (
	"encoding/json"
	"net/http"

	"university-erp-backend/internal/modules/admissions"
	"university-erp-backend/internal/platform/apperrors"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/response"

	"github.com/gorilla/mux"
)

// Handler holds HTTP handlers for the auth module.
type Handler struct {
	service           *Service
	admissionsService *admissionsmod.Service
}

func NewHandler(service *Service, admissionsService *admissionsmod.Service) *Handler {
	return &Handler{service: service, admissionsService: admissionsService}
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
		MobileOTP      string `json:"mobile_otp"`
		EmailOTP       string `json:"email_otp"`
		MobileVerified bool   `json:"mobile_verified"`
		EmailVerified  bool   `json:"email_verified"`
		EnquiryID      uint   `json:"enquiry_id"`
		// Legacy fields for backward compatibility
		MobileNumber string `json:"mobile_number"`
		Email        string `json:"email"`
		OTP          string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	// If using new format with mobile_otp and email_otp, handle admissions verification
	if req.MobileOTP != "" || req.EmailOTP != "" {
		// This is an admissions verification request
		var enquiryID uint
		
		// Use enquiry_id if provided, otherwise find by mobile or email, otherwise use latest unverified
		if req.EnquiryID != 0 {
			enquiryID = req.EnquiryID
		} else if req.MobileNumber != "" {
			enquiry, err := h.admissionsService.GetEnquiryByMobile(r.Context(), req.MobileNumber)
			if err != nil {
				response.Error(w, apperrors.NotFound("enquiry not found"))
				return
			}
			enquiryID = enquiry.ID
		} else if req.Email != "" {
			enquiry, err := h.admissionsService.GetEnquiryByEmail(r.Context(), req.Email)
			if err != nil {
				response.Error(w, apperrors.NotFound("enquiry not found"))
				return
			}
			enquiryID = enquiry.ID
		} else {
			// If no identifier provided, use latest unverified enquiry
			enquiry, err := h.admissionsService.GetLatestUnverifiedEnquiry(r.Context())
			if err != nil {
				response.Error(w, apperrors.NotFound("no unverified enquiry found"))
				return
			}
			enquiryID = enquiry.ID
		}
		
		// Verify mobile OTP if provided
		if req.MobileOTP != "" {
			if err := h.admissionsService.VerifyOTP(r.Context(), enquiryID, req.MobileOTP, "mobile"); err != nil {
				response.Error(w, err)
				return
			}
		}
		
		// Verify email OTP if provided
		if req.EmailOTP != "" {
			if err := h.admissionsService.VerifyOTP(r.Context(), enquiryID, req.EmailOTP, "email"); err != nil {
				response.Error(w, err)
				return
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "OTP verified successfully",
		})
		return
	}

	// Legacy format for auth verification
	resp, err := h.service.VerifyOTP(r.Context(), req.MobileNumber, req.Email, req.OTP)
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
		Email        string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, err)
		return
	}

	resp, err := h.service.ResendOTP(r.Context(), req.MobileNumber, req.Email)
	if err != nil {
		response.Error(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"otp":     resp.OTP,
		"mobile_number": req.MobileNumber,
		"email":   req.Email,
	})
}
