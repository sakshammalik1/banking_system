

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

var db *pg.DB

func initDB() {
	db = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "ramesh@rani1",
		Database: "banking",
	})
	
	// Check connection
	_, err := db.Exec("SELECT 1")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	fmt.Println("Successfully connected to database!")
}

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type Account struct {
	ID            int     `json:"id"`
	CustomerID    int     `json:"customer_id"`
	BranchID      int     `json:"branch_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
	Type          string  `json:"type"`
}

type Transaction struct {
	ID        int     `json:"id"`
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`
	Date      string  `json:"date"`
}

type Loan struct {
	ID           int     `json:"id"`
	CustomerID   int     `json:"customer_id"`
	BranchID     int     `json:"branch_id"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Balance      float64 `json:"balance"`
}

// API Handlers

// Create Customer
func createCustomer(c *gin.Context) {
	var customer Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Model(&customer).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// View Customer
func viewCustomer(c *gin.Context) {
	customerId := c.Param("customerId")
	var customer Customer

	err := db.Model(&customer).Where("id = ?", customerId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// Open Account
func openAccount(c *gin.Context) {
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.AccountNumber = fmt.Sprintf("ACC-%d", account.CustomerID) // Simplified account number generation
	_, err := db.Model(&account).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// View Account
func viewAccount(c *gin.Context) {
	accountId := c.Param("accountId")
	var account Account

	err := db.Model(&account).Where("id = ?", accountId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// Deposit Money
func depositMoney(c *gin.Context) {
	accountId := c.Param("accountId")
	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account Account
	err := db.Model(&account).Where("id = ?", accountId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	account.Balance += transaction.Amount
	_, err = db.Model(&account).Where("id = ?", accountId).Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction.AccountID = account.ID
	transaction.Type = "Deposit"
	_, err = db.Model(&transaction).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// Withdraw Money
func withdrawMoney(c *gin.Context) {
	accountId := c.Param("accountId")
	var transaction Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var account Account
	err := db.Model(&account).Where("id = ?", accountId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.Balance < transaction.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	account.Balance -= transaction.Amount
	_, err = db.Model(&account).Where("id = ?", accountId).Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction.AccountID = account.ID
	transaction.Type = "Withdraw"
	_, err = db.Model(&transaction).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// View Transactions
func viewTransactions(c *gin.Context) {
	accountId := c.Param("accountId")
	var transactions []Transaction

	err := db.Model(&transactions).Where("account_id = ?", accountId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// Take Loan
func takeLoan(c *gin.Context) {
	var loan Loan
	if err := c.ShouldBindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan.InterestRate = 12.0 // Fixed interest rate of 12%
	loan.Balance = loan.Amount

	_, err := db.Model(&loan).Insert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// View Loan
func viewLoan(c *gin.Context) {
	loanId := c.Param("loanId")
	var loan Loan

	err := db.Model(&loan).Where("id = ?", loanId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loan)
}

// Repay Loan
func repayLoan(c *gin.Context) {
	loanId := c.Param("loanId")
	var repayment struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&repayment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loan Loan
	err := db.Model(&loan).Where("id = ?", loanId).Select()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if repayment.Amount > loan.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repayment amount exceeds loan balance"})
		return
	}

	loan.Balance -= repayment.Amount
	_, err = db.Model(&loan).Where("id = ?", loanId).Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Repayment successful, new balance: %f", loan.Balance)})
}

func main() {
	initDB()
	defer db.Close()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow the frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
}
	router := gin.Default()
router.Use(cors.New(config))
	router.POST("/customers", createCustomer)
	router.GET("/customers/:customerId", viewCustomer)
	router.POST("/accounts", openAccount)
	router.GET("/accounts/:accountId", viewAccount)
	router.POST("/accounts/:accountId/deposit", depositMoney)
	router.POST("/accounts/:accountId/withdraw", withdrawMoney)
	router.GET("/accounts/:accountId/transactions", viewTransactions)
	router.POST("/loans", takeLoan)
	router.GET("/loans/:loanId", viewLoan)
	router.POST("/loans/:loanId/repay", repayLoan)

	router.Run(":8080")
} 
