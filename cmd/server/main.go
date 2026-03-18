package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/user/access-storage-server/internal/handler"
	"github.com/user/access-storage-server/internal/middleware"
	"github.com/user/access-storage-server/internal/repository"
	"github.com/user/access-storage-server/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	dsn := getEnv("DATABASE_DSN", "root:password@tcp(127.0.0.1:3306)/access_storage?parseTime=true")
	jwtSecret := getEnv("JWT_SECRET", "change-me-in-production")
	listenAddr := getEnv("LISTEN_ADDR", ":8080")

	db, err := repository.NewDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	seedAdminIfEmpty(userRepo)

	projectRepo := repository.NewProjectRepository(db)
	shareRepo := repository.NewShareRepository(db)

	authService := service.NewAuthService(userRepo, jwtSecret)
	projectService := service.NewProjectService(projectRepo, shareRepo, userRepo)
	adminService := service.NewAdminService(userRepo, shareRepo)

	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	adminHandler := handler.NewAdminHandler(adminService)

	authMiddleware := middleware.JWTAuth(jwtSecret)
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	// Auth routes
	mux.Handle("GET /api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("PUT /api/users/me", authMiddleware(http.HandlerFunc(authHandler.UpdateProfile)))

	// Protected routes
	mux.Handle("GET /api/projects/meta", authMiddleware(http.HandlerFunc(projectHandler.ListMeta)))
	mux.Handle("GET /api/projects", authMiddleware(http.HandlerFunc(projectHandler.List)))
	mux.Handle("GET /api/projects/{id}", authMiddleware(http.HandlerFunc(projectHandler.Get)))
	mux.Handle("POST /api/projects", authMiddleware(http.HandlerFunc(projectHandler.Create)))
	mux.Handle("PUT /api/projects/{id}", authMiddleware(http.HandlerFunc(projectHandler.Update)))
	mux.Handle("DELETE /api/projects/{id}", authMiddleware(http.HandlerFunc(projectHandler.Delete)))

	// Admin routes
	mux.Handle("GET /api/admin/users", authMiddleware(http.HandlerFunc(adminHandler.ListUsers)))
	mux.Handle("POST /api/admin/users", authMiddleware(http.HandlerFunc(adminHandler.CreateUser)))
	mux.Handle("PUT /api/admin/users/{id}", authMiddleware(http.HandlerFunc(adminHandler.UpdateUser)))
	mux.Handle("DELETE /api/admin/users/{id}", authMiddleware(http.HandlerFunc(adminHandler.DeleteUser)))
	mux.Handle("GET /api/admin/users/{id}/shares", authMiddleware(http.HandlerFunc(adminHandler.ListUserShares)))

	// Project sharing routes
	mux.Handle("POST /api/projects/{id}/share", authMiddleware(http.HandlerFunc(adminHandler.ShareProject)))
	mux.Handle("DELETE /api/projects/{id}/share/{userId}", authMiddleware(http.HandlerFunc(adminHandler.UnshareProject)))

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	finalHandler := corsMiddleware(rateLimiter.Middleware(mux))

	log.Printf("Server starting on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, finalHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func seedAdminIfEmpty(userRepo *repository.UserRepository) {
	count, err := userRepo.Count()
	if err != nil {
		log.Printf("[SEED] Warning: could not check user count: %v", err)
		return
	}
	if count > 0 {
		return
	}

	email := getEnv("SEED_ADMIN_EMAIL", "admin@local.local")

	passwordBytes := make([]byte, 8)
	if _, err := rand.Read(passwordBytes); err != nil {
		log.Fatalf("[SEED] Failed to generate password: %v", err)
	}
	password := hex.EncodeToString(passwordBytes)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[SEED] Failed to hash password: %v", err)
	}

	_, err = userRepo.CreateWithAdmin(email, string(hash), true)
	if err != nil {
		log.Fatalf("[SEED] Failed to create admin user: %v", err)
	}

	log.Printf("[SEED] ========================================")
	log.Printf("[SEED] Admin user created!")
	log.Printf("[SEED]   Email:    %s", email)
	log.Printf("[SEED]   Password: %s", password)
	log.Printf("[SEED] ========================================")
}
