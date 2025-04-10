package main

import (
	"go-expense-tracker/auth"
	"go-expense-tracker/expenses"
	"go-expense-tracker/helpers"
	"go-expense-tracker/initializers"
	"go-expense-tracker/renderer"
	"go-expense-tracker/templates/pages"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitDatabase()
	initializers.InitENV()
}

func main() {
	// Router
	r := gin.Default()
	r.LoadHTMLFiles("./templates/*.html")

	// Rendere for HTML
	ginHtmlRenderer := r.HTMLRender
	r.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable warning
	r.SetTrustedProxies(nil)

	// Home route
	r.GET("/", func(c *gin.Context) {
		page := renderer.New(c.Request.Context(), http.StatusOK, pages.Index("Expense Tracker"))
		c.Render(http.StatusOK, page)
	})

	// User auth routes
	r.POST("/signin", auth.SignIn)
	r.POST("/signup", auth.SignUp)
	r.POST("/signout", auth.SignOut)

	// User auth pages
	r.GET("/signin", func(c *gin.Context) {
		page := renderer.New(c.Request.Context(), http.StatusOK, pages.SignIn())
		c.Render(http.StatusOK, page)
	})
	r.GET("/signup", func(c *gin.Context) {
		page := renderer.New(c.Request.Context(), http.StatusOK, pages.SignUp())
		c.Render(http.StatusOK, page)
	})

	// Protected API routes group
	api := r.Group("/api")
	api.Use(auth.Middleware())
	{
		api.GET("/expenses", expenses.ViewAllExpenses)
		api.GET("/expenses/:id", expenses.GetExpenseByID)
		api.POST("/expenses", expenses.CreateNewExpense)
		api.PUT("/expenses/:id", expenses.UpdateExpenseByID)
		api.DELETE("/expenses/:id", expenses.DeleteExpenseByID)
	}

	// Expense pages
	views := r.Group("/")
	views.Use(auth.Middleware())
	{
		views.GET("/expenses", func(c *gin.Context) {
			userID := auth.GetUserIDFromCookie(c)
			expenses, err := helpers.GetAllExpensesHelper(userID, c)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error fetching expenses")
				return
			}
			page := renderer.New(c.Request.Context(), http.StatusOK, pages.Dashboard(expenses))
			c.Render(http.StatusOK, page)
		})
		views.GET("/expenses/:id", func(c *gin.Context) {
			userID := auth.GetUserIDFromCookie(c)
			id := c.Param("id")
			expense, err := helpers.GetExpenseByIDHelper(id, userID)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error fetching expense")
				return
			}
			page := renderer.New(c.Request.Context(), http.StatusOK, pages.ExpenseByIDPage(expense))
			c.Render(http.StatusOK, page)
		})
	}

	r.Run(":8080")
}
