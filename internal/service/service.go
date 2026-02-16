package service

import (
	"spends/internal/models"
	"spends/internal/repository"
	"time"
)

type expenseService struct {
	repo repository.Repository
}

func NewExpenseService(repo repository.Repository) ExpenseService {
	return &expenseService{repo: repo}
}

func (s *expenseService) RegisterUser(telegramID int64, username, firstName string) error {
	return s.repo.RegisterUser(telegramID, username, firstName)
}

func (s *expenseService) SaveExpense(userID int64, state *models.UserState) (int, error) {
	return s.repo.SaveExpense(userID, state)
}

func (s *expenseService) DeleteExpense(expenseID int, userID int64) error {
	return s.repo.DeleteExpense(expenseID, userID)
}

func (s *expenseService) GetExpense(expenseID int) (*models.Expense, error) {
	return s.repo.GetExpense(expenseID)
}

func (s *expenseService) GetAllExpenses(userID int64) ([]models.Expense, error) {
	return s.repo.GetAllExpenses(userID)
}

func (s *expenseService) GetExpensesByPeriod(userID int64, start, end time.Time) ([]models.Expense, error) {
	return s.repo.GetExpensesByPeriod(userID, start, end)
}

func (s *expenseService) GetExpensesByDate(userID int64, date string) ([]models.Expense, error) {
	return s.repo.GetExpensesByDate(userID, date)
}

func (s *expenseService) GetStatistics(userID int64) (*models.Statistics, error) {
	return s.repo.GetStatistics(userID)
}

func (s *expenseService) GetMonthCategoryStats(userID int64) (map[string]models.CategoryStats, float64, int, error) {
	return s.repo.GetMonthCategoryStats(userID)
}
