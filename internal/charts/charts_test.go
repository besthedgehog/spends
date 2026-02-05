package charts

import (
	"spends/internal/models"
	"testing"
	"time"
)

func FakeExpenses() []models.Expense {
	now := time.Now()
	expenses := []models.Expense{
		{
			ID:        1,
			UserID:    101,
			Name:      "Покупка продуктов",
			Amount:    200.50,
			Category:  "Еда",
			Priority:  1,
			CreatedAt: now.AddDate(0, 0, -4), // 4 дня назад
		},
		{
			ID:        2,
			UserID:    102,
			Name:      "Оплата коммунальных услуг",
			Amount:    1500.00,
			Category:  "Коммунальные",
			Priority:  2,
			CreatedAt: now.AddDate(0, 0, -3), // 3 дня назад
		},
		{
			ID:        3,
			UserID:    103,
			Name:      "Минимальный платеж по кредитной карте",
			Amount:    300.00,
			Category:  "Кредиты",
			Priority:  3,
			CreatedAt: now.AddDate(0, 0, -2), // 2 дня назад
		},
		{
			ID:        4,
			UserID:    104,
			Name:      "Ремонт автомобиля",
			Amount:    1200.75,
			Category:  "Транспорт",
			Priority:  2,
			CreatedAt: now.AddDate(0, 0, -1), // вчера
		},
		{
			ID:        5,
			UserID:    105,
			Name:      "Покупка нового ноутбука",
			Amount:    850.99,
			Category:  "Техника",
			Priority:  1,
			CreatedAt: now, // сегодня
		},
	}
	return expenses
}

func TestCreateCumulativeChart(T *testing.T) {
	img1 := CreateCumulativeChart(FakeExpenses())
	RenderPNG("CumulativeChart.png", img1)
}

func TestCreateCategoryPie(T *testing.T) {
	img2 := CreateCategoryPie(FakeExpenses())
	RenderPNG("CategoryPie.png", img2)
}

func TestCreatePriorityPie(T *testing.T) {
	img3 := CreatePriorityPie(FakeExpenses())
	RenderPNG("PriorityPie.png", img3)
}

func TestCreateDailyBarChart(T *testing.T) {
	img4 := CreateDailyBarChart(FakeExpenses())
	RenderPNG("DailyBarChart.png", img4)
}
