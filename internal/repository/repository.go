package repository

import (
	"spends/internal/models"
	"time"
)

// Repository - интерфейс для работы с данными
type Repository interface {
	// Пользователи
	RegisterUser(telegramID int64, username, firstName string) error
	GetUser(telegramID int64) (*models.User, error)

	// Расходы
	SaveExpense(userID int64, state *models.UserState) (int, error) // Теперь возвращает ID
	GetExpensesByDate(userID int64, date string) ([]models.Expense, error)
	GetExpensesByMonth(userID int64) ([]models.Expense, error)
	DeleteExpense(expenseID int, userID int64) error   // Новый метод
	GetExpense(expenseID int) (*models.Expense, error) // Новый метод

	// Экспорт
	GetExpensesByPeriod(userID int64, start, end time.Time) ([]models.Expense, error)
	GetAllExpenses(userID int64) ([]models.Expense, error)

	// Статистика
	GetStatistics(userID int64) (*models.Statistics, error)
	GetDayTotal(userID int64) (float64, int, error)
	GetMonthCategoryStats(userID int64) (map[string]models.CategoryStats, float64, int, error)

	// Управление
	Close() error
}
