package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"university-erp-backend/internal/config"
	academicmod "university-erp-backend/internal/modules/academic"
	admissionsmod "university-erp-backend/internal/modules/admissions"
	authmod "university-erp-backend/internal/modules/auth"
	coremod "university-erp-backend/internal/modules/core"
	exammod "university-erp-backend/internal/modules/exam"
	financemod "university-erp-backend/internal/modules/finance"
	hostelmod "university-erp-backend/internal/modules/hostel"
	hrmod "university-erp-backend/internal/modules/hr"
	librarymod "university-erp-backend/internal/modules/library"
	studentmod "university-erp-backend/internal/modules/student"
	transportmod "university-erp-backend/internal/modules/transport"
	"university-erp-backend/internal/platform/auth"
	"university-erp-backend/internal/platform/database"
	"university-erp-backend/internal/platform/eventbus"
	"university-erp-backend/internal/platform/middleware"
	"university-erp-backend/internal/platform/outbox"
	"university-erp-backend/internal/platform/swagger"
	wsplatform "university-erp-backend/internal/platform/websocket"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system env")
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

	// ─── Initialize Repositories ────────────────────────────────────────────────
	authRepo := authmod.NewRepository(db)
	coreRepo := coremod.NewRepository(db)
	studentRepo := studentmod.NewRepository(db)
	admissionsRepo := admissionsmod.NewRepository(db)
	academicRepo := academicmod.NewRepository(db)
	financeRepo := financemod.NewRepository(db)
	hrRepo := hrmod.NewRepository(db)
	examRepo := exammod.NewRepository(db)
	libraryRepo := librarymod.NewRepository(db)
	hostelRepo := hostelmod.NewRepository(db)
	transportRepo := transportmod.NewRepository(db)

	// ─── Initialize Services ────────────────────────────────────────────────────
	// Order matters: Finance & Student subscribe to events, so they must be
	// initialized before publishers that emit those events start up.
	// Since all subscriptions happen in NewService constructors, and the
	// outbox worker hasn't started yet, initialization order is safe.

	coreService := coremod.NewService(coreRepo)
	authService := authmod.NewService(authRepo, jwtMgr, bus, outboxWriter, db)
	studentService := studentmod.NewService(studentRepo, bus, outboxWriter, db)
	admissionsService := admissionsmod.NewService(admissionsRepo, bus, outboxWriter, db)
	academicService := academicmod.NewService(academicRepo, bus, outboxWriter, db)
	financeService := financemod.NewService(financeRepo, bus, outboxWriter, db)
	hrService := hrmod.NewService(hrRepo, bus, outboxWriter, db)
	examService := exammod.NewService(examRepo, bus, outboxWriter, db)
	libraryService := librarymod.NewService(libraryRepo, bus, outboxWriter, db)
	hostelService := hostelmod.NewService(hostelRepo, bus, outboxWriter, db)
	transportService := transportmod.NewService(transportRepo, bus, outboxWriter, db)

	// ─── Initialize Handlers ───────────────────────────────────────────────────
	coreHandler := coremod.NewHandler(coreService)
	authHandler := authmod.NewHandler(authService, admissionsService)
	studentHandler := studentmod.NewHandler(studentService)
	admissionsHandler := admissionsmod.NewHandler(admissionsService)
	academicHandler := academicmod.NewHandler(academicService)
	financeHandler := financemod.NewHandler(financeService)
	hrHandler := hrmod.NewHandler(hrService)
	examHandler := exammod.NewHandler(examService)
	libraryHandler := librarymod.NewHandler(libraryService)
	hostelHandler := hostelmod.NewHandler(hostelService)
	transportHandler := transportmod.NewHandler(transportService)

	// ─── Setup WebSocket Hub ────────────────────────────────────────────────────
	wsHub := wsplatform.NewHub()
	go wsHub.Run()
	wsHandler := wsplatform.NewHandler(wsHub, jwtMgr)

	eventsToSubscribe := []string{
		eventbus.EventUserRegistered,
		eventbus.EventUserLoggedIn,
		eventbus.EventStudentEnrolled,
		eventbus.EventApplicationSubmitted,
		eventbus.EventApplicationApproved,
		eventbus.EventSeatAllocated,
		eventbus.EventInvoiceGenerated,
		eventbus.EventPaymentCompleted,
		eventbus.EventPaymentFailed,
		eventbus.EventRefundProcessed,
		eventbus.EventResultPublished,
		eventbus.EventHostelAllocated,
		eventbus.EventMaintenanceRequested,
		eventbus.EventNotificationCreated,
	}

	for _, evtType := range eventsToSubscribe {
		wsHubForClosure := wsHub
		bus.Subscribe(evtType, func(ctx context.Context, ev eventbus.Event) error {
			frontendType := ev.Type
			if ev.Type == eventbus.EventPaymentCompleted {
				frontendType = "finance.payment_received"
			}
			wsHubForClosure.Broadcast(frontendType, ev.Payload)
			return nil
		})
	}

	// ─── Setup Router ───────────────────────────────────────────────────────────
	r := mux.NewRouter()

	r.Handle("/ws", wsHandler)

	// Global middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS([]string{"*"}))

	// Auth middleware
	authMW := middleware.Authenticate(jwtMgr)

	// Register all module routes
	coreHandler.RegisterRoutes(r, authMW)
	authHandler.RegisterRoutes(r, authMW)
	studentHandler.RegisterRoutes(r, authMW)
	admissionsHandler.RegisterRoutes(r, authMW)
	academicHandler.RegisterRoutes(r, authMW)
	financeHandler.RegisterRoutes(r, authMW)
	hrHandler.RegisterRoutes(r, authMW)
	examHandler.RegisterRoutes(r, authMW)
	libraryHandler.RegisterRoutes(r, authMW)
	hostelHandler.RegisterRoutes(r, authMW)
	transportHandler.RegisterRoutes(r, authMW)

	// ─── Swagger UI ─────────────────────────────────────────────────────────────
	swaggerHandler := swagger.NewHandler("docs/swagger.yaml")
	swaggerMux := http.NewServeMux()
	swaggerHandler.RegisterRoutes(swaggerMux)
	r.PathPrefix("/swagger").Handler(swaggerMux)

	// ─── Start Outbox Worker ────────────────────────────────────────────────────
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		outboxWorker.Start(ctx)
	}()

	slog.Info("transactional outbox worker started",
		"poll_interval_sec", cfg.OutboxPollInterval,
		"batch_size", cfg.OutboxBatchSize,
	)

	// ─── Start HTTP Server ──────────────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("university ERP backend starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// ─── Graceful Shutdown ──────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	slog.Info("shutdown signal received", "signal", sig.String())

	// Stop accepting new HTTP requests
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced shutdown", "error", err)
	}

	// Stop outbox worker
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	slog.Info("server stopped gracefully")
}
