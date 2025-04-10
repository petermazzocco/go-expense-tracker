package helpers

import (
	"errors"
	"go-expense-tracker/initializers"
	"go-expense-tracker/models"

	"github.com/gin-gonic/gin"
)

func GetAllExpensesHelper(userID int, c *gin.Context) ([]models.Expense, error) {
	expenses := []models.Expense{}
	if err := initializers.DB.Where("user_id = ?", userID).Scopes(Paginate(c)).Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func GetExpenseByIDHelper(id string, userID int) (*models.Expense, error) {
	expense := models.Expense{}
	result := initializers.DB.Where("id = ?", id).Where("user_id = ?", userID).Find(&expense)
	if result.Error != nil {
		return nil, errors.New("Database error")
	}
	return &expense, nil
}
