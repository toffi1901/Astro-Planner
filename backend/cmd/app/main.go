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
	profileuc "astro-planner/backend/internal/usecases/profile"
	satuc "astro-planner/backend/internal/usecases/satellite"

	"github.com/go-chi/cors"
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
	satRepo := repository.NewSatelliteRepository(pool)

	registerUC := authuc.NewRegisterUseCase(userRepo, cfg.Pepper)
	loginUC := authuc.NewLoginUseCase(userRepo, sessionRepo, limiterRepo, cfg.Pepper, cfg.JWTSecret)
	logoutUC := authuc.NewLogoutUseCase(sessionRepo)
	refreshUC := authuc.NewRefreshUseCase(sessionRepo, cfg.JWTSecret)

	profileUC := profileuc.NewProfileUseCase(userRepo)
	predictUC := satuc.NewPredictUseCase(satRepo, userRepo)

	authHandler := &api.AuthHandler{
		RegisterUC: registerUC,
		LoginUC:    loginUC,
		LogoutUC:   logoutUC,
		RefreshUC:  refreshUC,
	}
	profileHandler := &api.ProfileHandler{
		ProfileUC: profileUC,
	}
	satHandler := &api.SatelliteHandler{
		Repo:      satRepo,
		PredictUC: predictUC,
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// Перечисляем порты, на которых будет открываться твой фронтенд
		// 5500 - Live Server в VS Code, 5173 - Vite, 3000 - React
		AllowedOrigins:   []string{"http://localhost:5500", "http://127.0.0.1:5500", "http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // для рефреш токен
		MaxAge:           300,
	}))
	r.Use(api.Logging)

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/refresh", authHandler.Refresh)
	r.Get("/api/satellites", satHandler.List)

	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware(sessionRepo, cfg.JWTSecret))

		r.Post("/api/auth/logout", authHandler.Logout)
		r.Get("/api/profile", profileHandler.Get)
		r.Put("/api/profile", profileHandler.UpdateLocation)

		r.Post("/api/satellites/{norad_id}/predict", satHandler.Predict)
	})

	log.Println("Server on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
