package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/auth"
	"github.com/KaoriNakajima/sturdyticket/backend/internal/event"
	eventpg "github.com/KaoriNakajima/sturdyticket/backend/internal/event/postgres"
	"github.com/KaoriNakajima/sturdyticket/backend/pkg/response"
)

func main() {
	ctx := context.Background()

	// Connect to PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize Firebase Admin SDK
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	firebaseApp, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalf("failed to initialize firebase: %v", err)
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("failed to initialize firebase auth: %v", err)
	}

	authMiddleware := auth.NewMiddleware(authClient)

	// Initialize event domain
	eventRepo := eventpg.NewRepository(pool)
	eventUseCase := event.NewUseCase(eventRepo)
	eventHandler := event.NewHandler(eventUseCase)

	r := chi.NewRouter()

	// TODO: set up middleware (logging, recovery, CORS, rate limiting, bot detection)

	// Public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	eventHandler.RegisterRoutes(r)

	// Protected routes (require Firebase auth)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
			uid := auth.UserIDFromContext(r.Context())
			response.JSON(w, http.StatusOK, map[string]string{"uid": uid})
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
