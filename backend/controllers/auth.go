package controllers

import (
	"context"
	"expenceTracker/backend/config"
	"expenceTracker/backend/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey string

// init function to read the JWT secret key from the .env file
func init()  {
	err := godotenv.Load(".env")
	if err!=nil{
		panic("Error loading .env file")
	}
	jwtKey = os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		panic("JWT_SECRET is not set in .env")
	}
}

func RegisterUser(c *gin.Context) {
	var user models.User

	// Bind JSON to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing the password"})
		return
	}
	
	// Generate a unique ID for the user
	userID := uuid.New()

		// Insert the user into the database
	query := `
		INSERT INTO auth.users (
			id, email, encrypted_password, role, created_at, updated_at, is_anonymous, is_sso_user
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = config.DB.Exec(
		context.Background(),
		query,
		userID,                 // id
		user.Email,             // email
		string(hashedPassword), // encrypted_password
		"user",                 // role
		time.Now(),             // created_at
		time.Now(),             // updated_at
		false,                  // is_anonymous
		false,                  // is_sso_user
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving the user to the database", "details": err.Error()})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginUser function to authenticate a user
func LoginUser(c *gin.Context) {
	var credentials models.User

	// Parse JSON input
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Fetch the user from the database
	var user models.User
	query := `SELECT id, email, encrypted_password FROM auth.users WHERE email = $1`
	err := config.DB.QueryRow(context.Background(), query, credentials.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Printf("invalid user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		log.Printf("invalid password: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token for the authenticated user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with the JWT token and success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}