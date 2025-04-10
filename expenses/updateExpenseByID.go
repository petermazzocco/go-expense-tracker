package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateExpenseByID(c *gin.Context) {
	// Get id value from the param path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	title := c.PostForm("title")
	category := c.PostForm("category")
	amount := c.PostForm("amount")
	if title == "" || category == "" || amount == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Category, and Amount are required"})
		return
	}

	// Get the user ID
	userID := auth.GetUserIDFromCookie(c)

	// First find the existing expense
	var existingExpense models.Expense
	result := initializers.DB.Where("id = ?", id).Where("user_id = ?", userID).First(&existingExpense)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find expense by ID or userID is not valid for this expense"})
		return
	}

	// Now bind the updated data

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	existingExpense.Amount = float64(amountInt)
	existingExpense.Category = category
	existingExpense.Title = title

	// Save the updated expense
	result = initializers.DB.Save(&existingExpense)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}
	c.Header("HX-Redirect", "/expenses/"+id)
	// Return expense details
	c.JSON(http.StatusOK, gin.H{
		"expense": existingExpense,
		"message": "Expense updated successfully",
	})
}
