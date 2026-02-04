package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &SQLiteRepository{db: db}
	if err := repo.init(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepository) init() error {
	// Таблица пользователей
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			telegram_id INTEGER UNIQUE NOT NULL,
			username TEXT,
			first_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Таблица расходов
	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			amount REAL NOT NULL,
			category TEXT NOT NULL,
			priority INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(telegram_id)
		)
	`)
	return err
}

func (r *SQLiteRepository) RegisterUser(telegramID int64, username, firstName string) error {
	_, err := r.db.Exec(`
		INSERT INTO users (telegram_id, username, first_name)
		VALUES (?, ?, ?)
		ON CONFLICT(telegram_id) DO UPDATE SET
			username = excluded.username,
			first_name = excluded.first_name
	`, telegramID, username, firstName)
	return err
}

func (r *SQLiteRepository) GetUser(telegramID int64) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(`
		SELECT id, telegram_id, username, first_name, created_at
		FROM users WHERE telegram_id = ?
	`, telegramID).Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (r *SQLiteRepository) GetExpensesByDate(userID int64, date string) ([]Expense, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, name, amount, category, priority, created_at
		FROM expenses
		WHERE user_id = ? AND DATE(created_at) = DATE(?)
		ORDER BY created_at DESC
	`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.UserID, &e.Name, &e.Amount, &e.Category, &e.Priority, &e.CreatedAt)
		if err != nil {
			continue
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *SQLiteRepository) GetExpensesByMonth(userID int64) ([]Expense, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, name, amount, category, priority, created_at
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.ID, &e.UserID, &e.Name, &e.Amount, &e.Category, &e.Priority, &e.CreatedAt)
		if err != nil {
			continue
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *SQLiteRepository) GetStatistics(userID int64) (*Statistics, error) {
	stats := &Statistics{}

	// Общая статистика
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0), COUNT(*)
		FROM expenses WHERE user_id = ?
	`, userID).Scan(&stats.TotalSum, &stats.Count)
	if err != nil {
		return nil, err
	}

	// Сегодня
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND DATE(created_at) = DATE('now')
	`, userID).Scan(&stats.TodayTotal)
	if err != nil {
		return nil, err
	}

	// Неделя
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', '-7 days')
	`, userID).Scan(&stats.WeekTotal)
	if err != nil {
		return nil, err
	}

	// Месяц
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
	`, userID).Scan(&stats.MonthTotal)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *SQLiteRepository) GetDayTotal(userID int64) (float64, int, error) {
	var total float64
	var count int
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0), COUNT(*)
		FROM expenses
		WHERE user_id = ? AND DATE(created_at) = DATE('now')
	`, userID).Scan(&total, &count)
	return total, count, err
}

func (r *SQLiteRepository) GetMonthCategoryStats(userID int64) (map[string]CategoryStats, float64, int, error) {
	// Общая сумма за месяц
	var monthTotal float64
	var count int
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0), COUNT(*)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
	`, userID).Scan(&monthTotal, &count)
	if err != nil {
		return nil, 0, 0, err
	}

	// По категориям
	rows, err := r.db.Query(`
		SELECT category, SUM(amount), COUNT(*)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
		GROUP BY category
		ORDER BY SUM(amount) DESC
	`, userID)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	categoryStats := make(map[string]CategoryStats)
	for rows.Next() {
		var cs CategoryStats
		err := rows.Scan(&cs.Category, &cs.Amount, &cs.Count)
		if err != nil {
			continue
		}
		categoryStats[cs.Category] = cs
	}

	return categoryStats, monthTotal, count, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

// Обновляем функцию SaveExpense - теперь возвращает ID
//
// Для того, чтобы следующей кнопкной можно было отменить внесённую трату
func (r *SQLiteRepository) SaveExpense(userID int64, state *UserState) (int, error) {
	result, err := r.db.Exec(`
		INSERT INTO expenses (user_id, name, amount, category, priority)
		VALUES (?, ?, ?, ?, ?)
	`, userID, state.Name, state.Amount, state.Category, state.Priority)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// Новый метод - удаление траты
func (r *SQLiteRepository) DeleteExpense(expenseID int, userID int64) error {
	_, err := r.db.Exec(`
		DELETE FROM expenses
		WHERE id = ? AND user_id = ?
	`, expenseID, userID)
	return err
}

// Новый метод - получение одной траты
func (r *SQLiteRepository) GetExpense(expenseID int) (*Expense, error) {
	expense := &Expense{}
	err := r.db.QueryRow(`
		SELECT id, user_id, name, amount, category, priority, created_at
		FROM expenses WHERE id = ?
	`, expenseID).Scan(
		&expense.ID, &expense.UserID, &expense.Name,
		&expense.Amount, &expense.Category, &expense.Priority,
		&expense.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return expense, nil
}
