package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"astro-planner/backend/internal/api"
	"astro-planner/backend/internal/config"
	"astro-planner/backend/internal/repository"
	"astro-planner/backend/internal/storage"
	authuc "astro-planner/backend/internal/usecases/auth"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := storage.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("postgres: ", err)
	}
	defer pool.Close()

	rdb, err := storage.NewRedis(ctx, cfg.RedisAddr)
	if err != nil {
		log.Fatal("redis: ", err)
	}
	defer rdb.Close()

	userRepo := repository.NewPostgresUserRepository(pool)
	sessionRepo := repository.NewSessionRepository(rdb)
	limiterRepo := repository.NewLimiterRepository(rdb)

	registerUC := authuc.NewRegisterUseCase(userRepo, cfg.Pepper)
	loginUC := authuc.NewLoginUseCase(userRepo, sessionRepo, limiterRepo, cfg.Pepper, cfg.JWTSecret)
	logoutUC := authuc.NewLogoutUseCase(sessionRepo)

	authHandler := &api.AuthHandler{
		RegisterUC: registerUC,
		LoginUC:    loginUC,
		LogoutUC:   logoutUC,
	}

	r := chi.NewRouter()
	r.Use(api.Logging)

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware(sessionRepo, cfg.JWTSecret))
		r.Post("/api/auth/logout", authHandler.Logout)
	})

	log.Println("Server on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
