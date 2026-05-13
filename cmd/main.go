package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/historian/backend/internal/config"
	"github.com/historian/backend/internal/handlers"
	"github.com/historian/backend/internal/middleware"
	"github.com/historian/backend/internal/repository"
	"github.com/historian/backend/internal/services"
)

func main() {
	cfg := config.Load()

	pgPool, err := repository.NewPostgresPool(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pgPool.Close()

	chConn, err := repository.NewClickHouseConn(cfg.ClickHouseDSN)
	if err != nil {
		log.Fatalf("clickhouse: %v", err)
	}
	defer chConn.Close()

	kafkaProducer := repository.NewKafkaProducer(cfg.KafkaBrokers, "analytics-events")
	defer kafkaProducer.Close()

	userRepo := repository.NewUserRepo(pgPool)
	//programRepo := repository.NewProgramRepo(pgPool)
	testRepo := repository.NewTestRepo(pgPool)
	analyticsRepo := repository.NewAnalyticsRepo(chConn)

	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret)
	//programSvc := services.NewProgramService(programRepo)
	testSvc := services.NewTestService(testRepo, analyticsRepo, kafkaProducer)
	//analyticsSvc := services.NewAnalyticsService(analyticsRepo)

	mux := http.NewServeMux()

	authH := handlers.NewAuthHandler(authSvc)
	//programH := handlers.NewProgramHandler(programSvc)
	testH := handlers.NewTestHandler(testSvc)
	//analyticsH := handlers.NewAnalyticsHandler(analyticsSvc)

	// Auth
	mux.HandleFunc("POST /api/auth/register", authH.Register)
	mux.HandleFunc("POST /api/auth/login", authH.Login)
	mux.HandleFunc("POST /api/auth/logout", authH.Logout)
	mux.HandleFunc("GET /api/auth/me", middleware.Auth(authSvc, authH.Me))

	// Programs (teacher)
	//mux.HandleFunc("GET /api/programs", middleware.Auth(authSvc, programH.List))
	//mux.HandleFunc("GET /api/programs/{id}", middleware.Auth(authSvc, programH.Get))
	//mux.HandleFunc("POST /api/programs", middleware.Role(authSvc, "teacher", programH.Create))
	//mux.HandleFunc("PUT /api/programs/{id}", middleware.Role(authSvc, "teacher", programH.Update))
	//mux.HandleFunc("DELETE /api/programs/{id}", middleware.Role(authSvc, "teacher", programH.Delete))
	//mux.HandleFunc("GET /api/programs/{id}/report", middleware.Role(authSvc, "teacher", programH.GenerateReport))

	// Tests (student)
	mux.HandleFunc("GET /api/tests", middleware.Auth(authSvc, testH.List))
	mux.HandleFunc("GET /api/tests/{id}", middleware.Auth(authSvc, testH.Get))
	mux.HandleFunc("POST /api/tests/{id}/submit", middleware.Auth(authSvc, testH.Submit))
	mux.HandleFunc("GET /api/tests/results/my", middleware.Auth(authSvc, testH.MyResults))

	// Analytics
	//mux.HandleFunc("GET /api/analytics/test/{id}", middleware.Auth(authSvc, analyticsH.TestAnalytics))
	//mux.HandleFunc("GET /api/analytics/my", middleware.Auth(authSvc, analyticsH.MyAnalytics))

	handler := middleware.CORS(middleware.Logger(mux))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	go func() {
		log.Printf("Server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
