package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"

	"github.com/KaoriNakajima/sturdyticket/backend/internal/auth"
	"github.com/KaoriNakajima/sturdyticket/backend/internal/booking"
	bookingpg "github.com/KaoriNakajima/sturdyticket/backend/internal/booking/postgres"
	"github.com/KaoriNakajima/sturdyticket/backend/internal/event"
	eventpg "github.com/KaoriNakajima/sturdyticket/backend/internal/event/postgres"
	"github.com/KaoriNakajima/sturdyticket/backend/internal/session"
	sessionredis "github.com/KaoriNakajima/sturdyticket/backend/internal/session/redis"
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

	// Initialize Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	redisOpts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse REDIS_URL: %v", err)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// Initialize session domain
	sessionStore := sessionredis.NewStore(redisClient)
	sessionService := session.NewService(sessionStore, 30*time.Second, 100)
	sessionHandler := session.NewHandler(sessionService)

	// Start session expiry subscriber
	go sessionredis.StartSubscriber(ctx, redisClient)

	// Initialize event domain
	eventRepo := eventpg.NewRepository(pool)
	eventUseCase := event.NewUseCase(eventRepo, sessionService)
	eventHandler := event.NewHandler(eventUseCase)

	// Initialize booking domain
	bookingRepo := bookingpg.NewRepository(pool)
	bookingUseCase := booking.NewUseCase(bookingRepo, eventRepo)
	bookingHandler := booking.NewHandler(bookingUseCase)

	r := chi.NewRouter()

	// TODO: set up middleware (logging, recovery, CORS, rate limiting, bot detection)

	// Public routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	eventHandler.RegisterPublicRoutes(r)

	// Protected routes (require Firebase auth)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Get("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
			uid := auth.UserIDFromContext(r.Context())
			response.JSON(w, http.StatusOK, map[string]string{"uid": uid})
		})

		eventHandler.RegisterProtectedRoutes(r)
		bookingHandler.RegisterRoutes(r)
		sessionHandler.RegisterRoutes(r)
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
