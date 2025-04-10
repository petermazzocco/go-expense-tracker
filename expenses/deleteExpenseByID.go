package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteExpenseByID(c *gin.Context) {
	// Get id value from the param path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Get the user ID
	userID := auth.GetUserIDFromCookie(c)

	// Find the post where the ID matches the param and userID matches the cookie userID
	var expense models.Expense
	result := initializers.DB.Where("id = ?", id).Where("user_id = ?", userID).Find(&expense)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find expense by ID or userID is not valid for this expense"})
		return
	}

	// Set deleted at value with the soft delete from gorm
	if err := initializers.DB.Delete(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.Header("HX-Redirect", "/expenses")
	// Return expense details
	c.JSON(http.StatusOK, gin.H{
		"expense": expense,
		"message": "Expense deleted successfully",
	})
}
