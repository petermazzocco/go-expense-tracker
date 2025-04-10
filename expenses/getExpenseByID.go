package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetExpenseByID(c *gin.Context) {
	// Get id value from the param path
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Get the user ID
	userID := auth.GetUserIDFromCookie(c)

	// Find the post where the ID matches the param and userID matches the cookie userID
	expense, err := helpers.GetExpenseByIDHelper(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get expense"})
		return
	}
	// Return expense details
	c.JSON(http.StatusOK, gin.H{
		"expense": expense,
		"message": "Expense deleted successfully",
	})
}
