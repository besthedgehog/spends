package service

import (
	"spends/internal/models"
	"time"
)

// ExpenseService - интерфейс бизнес-логики
type ExpenseService interface {
	// Users
	RegisterUser(telegramID int64, username, firstName string) error

	// Expenses
	SaveExpense(userID int64, state *models.UserState) (int, error)
	DeleteExpense(expenseID int, userID int64) error
	GetExpense(expenseID int) (*models.Expense, error)

	// Queries
	GetAllExpenses(userID int64) ([]models.Expense, error)
	GetExpensesByPeriod(userID int64, start, end time.Time) ([]models.Expense, error)
	GetExpensesByDate(userID int64, date string) ([]models.Expense, error)

	// Stats
	GetStatistics(userID int64) (*models.Statistics, error)
	GetMonthCategoryStats(userID int64) (map[string]models.CategoryStats, float64, int, error)
}
