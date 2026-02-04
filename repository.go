package main

// Repository - интерфейс для работы с данными
type Repository interface {
	// Пользователи
	RegisterUser(telegramID int64, username, firstName string) error
	GetUser(telegramID int64) (*User, error)

	// Расходы
	SaveExpense(userID int64, state *UserState) (int, error) // Теперь возвращает ID
	GetExpensesByDate(userID int64, date string) ([]Expense, error)
	GetExpensesByMonth(userID int64) ([]Expense, error)
	DeleteExpense(expenseID int, userID int64) error // Новый метод
	GetExpense(expenseID int) (*Expense, error)      // Новый метод

	// Статистика
	GetStatistics(userID int64) (*Statistics, error)
	GetDayTotal(userID int64) (float64, int, error)
	GetMonthCategoryStats(userID int64) (map[string]CategoryStats, float64, int, error)

	// Управление
	Close() error
}

type CategoryStats struct {
	Category string
	Amount   float64
	Count    int
}
