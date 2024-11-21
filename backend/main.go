package main

import (
	"expenceTracker/backend/config"
	"expenceTracker/backend/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	config.ConnectDB()
	defer config.CloseDB()

	// Initialize Gin router
	r := gin.Default()

	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Expense Tracker API",
			"routes": []map[string]string{
				{"route": "/register", "method": "POST", "description": "Register a new user"},
				{"route": "/login", "method": "POST", "description": "Login a user and get a JWT"},
				{"route": "/expenses", "method": "GET", "description": "Get all expenses for the authenticated user"},
				{"route": "/expenses", "method": "POST", "description": "Add a new expense for the authenticated user"},
			},
		})
	})

	// Register routes
	routes.AuthRoutes(r)
	routes.ExpenseRoutes(r)

		// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Start the server
	r.Run(":8080")
}
