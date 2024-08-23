// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	//"os"

// 	"github.com/gorilla/mux"
// 	_ "github.com/lib/pq"
// )

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "ramesh@rani1"
// 	dbname   = "banking"
// )

// var db *sql.DB

// func initDB() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)

// 	var err error
// 	db, err = sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Successfully connected to database!")
// }

// type Customer struct {
// 	ID      int    `json:"id"`
// 	Name    string `json:"name"`
// 	Address string `json:"address"`
// 	Phone   string `json:"phone"`
// 	Email   string `json:"email"`
// }

// type Account struct {
// 	ID           int     `json:"id"`
// 	CustomerID   int     `json:"customer_id"`
// 	BranchID     int     `json:"branch_id"`
// 	AccountNumber string `json:"account_number"`
// 	Balance      float64 `json:"balance"`
// 	Type         string  `json:"type"`
// }

// type Transaction struct {
// 	ID        int     `json:"id"`
// 	AccountID int     `json:"account_id"`
// 	Amount    float64 `json:"amount"`
// 	Type      string  `json:"type"`
// 	Date      string  `json:"date"`
// }

// type Loan struct {
// 	ID           int     `json:"id"`
// 	CustomerID   int     `json:"customer_id"`
// 	BranchID     int     `json:"branch_id"`
// 	Amount       float64 `json:"amount"`
// 	InterestRate float64 `json:"interest_rate"`
// 	StartDate    string  `json:"start_date"`
// 	EndDate      string  `json:"end_date"`
// 	Balance      float64 `json:"balance"`
// }

// // API Handlers

// // Create Customer
// func createCustomer(w http.ResponseWriter, r *http.Request) {
// 	var customer Customer
// 	err := json.NewDecoder(r.Body).Decode(&customer)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	sqlStatement := `INSERT INTO customers (name, address, phone, email) VALUES ($1, $2, $3, $4) RETURNING id`
// 	err = db.QueryRow(sqlStatement, customer.Name, customer.Address, customer.Phone, customer.Email).Scan(&customer.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(customer)
// }

// // View Customer
// func viewCustomer(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	customerId := params["customerId"]

// 	var customer Customer
// 	sqlStatement := `SELECT id, name, address, phone, email FROM customers WHERE id=$1`
// 	row := db.QueryRow(sqlStatement, customerId)
// 	err := row.Scan(&customer.ID, &customer.Name, &customer.Address, &customer.Phone, &customer.Email)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(customer)
// }

// // Open Account
// func openAccount(w http.ResponseWriter, r *http.Request) {
// 	var account Account
// 	err := json.NewDecoder(r.Body).Decode(&account)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	account.AccountNumber = fmt.Sprintf("ACC-%d", account.CustomerID) // Simplified account number generation
// 	sqlStatement := `INSERT INTO accounts (customer_id, branch_id, account_number, balance, type) VALUES ($1, $2, $3, $4, $5) RETURNING id`
// 	err = db.QueryRow(sqlStatement, account.CustomerID, account.BranchID, account.AccountNumber, account.Balance, account.Type).Scan(&account.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(account)
// }

// // View Account
// func viewAccount(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	accountId := params["accountId"]

// 	var account Account
// 	sqlStatement := `SELECT id, customer_id, branch_id, account_number, balance, type FROM accounts WHERE id=$1`
// 	row := db.QueryRow(sqlStatement, accountId)
// 	err := row.Scan(&account.ID, &account.CustomerID, &account.BranchID, &account.AccountNumber, &account.Balance, &account.Type)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(account)
// }

// // Deposit Money
// func depositMoney(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	accountId := params["accountId"]

// 	var transaction Transaction
// 	err := json.NewDecoder(r.Body).Decode(&transaction)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var balance float64
// 	err = db.QueryRow(`SELECT balance FROM accounts WHERE id=$1`, accountId).Scan(&balance)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	newBalance := balance + transaction.Amount
// 	_, err = db.Exec(`UPDATE accounts SET balance=$1 WHERE id=$2`, newBalance, accountId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	sqlStatement := `INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3) RETURNING id`
// 	err = db.QueryRow(sqlStatement, accountId, transaction.Amount, "Deposit").Scan(&transaction.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(transaction)
// }

// // Withdraw Money
// func withdrawMoney(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	accountId := params["accountId"]

// 	var transaction Transaction
// 	err := json.NewDecoder(r.Body).Decode(&transaction)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var balance float64
// 	err = db.QueryRow(`SELECT balance FROM accounts WHERE id=$1`, accountId).Scan(&balance)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if balance < transaction.Amount {
// 		http.Error(w, "Insufficient balance", http.StatusBadRequest)
// 		return
// 	}

// 	newBalance := balance - transaction.Amount
// 	_, err = db.Exec(`UPDATE accounts SET balance=$1 WHERE id=$2`, newBalance, accountId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	sqlStatement := `INSERT INTO transactions (account_id, amount, type) VALUES ($1, $2, $3) RETURNING id`
// 	err = db.QueryRow(sqlStatement, accountId, transaction.Amount, "Withdraw").Scan(&transaction.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(transaction)
// }

// // View Transactions
// func viewTransactions(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	accountId := params["accountId"]

// 	rows, err := db.Query(`SELECT id, account_id, amount, type, date FROM transactions WHERE account_id=$1`, accountId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var transactions []Transaction
// 	for rows.Next() {
// 		var transaction Transaction
// 		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Type, &transaction.Date)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		transactions = append(transactions, transaction)
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(transactions)
// }

// // Take Loan
// func takeLoan(w http.ResponseWriter, r *http.Request) {
// 	var loan Loan
// 	err := json.NewDecoder(r.Body).Decode(&loan)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	loan.InterestRate = 12.0 // Fixed interest rate of 12%
// 	loan.Balance = loan.Amount

// 	sqlStatement := `INSERT INTO loans (customer_id, branch_id, amount, interest_rate, start_date, end_date, balance) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
// 	err = db.QueryRow(sqlStatement, loan.CustomerID, loan.BranchID, loan.Amount, loan.InterestRate, loan.StartDate, loan.EndDate, loan.Balance).Scan(&loan.ID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(loan)
// }

// // View Loan
// func viewLoan(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	loanId := params["loanId"]

// 	var loan Loan
// 	sqlStatement := `SELECT id, customer_id, branch_id, amount, interest_rate, start_date, end_date, balance FROM loans WHERE id=$1`
// 	row := db.QueryRow(sqlStatement, loanId)
// 	err := row.Scan(&loan.ID, &loan.CustomerID, &loan.BranchID, &loan.Amount, &loan.InterestRate, &loan.StartDate, &loan.EndDate, &loan.Balance)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(loan)
// }

// // Repay Loan
// func repayLoan(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	loanId := params["loanId"]

// 	var repayment struct {
// 		Amount float64 `json:"amount"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&repayment)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var balance float64
// 	err = db.QueryRow(`SELECT balance FROM loans WHERE id=$1`, loanId).Scan(&balance)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if repayment.Amount > balance {
// 		http.Error(w, "Repayment amount exceeds loan balance", http.StatusBadRequest)
// 		return
// 	}

// 	newBalance := balance - repayment.Amount
// 	_, err = db.Exec(`UPDATE loans SET balance=$1 WHERE id=$2`, newBalance, loanId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Repayment successful, new balance: %f", newBalance)
// }

// func main() {
// 	initDB()

// 	router := mux.NewRouter()
// 	router.HandleFunc("/customers", createCustomer).Methods("POST")
// 	router.HandleFunc("/customers/{customerId}", viewCustomer).Methods("GET")
// 	router.HandleFunc("/accounts", openAccount).Methods("POST")
// 	router.HandleFunc("/accounts/{accountId}", viewAccount).Methods("GET")
// 	router.HandleFunc("/accounts/{accountId}/deposit", depositMoney).Methods("POST")
// 	router.HandleFunc("/accounts/{accountId}/withdraw", withdrawMoney).Methods("POST")
// 	router.HandleFunc("/accounts/{accountId}/transactions", viewTransactions).Methods("GET")
// 	router.HandleFunc("/loans", takeLoan).Methods("POST")
// 	router.HandleFunc("/loans/{loanId}", viewLoan).Methods("GET")
// 	router.HandleFunc("/loans/{loanId}/repay", repayLoan).Methods("POST")

// 	log.Fatal(http.ListenAndServe(":8080", router))
// }

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