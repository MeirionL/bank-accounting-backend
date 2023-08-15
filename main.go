package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type transaction struct {
	Payment   float64 `json:"payment"`
	PreTotal  float64 `json:"pretotal"`
	Recipient string  `json:"recipient"`
}

func getTransactions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, transactions)
}

var transactions = []transaction{
	{Payment: 56.99, PreTotal: 200, Recipient: "John"},
	{Payment: 17.99, PreTotal: 300, Recipient: "Connor"},
	{Payment: 39.99, PreTotal: 132.47, Recipient: "Doyle"},
}

func postTransactions(c *gin.Context) {
	var newTransaction transaction
	if err := c.BindJSON(&newTransaction); err != nil {
		return
	}
	transactions = append(transactions, newTransaction)
	c.IndentedJSON(http.StatusCreated, newTransaction)
}

func getTransactionByRecipient(c *gin.Context) {
	recipient := c.Param("recipient")

	for _, t := range transactions {
		if t.Recipient == recipient {
			c.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "recipient not found"})
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the environment")
	}

	router := gin.Default() //initialising router
	router.GET("/transactions", getTransactions)
	router.GET("/transactions/:recipient", getTransactionByRecipient)
	router.POST("/transactions", postTransactions)
	addr := "localhost:" + portString
	router.Run(addr) //attaching it to an http.Server
}
