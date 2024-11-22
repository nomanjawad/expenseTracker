package controllers

import (
	"context"
	"expenceTracker/backend/config"
	"expenceTracker/backend/models"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetExpenses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Define the SQL query
	query := `
		SELECT id, amount, purpose, account, date 
		FROM public.expenses 
		WHERE user_id = $1 
		ORDER BY date DESC
	`

	// Execute the query
	rows, err := config.DB.Query(context.Background(), query, userID)
	if err != nil {
		log.Printf("Query error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	defer rows.Close()

	// Parse the result into a slice of expenses
	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Purpose, &expense.Account, &expense.Date)
		if err != nil {
			log.Printf("Scan error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process expense data"})
			return
		}
		expenses = append(expenses, expense)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}


// AddExpense adds a new expense for the authenticated user
func AddExpense(c *gin.Context) {
	// Extract user_id from the context (set by the authentication middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse the expense details from the request body
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate required fields
	if expense.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}
	if expense.Account == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account is required"})
		return
	}

	// Generate a new UUID for the expense
	expenseID := uuid.New()

	// Insert the expense into the database
	query := `INSERT INTO public.expenses (id, user_id, amount, purpose, account, date) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := config.DB.Exec(
		context.Background(),
		query,
		expenseID,        // id
		userID,           // user_id
		expense.Amount,   // amount
		expense.Purpose,  // purpose
		expense.Account,  // account
		expense.Date,     // date
	)
	if err != nil {
		log.Printf("Adding error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add expense"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":    "Expense added successfully",
		"expense_id": expenseID,
	})
}
