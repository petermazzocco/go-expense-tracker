package auth

import (
	"encoding/base64"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUserIDFromCookie(c *gin.Context) int {
	userID, exists := c.MustGet("userID").(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return 0
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error converting UserID string to int"})
		return 0
	}
	return userIDInt
}

func SignIn(c *gin.Context) {
	// Get request data
	var request struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required,password"`
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Find user in DB
	var user models.User
	if err := initializers.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare stored hashed password with the given password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
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
	encryptionKey := DeriveKey(request.Password, salt)
	encodedKey := EncodeToBase64(encryptionKey)
	// Generate a session token (JWT)
	token, err := GenerateJWT(strconv.FormatUint(uint64(user.ID), 10), encodedKey)
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
	c.Header("HX-Redirect", "/dashboard")
	// Return the user object and a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
	})
}

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
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
