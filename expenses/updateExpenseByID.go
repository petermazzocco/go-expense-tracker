package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateExpenseByID(c *gin.Context) {
	// Get id value from the param path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
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
	var updatedData models.Expense
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	existingExpense.Amount = updatedData.Amount
	existingExpense.Category = updatedData.Category
	existingExpense.Title = updatedData.Title

	// Save the updated expense
	result = initializers.DB.Save(&existingExpense)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	// Return expense details
	c.JSON(http.StatusOK, gin.H{
		"expense": existingExpense,
		"message": "Expense updated successfully",
	})
}
