package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var Total float64 = 885.03
var CurrId int = 3

type transaction struct {
	ID            string  `json:"id"`
	Payment       float64 `json:"payment"`
	PreTotal      float64 `json:"preTotal"`
	PostTotal     float64 `json:"postTotal"`
	Recipient     string  `json:"recipient"`
	AccountNumber string  `json:"accountNumber"`
	SortCode      string  `json:"sortCode"`
	PaidBefore    bool    `json:"paidBefore"`
	Time          string  `json:"time"`
	Date          string  `json:"data"`
	Message       string  `json:"message"`
}

var transactions = []transaction{
	{ID: "1", Payment: 56.99, PreTotal: 1000, PostTotal: 943.01, Recipient: "John", AccountNumber: "95028403", SortCode: "23-87-04", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02")},
	{ID: "2", Payment: 17.99, PreTotal: 943.01, PostTotal: 925.02, Recipient: "Connor", AccountNumber: "59058024", SortCode: "67-98-23", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02"), Message: "Yes mate"},
	{ID: "3", Payment: 39.99, PreTotal: 925.02, PostTotal: 885.03, Recipient: "Doyle", AccountNumber: "58301109", SortCode: "56-17-56", Time: time.Now().Format("15:04"), Date: time.Now().Format("2006-01-02"), Message: "Payment for XYZ"},
}

func getBalance(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"balance": Total})
}

func getLastTransaction(c *gin.Context) {
	transaction, err := getTransactionById(fmt.Sprint(CurrId))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
	}

	c.IndentedJSON(http.StatusOK, transaction)
}

func getTransactions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, transactions)
}

func postTransactions(c *gin.Context) {
	var newTransaction transaction
	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}

	CurrId++
	newTransaction.ID = strconv.Itoa(CurrId)

	newTransaction.PreTotal = Total
	Total -= newTransaction.Payment
	newTransaction.PostTotal = Total

	if hasPrevTransaction(newTransaction.AccountNumber, newTransaction.SortCode) {
		newTransaction.PaidBefore = true
	} else {
		newTransaction.PaidBefore = false
	}

	newTransaction.Time = fmt.Sprint(time.Now().Format("15:04"))
	newTransaction.Date = fmt.Sprint(time.Now().Format("2006-01-02"))

	transactions = append(transactions, newTransaction)
	c.IndentedJSON(http.StatusCreated, newTransaction)
}

func hasPrevTransaction(an, sc string) bool {
	for _, t := range transactions {
		if t.AccountNumber == an && t.SortCode == sc {
			return true
		}
	}
	return false
}

func transactionById(c *gin.Context) {
	id := c.Param("id")
	transaction, err := getTransactionById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, transaction)
}

func getTransactionById(id string) (*transaction, error) {
	for i, t := range transactions {
		if t.ID == id {
			return &transactions[i], nil
		}
	}

	return nil, errors.New("transaction not found")
}

func main() {
	godotenv.Load(".env")

	portString := "localhost:" + os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the environment")
	}

	router := gin.Default() //initialising router

	log.Printf("Server starting on port %v", portString)

	v1Router := router.Group("/v1")

	v1Router.GET("/healthz", func(c *gin.Context) { // A path to call to check if your server is still live
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1Router.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An eror occurred: %s", err)})
		}
	})

	router.GET("/balance", getBalance)
	router.GET("/lastTransaction", getLastTransaction)
	router.GET("/transactions", getTransactions) // sending a GET request to /transactions calls the getTransactions function
	router.GET("/transactions/:id", transactionById)
	router.POST("/transactions", postTransactions)
	router.Run(portString) //attaching it to an http.Server
}
