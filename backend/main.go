package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fintrack/internal/config"
	"fintrack/internal/handler"
	"fintrack/internal/repository"
	"fintrack/internal/service"
	"fintrack/internal/websocket"
	"fintrack/pkg/cache"
	"fintrack/pkg/database"
	"fintrack/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	redisClient, err := cache.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	hub := websocket.NewHub()
	go hub.Run()

	// Wire up layers: repository → service → handler
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	accountService := service.NewAccountService(accountRepo, redisClient)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, redisClient, hub)

	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	wsHandler := handler.NewWebSocketHandler(hub)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://fin-track-ten-phi.vercel.app", "https://fin-track-chanakyasarmas-projects.vercel.app", "https://*.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Public routes
	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		r.Get("/api/v1/accounts", accountHandler.List)
		r.Post("/api/v1/accounts", accountHandler.Create)
		r.Get("/api/v1/accounts/{id}", accountHandler.Get)

		r.Get("/api/v1/transactions", transactionHandler.List)
		r.Post("/api/v1/transactions", transactionHandler.Create)
		r.Get("/api/v1/transactions/summary", transactionHandler.Summary)

		r.Get("/api/v1/ws", wsHandler.ServeWS)
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited cleanly")
}
