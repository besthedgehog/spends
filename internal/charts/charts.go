package charts

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"spends/internal/models"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func CreateCategoryPie(data []models.Expense) []byte {
	sums := make(map[string]float64)
	for _, e := range data {
		sums[e.Category] += e.Amount
	}

	var values []chart.Value
	for category, amount := range sums {
		values = append(values, chart.Value{
			Label: category,
			Value: amount,
		})
	}

	pie := chart.PieChart{
		Title:  "Траты за период времени",
		Width:  800,
		Height: 600,
		Values: values,
	}

	var buf bytes.Buffer
	err := pie.Render(chart.PNG, &buf)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func CreatePriorityPie(data []models.Expense) []byte {
	sums := make(map[int]float64)

	for _, e := range data {
		sums[e.Priority] += e.Amount
	}

	var values []chart.Value
	for priority, amount := range sums {
		values = append(values, chart.Value{
			Label: fmt.Sprintf("Приоритет %d", priority-1), // Преобразуем int в string
			Value: amount,
		})
	}
	pie := chart.PieChart{
		Title:  "Траты по приоритетам",
		Width:  800,
		Height: 600,
		Values: values,
	}
	var buf bytes.Buffer
	err := pie.Render(chart.PNG, &buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Столбчатая диаграмма: траты по дням
func CreateDailyBarChart(data []models.Expense) []byte {
	// Группируем траты по дням
	dailySums := make(map[string]float64)
	for _, e := range data {
		date := e.CreatedAt.Format("2006-01-02")
		dailySums[date] += e.Amount
	}

	// Сортируем даты
	var dates []string
	for date := range dailySums {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Создаем значения для графика
	var bars []chart.Value
	for _, date := range dates {
		// Форматируем дату для отображения
		t, _ := time.Parse("2006-01-02", date)
		label := t.Format("02.01")

		bars = append(bars, chart.Value{
			Label: label,
			Value: dailySums[date],
		})
	}

	barChart := chart.BarChart{
		Title:    "Траты по дням",
		Width:    800,
		Height:   600,
		Bars:     bars,
		BarWidth: 60,
	}

	var buf bytes.Buffer
	err := barChart.Render(chart.PNG, &buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// График накопления по дням
func CreateCumulativeChart(data []models.Expense) []byte {
	// Группируем траты по дням
	dailySums := make(map[time.Time]float64)
	for _, e := range data {
		date := time.Date(e.CreatedAt.Year(), e.CreatedAt.Month(), e.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)
		dailySums[date] += e.Amount
	}

	// Сортируем даты
	var dates []time.Time
	for date := range dailySums {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	// Создаем массивы для графика
	var xValues []time.Time
	var yValues []float64
	cumulative := 0.0

	for _, date := range dates {
		cumulative += dailySums[date]
		xValues = append(xValues, date)
		yValues = append(yValues, cumulative)
	}

	graph := chart.Chart{
		Title:  "Накопление расходов по дням",
		Width:  800,
		Height: 600,
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeValueFormatterWithFormat("02.01"),
		},
		YAxis: chart.YAxis{},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					StrokeColor: drawing.ColorBlue,
					StrokeWidth: 2,
				},
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	var buf bytes.Buffer
	err := graph.Render(chart.PNG, &buf)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// RenderPNG принимает название диаграммы и набор байтов
//
// Сохраняет файл в png
func RenderPNG(name string, buf []byte) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(buf)

	if err != nil {
		panic(err)
	}
}
