package controllers

import (
	"context"
	"expenceTracker/backend/config"
	"expenceTracker/backend/models"
	"log"
	"time"

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

	// Query expenses from the public schema
	query := `SELECT id, amount, purpose, account, created_at FROM public.expenses WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := config.DB.Query(context.Background(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	defer rows.Close()

	// Parse the results into a slice
	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.ID, &expense.Amount, &expense.Purpose, &expense.Account, &expense.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process expense data"})
			return
		}
		expenses = append(expenses, expense)
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

		// Handle the `date` field: Default to the current time if it's empty
	if expense.Date == "" {
		expense.Date = time.Now().Format(time.RFC3339)
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
