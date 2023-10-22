package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MeirionL/boing-block/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // Imported database driver we need
)

type apiConfig struct {
	DB        *database.Queries
	jwtSecret string
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	cfg := apiConfig{
		DB:        database.New(db),
		jwtSecret: jwtSecret,
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

	router.Get("/healthz", handlerReadiness)
	router.Get("/err", handlerErr)

	router.Post("/users", cfg.handlerCreateUser)
	router.Get("/users", cfg.handlerGetUsers)
	router.Get("/users/{userID}", cfg.handlerGetUserByID)
	router.Post("/login", cfg.handlerUsersLogin)

	// Router that implements user authentication
	authRouter := chi.NewRouter()
	authRouter.Use(cfg.middlewareAuth)

	authRouter.Put("/users", cfg.handlerUsersUpdate)
	authRouter.Delete("/users", cfg.handlerDeleteUser)

	authRouter.Post("/revoke", cfg.handlerRevoke)
	authRouter.Get("/revoke", cfg.handlerGetRevokedTokens)
	authRouter.Post("/refresh", cfg.handlerRefresh)

	authRouter.Post("/accounts", cfg.handlerCreateAccount)
	authRouter.Get("/accounts", cfg.handlerGetAccounts)
	authRouter.Put("/accounts", cfg.handlerUpdateAccount)
	authRouter.Delete("/accounts/{accountID}", cfg.handlerDeleteAccount)

	authRouter.Get("/balances", cfg.handlerGetAccountsBalances)

	authRouter.Post("/transactions", cfg.handlerCreateTransaction)
	authRouter.Get("/transactions", cfg.handlerGetTransactions)
	authRouter.Delete("/transactions/{transactionID}", cfg.handlerDeleteTransaction)

	authRouter.Get("/others", cfg.handlerGetOthersAccounts)

	router.Mount("/auth", authRouter)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	log.Fatal(srv.ListenAndServe())
}
