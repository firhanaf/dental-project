package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/yourusername/dental-api/config"
	"github.com/yourusername/dental-api/internal/handler"
	"github.com/yourusername/dental-api/internal/middleware"
	"github.com/yourusername/dental-api/internal/repository"
	"github.com/yourusername/dental-api/internal/service"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"github.com/yourusername/dental-api/pkg/storage"
)

func main() {
	// Load .env (diabaikan jika tidak ada — untuk production pakai env langsung)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// ── Database ─────────────────────────────────────────────
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}
	log.Println("✓ Database connected")

	// ── Dependencies ─────────────────────────────────────────
	jwtMgr  := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)
	store   := storage.NewLocal(cfg.UploadDir, cfg.MaxFileSizeMB)
	authMw  := middleware.NewAuth(jwtMgr)

	// Repositories
	userRepo       := repository.NewUserRepo(pool)
	branchRepo     := repository.NewBranchRepo(pool)
	patientRepo    := repository.NewPatientRepo(pool)
	visitRepo      := repository.NewVisitRepo(pool)
	attachmentRepo := repository.NewAttachmentRepo(pool)

	// Services
	authSvc       := service.NewAuthService(userRepo, jwtMgr)
	branchSvc     := service.NewBranchService(branchRepo)
	patientSvc    := service.NewPatientService(patientRepo)
	visitSvc      := service.NewVisitService(visitRepo, patientRepo)
	attachmentSvc := service.NewAttachmentService(attachmentRepo, store)
	exportSvc     := service.NewExportService(patientRepo, visitRepo)
	userMgmtSvc   := service.NewUserMgmtService(userRepo)

	// Handlers
	authH       := handler.NewAuthHandler(authSvc)
	branchH     := handler.NewBranchHandler(branchSvc)
	patientH    := handler.NewPatientHandler(patientSvc)
	visitH      := handler.NewVisitHandler(visitSvc)
	attachmentH := handler.NewAttachmentHandler(attachmentSvc)
	exportH     := handler.NewExportHandler(exportSvc)
	userMgmtH   := handler.NewUserMgmtHandler(userMgmtSvc)

	// ── Router ───────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check (tidak perlu auth)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// ── API Routes ───────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {

		// Public
		r.Post("/auth/login",  authH.Login)

		// Protected — semua perlu login
		r.Group(func(r chi.Router) {
			r.Use(authMw.Authenticate)

			// Auth
			r.Get("/auth/me",     authH.Me)
			r.Post("/auth/logout", authH.Logout)

			// Branches
			r.Get("/branches", branchH.List)

			// Patients — READ (write & readonly)
			r.Get("/patients",          patientH.List)
			r.Get("/patients/{id}",     patientH.GetByID)
			r.Get("/patients/{id}/visits",     visitH.ListByPatient)
			r.Get("/patients/{id}/attachments", attachmentH.ListByPatient)

			// Visits — READ
			r.Get("/visits/{id}", visitH.GetByID)

			// Attachments — READ & DOWNLOAD
			r.Get("/attachments/{id}",          attachmentH.GetByID)
			r.Get("/attachments/{id}/download", attachmentH.Download)

			// Export — semua role bisa export
			r.Get("/export/patients", exportH.ExportPatients)
			r.Get("/export/visits",   exportH.ExportVisits)

			// ── Write-only routes ────────────────────────────
			r.Group(func(r chi.Router) {
				r.Use(authMw.RequireWrite)

				// Patients
				r.Post("/patients",        patientH.Create)
				r.Put("/patients/{id}",    patientH.Update)
				r.Delete("/patients/{id}", patientH.Delete)

				// Visits
				r.Post("/visits",        visitH.Create)
				r.Put("/visits/{id}",    visitH.Update)
				r.Delete("/visits/{id}", visitH.Delete)

				// Attachments
				r.Post("/attachments",        attachmentH.Upload)
				r.Delete("/attachments/{id}", attachmentH.Delete)
			})

			// ── SuperAdmin-only routes ───────────────────────
			r.Group(func(r chi.Router) {
				r.Use(authMw.RequireSuperAdmin)

				r.Get("/users",          userMgmtH.List)
				r.Post("/users",         userMgmtH.Create)
				r.Put("/users/{id}",     userMgmtH.Update)
				r.Delete("/users/{id}",  userMgmtH.Deactivate)
			})
		})
	})

	// ── Start Server ─────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second, // lebih lama untuk export file besar
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("✓ Server running on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
	log.Println("Server stopped.")
}
