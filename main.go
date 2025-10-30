package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/err0r500/go-realworld-clean/auth"
	"github.com/err0r500/go-realworld-clean/driven/db"
	"github.com/err0r500/go-realworld-clean/driven/identityProvider"
	"github.com/err0r500/go-realworld-clean/driven/profileRepo"
	"github.com/err0r500/go-realworld-clean/driving/api/profile"
	"github.com/err0r500/go-realworld-clean/usecase"
)

func main() {
	// Initialize PostgreSQL connection
	database, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	if err := db.RunMigrations(database, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	profileRepository := profileRepo.NewPostgresProfileRepo(database)

	// Initialize Kratos clients
	kratosURL := "http://localhost:4433"
	kratosAdminURL := "http://localhost:4434"
	authHandler := auth.NewAuthHandlerWithAdmin(kratosURL, kratosAdminURL)

	// Initialize identity provider
	kratosIdentityProvider := identityProvider.NewKratosIdentityProvider(authHandler.GetKratosAdminClient())

	// Initialize use cases
	profileUseCase := usecase.NewProfileUseCases(profileRepository, kratosIdentityProvider)

	// Initialize HTTP handlers
	profileHandler := profile.NewProfileHandler(profileUseCase)

	// Setup HTTP mux and register endpoints
	mux := http.NewServeMux()

	// Register auth endpoints
	authHandler.RegisterUser(mux)
	authHandler.Login(mux)
	authHandler.GetCurrentUser(mux)
	authHandler.UpdateUser(mux)

	// Register profile endpoints
	profileHandler.GetProfile(mux)
	profileHandler.FollowUser(mux)
	profileHandler.UnfollowUser(mux)

	// Wrap with auth middleware
	handler := authHandler.AuthMiddleware(mux)

	// Start server
	addr := ":8080"
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
