package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNewExpense(c *gin.Context) {
	// Validate request body
	var expense models.Expense
	if err := c.BindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	var user models.User
	userID := auth.GetUserIDFromCookie(c)

	// Check if user exists in the db
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	expense.UserID = uint(userID)
	// Create the expense for the specific user
	if err := initializers.DB.Create(&expense).Where("userID = ?", userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create expense"})
		return
	}

	// Return expense details and success message
	c.JSON(http.StatusOK, gin.H{
		"expense": expense,
		"message": "Expense created successfully",
	})
}
