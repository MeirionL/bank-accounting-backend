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

type transaction struct {
	ID             int     `json:"id"`
	PreBalance     float64 `json:"pre_balance"`
	PaymentAmount  float64 `json:"payment"`
	PostBalance    float64 `json:"post_balance"`
	PartnerAccount string  `json:"partner_account"`
	AccountNumber  string  `json:"account_number"`
	SortCode       string  `json:"sort_code"`
	NewPartner     bool    `json:"paid_before"`
	TimeOf         struct {
		Time time.Time `json:"time"`
		Date time.Time `json:"data"`
	}
	Message string `json:"message"`
}

type apiConfig struct {
	DB *database.Queries // Calling the database package in my directory
}

// var transactions = []transaction{
// 	{ID: "1", Payment: 56.99, PreTotal: 1000, PostTotal: 943.01, Recipient: "John", AccountNumber: "95028403", SortCode: "23-87-04", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02")},
// 	{ID: "2", Payment: 17.99, PreTotal: 943.01, PostTotal: 925.02, Recipient: "Connor", AccountNumber: "59058024", SortCode: "67-98-23", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02"), Message: "Yes mate"},
// 	{ID: "3", Payment: 39.99, PreTotal: 925.02, PostTotal: 885.03, Recipient: "Doyle", AccountNumber: "58301109", SortCode: "56-17-56", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02"), Message: "Payment for XYZ"},
// }

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL) // Loading in our database
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	apiCfg := apiConfig{ // An API config we can pass to our
		DB: database.New(db), // handler so that it has access to
	} // our database.

	router := chi.NewRouter() //initialising router

	router.Use(cors.Handler(cors.Options{ //Tells our server to send extra http headers in our responses that will tell browsers we allow this stuff.
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
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/transactions/last", apiCfg.handlerLastTransactionGet)
	v1Router.Get("/transactions", apiCfg.handlerTransactionsRetrieve)
	v1Router.Get("/transactions/{id}", apiCfg.handlerTransactionGet)
	v1Router.Post("/transactions", apiCfg.handlerTransactionsCreate)
	router.Get("/balance", apiCfg.handlerGetBalance)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server starting on port %v", portString)
}

// func getLastTransaction(c *gin.Context) {
// 	transaction, err := getTransactionById(fmt.Sprint(CurrId))
// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
// 	}

// 	c.IndentedJSON(http.StatusOK, transaction)
// }

// func getTransactions(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, transactions)
// }

// func postTransactions(c *gin.Context) {
// 	var newTransaction transaction
// 	if err := c.BindJSON(&newTransaction); err != nil {
// 		return
// 	}

// 	CurrId++
// 	newTransaction.ID = strconv.Itoa(CurrId)

// 	newTransaction.PreTotal = Total
// 	Total -= newTransaction.Payment
// 	newTransaction.PostTotal = Total

// 	if hasPrevTransaction(newTransaction.AccountNumber, newTransaction.SortCode) {
// 		newTransaction.PaidBefore = true
// 	} else {
// 		newTransaction.PaidBefore = false
// 	}

// 	newTransaction.Time = fmt.Sprint(time.Now().Format("15:04"))
// 	newTransaction.Date = fmt.Sprint(time.Now().Format("2006-01-02"))

// 	transactions = append(transactions, newTransaction)
// 	c.IndentedJSON(http.StatusCreated, newTransaction)
// }

// func hasPrevTransaction(an, sc string) bool {
// 	for _, t := range transactions {
// 		if t.AccountNumber == an && t.SortCode == sc {
// 			return true
// 		}
// 	}
// 	return false
// }

// func transactionById(c *gin.Context) {
// 	id := c.Param("id")
// 	transaction, err := getTransactionById(id)

// 	if err != nil {
// 		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
// 		return
// 	}

// 	c.IndentedJSON(http.StatusOK, transaction)
// }

// func getTransactionById(id string) (*transaction, error) {
// 	for i, t := range transactions {
// 		if t.ID == id {
// 			return &transactions[i], nil
// 		}
// 	}

// 	return nil, errors.New("transaction not found")
// }
