package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type transaction struct {
	ID        string  `json:"id"`
	Payment   float64 `json:"payment"`
	PreTotal  float64 `json:"pretotal"`
	Recipient string  `json:"recipient"`
}

var transactions = []transaction{
	{ID: "1", Payment: 56.99, PreTotal: 200, Recipient: "John"},
	{ID: "2", Payment: 17.99, PreTotal: 300, Recipient: "Connor"},
	{ID: "3", Payment: 39.99, PreTotal: 132.47, Recipient: "Doyle"},
}

func getTransactions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, transactions)
}

func postTransactions(c *gin.Context) {
	// NEEDS LOGIC FOR FINDING APPROPRIATE PARAMETERS AND SUBTRACTING THEM
	id, ok := c.GetQuery("id")
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing transaction id"})
		return
	}

	transaction, err := getTransactionById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	var newTransaction transaction
	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}
	transactions = append(transactions, newTransaction)
	c.IndentedJSON(http.StatusCreated, newTransaction)
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
	for i, b := range transactions {
		if b.ID == id {
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

	router := gin.Default()                      //initialising router
	router.GET("/transactions", getTransactions) // sending a GET request to /transactions calls the getTransactions function
	router.GET("/transactions/:id", transactionById)
	router.POST("/transactions", postTransactions)
	router.Run(portString) //attaching it to an http.Server
}
