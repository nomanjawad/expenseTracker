package models

import "time"

type Expense struct {
	ID          string  	`json:"id"`
	Amount      float64 	`json:"amount" binding:"required"`
	Purpose     string  	`json:"purpose" binding:"required"`
	Account		string  	`json:"account_name" binding:"required"`
	CreatedAt   string  	`json:"created_at"`
	Date    	time.Time  	`json:"date"`
}
