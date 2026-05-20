package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
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
	ContextUserID    contextKey = "user_id"
	ContextUsername   contextKey = "username"
	ContextRoles     contextKey = "roles"
	ContextTraceID   contextKey = "trace_id"
	ContextScopeDept contextKey = "scope_department_id"
)

func GetUserID(ctx context.Context) uint {
	if v, ok := ctx.Value(ContextUserID).(uint); ok {
		return v
	}
	return 0
}

func GetUsername(ctx context.Context) string {
	if v, ok := ctx.Value(ContextUsername).(string); ok {
		return v
	}
	return ""
}

func GetRoles(ctx context.Context) []string {
	if v, ok := ctx.Value(ContextRoles).([]string); ok {
		return v
	}
	return nil
}

func GetTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextTraceID).(string); ok {
		return v
	}
	return ""
}

func GetScopeDepartmentID(ctx context.Context) *uint {
	if v, ok := ctx.Value(ContextScopeDept).(*uint); ok {
		return v
	}
	return nil
}

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
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Trace-ID")
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

// ─── Scoped Role-Based Access Control ─────────────────────────────────────────

// RequirePermission checks that the user has a specific resource:action permission.
// It queries security.role_permissions joined with shared.user_roles.
func RequirePermission(db *gorm.DB, resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r.Context())
			if userID == 0 {
				respondError(w, apperrors.Unauthorized("not authenticated"))
				return
			}

			var count int64
			db.Table("security.role_permissions rp").
				Joins("JOIN shared.user_roles ur ON ur.role_id = rp.role_id").
				Joins("JOIN security.permissions p ON p.id = rp.permission_id").
				Where("ur.user_id = ? AND p.resource = ? AND p.action = ?", userID, resource, action).
				Count(&count)

			if count == 0 {
				slog.Warn("permission denied",
					"trace_id", GetTraceID(r.Context()),
					"user_id", userID,
					"resource", resource,
					"action", action,
				)
				respondError(w, apperrors.Forbidden("insufficient permissions for " + resource + ":" + action))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRoles checks that the user has at least one of the allowed roles.
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

// InjectScope loads the user's department scope from hr.employees and injects it into context.
// This enables scoped RBAC: a Department Head can only operate within their department.
func InjectScope(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r.Context())
			if userID == 0 {
				next.ServeHTTP(w, r)
				return
			}

			var deptID *uint
			db.Table("hr.employees").
				Select("department_id").
				Where("user_id = ? AND is_active = true", userID).
				Scan(&deptID)

			ctx := context.WithValue(r.Context(), ContextScopeDept, deptID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ─── Audit Logging ───────────────────────────────────────────────────────────

func AuditLog(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)

			if r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "HEAD" {
				return
			}

			userID := GetUserID(r.Context())
			var uid *uint
			if userID > 0 {
				uid = &userID
			}

			entry := domain.AuditLog{
				UserID:        uid,
				Action:        r.Method,
				SchemaName:    "",
				AffectedTable: r.URL.Path,
				RecordID:      "",
				IPAddress:     r.RemoteAddr,
				UserAgent:     r.UserAgent(),
				CreatedAt:     start,
			}

			if err := db.Create(&entry).Error; err != nil {
				slog.Error("audit log write failed", "error", err, "trace_id", GetTraceID(r.Context()))
			}
		})
	}
}

// ─── Request Logger with Trace ID ─────────────────────────────────────────────

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := generateTraceID()
		ctx := context.WithValue(r.Context(), ContextTraceID, traceID)
		r = r.WithContext(ctx)

		start := time.Now()
		next.ServeHTTP(w, r)

		slog.Info("request",
			"trace_id", traceID,
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
			"user_id", GetUserID(ctx),
		)
	})
}

func generateTraceID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 12)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
