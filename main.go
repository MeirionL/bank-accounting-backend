package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MeirionL/personal-finance-app/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Imported database driver we need
)

type apiConfig struct {
	DB        *database.Queries // Calling the database package in my directory
	jwtSecret string
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL) // Loading in our database.
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	cfg := apiConfig{ // An API config we can pass to our
		DB:        database.New(db), // handlers so that they have access to
		jwtSecret: jwtSecret,        // our database.
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", cfg.handlerCreateUser)
	v1Router.Get("/users", cfg.handlerGetUsers)
	v1Router.Get("/users/{userID}", cfg.handlerGetUserByID)
	v1Router.Post("/login", cfg.handlerUsersLogin)

	authRouter := chi.NewRouter()

	authRouter.Use(cfg.middlewareAuth)

	authRouter.Put("/users", cfg.handlerUsersUpdate)
	authRouter.Delete("/users", cfg.handlerDeleteUser)

	authRouter.Post("/accounts", cfg.handlerCreateAccount)
	authRouter.Get("/accounts", cfg.handlerGetAccounts)
	authRouter.Put("/accounts", cfg.handlerUpdateAccount)
	authRouter.Delete("/accounts/{accountID}", cfg.handlerDeleteAccount)

	authRouter.Get("/balances", cfg.handlerGetAccountsBalances)

	authRouter.Post("/transactions", cfg.handlerCreateTransaction)
	authRouter.Get("/transactions", cfg.handlerGetTransactions)
	authRouter.Delete("/transactions/{transactionID}", cfg.handlerDeleteTransaction)

	authRouter.Get("/others", cfg.handlerGetOthersAccounts)

	router.Mount("/v1", v1Router)
	v1Router.Mount("/auth", authRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	log.Fatal(srv.ListenAndServe())
}
