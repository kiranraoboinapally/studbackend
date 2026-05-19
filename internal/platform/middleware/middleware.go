package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"university-erp-backend/internal/domain"
	"university-erp-backend/internal/platform/apperrors"
	"university-erp-backend/internal/platform/auth"

	"gorm.io/gorm"
)

type contextKey string

const (
	ContextUserID   contextKey = "user_id"
	ContextUsername  contextKey = "username"
	ContextRoles    contextKey = "roles"
)

// GetUserID extracts the authenticated user ID from context.
func GetUserID(ctx context.Context) uint {
	if v, ok := ctx.Value(ContextUserID).(uint); ok {
		return v
	}
	return 0
}

// GetRoles extracts roles from context.
func GetRoles(ctx context.Context) []string {
	if v, ok := ctx.Value(ContextRoles).([]string); ok {
		return v
	}
	return nil
}

// respondError writes a JSON error response.
func respondError(w http.ResponseWriter, err *apperrors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   err.Message,
		"detail":  err.Detail,
	})
}

// ─── CORS ────────────────────────────────────────────────────────────────────

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ─── Authentication ──────────────────────────────────────────────────────────

func Authenticate(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				respondError(w, apperrors.Unauthorized("missing authorization header"))
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				respondError(w, apperrors.Unauthorized("invalid authorization format"))
				return
			}

			claims, err := jwtMgr.ValidateToken(parts[1])
			if err != nil {
				respondError(w, apperrors.Unauthorized("invalid or expired token"))
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextUsername, claims.Username)
			ctx = context.WithValue(ctx, ContextRoles, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── Role Authorization ──────────────────────────────────────────────────────

func RequireRoles(allowed ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := GetRoles(r.Context())
			for _, userRole := range roles {
				for _, a := range allowed {
					if userRole == a {
						next.ServeHTTP(w, r)
						return
					}
				}
			}
			respondError(w, apperrors.Forbidden("insufficient permissions"))
		})
	}
}

// ─── Audit Logging ───────────────────────────────────────────────────────────

func AuditLog(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			// Only audit write operations
			if r.Method == "GET" || r.Method == "OPTIONS" {
				return
			}

			userID := GetUserID(r.Context())
			var uid *uint
			if userID > 0 {
				uid = &userID
			}

			entry := domain.AuditLog{
				UserID:    uid,
				Action:    r.Method,
				SchemaName: "",
				AffectedTable: r.URL.Path,
				IPAddress: r.RemoteAddr,
				UserAgent: r.UserAgent(),
				CreatedAt: start,
			}

			if err := db.Create(&entry).Error; err != nil {
				log.Printf("⚠️  Audit log write failed: %v", err)
			}
		})
	}
}

// ─── Request Logger ──────────────────────────────────────────────────────────

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%-6s %-40s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
