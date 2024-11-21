package routes

import (
	"expenceTracker/backend/controllers"
	"expenceTracker/backend/middleware"

	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(router *gin.Engine) {
	expenseGroup := router.Group("/expenses")
	expenseGroup.Use(middleware.Authenticate())
	{
		expenseGroup.GET("/", controllers.GetExpenses)
		expenseGroup.POST("/", controllers.AddExpense)
	}
}
