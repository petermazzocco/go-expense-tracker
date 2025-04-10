package expenses

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/helpers"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewAllExpenses(c *gin.Context) {
	var expenses []models.Expense
	var totalCount int64
	userID := auth.GetUserIDFromCookie(c)

	// Get total count for pagination info
	initializers.DB.Model(&models.Expense{}).Where("user_id = ?", userID).Count(&totalCount)

	expenses, err := helpers.GetAllExpensesHelper(userID, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	// Return expenses with pagination info
	c.JSON(http.StatusOK, gin.H{
		"expenses": expenses,
		"message":  "All expenses fetched successfully",
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  totalCount,
			"total_pages":  int(math.Ceil(float64(totalCount) / float64(pageSize))),
		},
	})
}
