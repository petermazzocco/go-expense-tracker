package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateNewExpense(c *gin.Context) {
	// Validate request body
	var expense models.Expense
	title := c.PostForm("title")
	category := c.PostForm("category")
	amount := c.PostForm("amount")

	// Get user ID from JWT token
	var user models.User
	userID := auth.GetUserIDFromCookie(c)

	// Check if user exists in the db
	if err := initializers.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	expense.UserID = uint(userID)
	expense.Title = title
	expense.Category = category
	expenseAmount, err := strconv.Atoi(amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	expense.Amount = float64(expenseAmount)

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
