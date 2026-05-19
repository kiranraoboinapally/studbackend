package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"university-erp-backend/internal/config"
	authmod "university-erp-backend/internal/modules/auth"
	studentmod "university-erp-backend/internal/modules/student"
	"university-erp-backend/internal/platform/auth"
	"university-erp-backend/internal/platform/database"
	"university-erp-backend/internal/platform/eventbus"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/outbox"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	// Load config
	cfg := config.Load()

	// Connect DB
	db := database.Connect(cfg)

	// Initialize platform components
	jwtMgr := auth.NewJWTManager(cfg)
	bus := eventbus.New()
	outboxWriter := outbox.NewWriter(db)
	outboxWorker := outbox.NewWorker(db, bus, cfg.OutboxPollInterval, cfg.OutboxBatchSize)

	// Start outbox worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	outboxWorker.Start(ctx)

	// Initialize repositories
	authRepo := authmod.NewRepository(db)
	studentRepo := studentmod.NewRepository(db)

	// Initialize services
	authService := authmod.NewService(authRepo, jwtMgr, bus, outboxWriter, db)
	studentService := studentmod.NewService(studentRepo, bus, outboxWriter, db)

	// Initialize handlers
	authHandler := authmod.NewHandler(authService)
	studentHandler := studentmod.NewHandler(studentService)

	// Setup Router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS([]string{"*"}))

	// Create auth middleware
	authMW := middleware.Authenticate(jwtMgr)

	// Register routes
	authHandler.RegisterRoutes(r, authMW)
	studentHandler.RegisterRoutes(r, authMW)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 University ERP Backend running on :%s", port)
		log.Println("📚 API Base: http://localhost:" + port + "/api/v1")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	// Stop outbox worker
	cancel()

	log.Println("✅ Server stopped")
}
