package auth

import (
	"encoding/base64"
	"fmt"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserIDFromCookie(c *gin.Context) int {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - no userID in context"})
		return 0
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - userID not a string"})
		return 0
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error converting UserID string '%s' to int: %v", userID, err)})
		return 0
	}

	return userIDInt
}

func SignIn(c *gin.Context) {
	// Get request data
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate required fields
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, name, and password are required"})
		return
	}

	// Find user in DB
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare stored hashed password with the given password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Decode user's salt
	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding salt"})
		return
	}

	// Derive the encryption key
	encryptionKey := DeriveKey(password, salt)
	encodedKey := EncodeToBase64(encryptionKey)

	// Format user ID correctly for JWT
	userIDString := strconv.FormatUint(uint64(user.ID), 10)

	// Generate a session token (JWT)
	token, err := GenerateJWT(userIDString, encodedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Set the JWT as an HTTPOnly cookie
	c.SetCookie(
		"token",
		token,
		3600*24*30, // one month
		"/",
		"",
		false, // HTTPS only
		true,  // not accessible via JavaScript
	)
	// Add redirect header for HTMX
	c.Header("HX-Redirect", "/expenses")
	// Return the user object and a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})
}

func SignUp(c *gin.Context) {
	email := c.PostForm("email")
	name := c.PostForm("name")
	password := c.PostForm("password")

	// Validate required fields
	if email == "" || name == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, name, and password are required"})
		return
	}

	user := models.User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	// Generate a salt for the user
	salt, err := GenerateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating salt"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}
	user.Password = string(hashedPassword)

	// Store salt in database
	user.Salt = EncodeToBase64(salt)

	// Save user to DB
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Add redirect header for HTMX
	c.Header("HX-Redirect", "/expenses")
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func SignOut(c *gin.Context) {
	// Clear the token cookie
	c.SetCookie(
		"token",
		"",
		-1, // expire immediately
		"/",
		"",
		false, // HTTPS only
		true,  // not accessible via JavaScript
	)

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}
