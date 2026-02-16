package models

import "time"

const (
	StateNone = iota
	StateWaitingCategory
	StateWaitingPriority
)

type UserState struct {
	State    int
	Name     string
	Amount   float64
	Category string
	Priority int
}

type User struct {
	ID         int
	TelegramID int64
	Username   string
	FirstName  string
	CreatedAt  time.Time
}

type Expense struct {
	ID        int
	UserID    int64
	Name      string
	Amount    float64
	Category  string
	Priority  int
	CreatedAt time.Time
}

type Statistics struct {
	TodayTotal float64
	WeekTotal  float64
	MonthTotal float64
	TotalSum   float64
	Count      int
}

type CategoryStats struct {
	Category string
	Amount   float64
	Count    int
}
