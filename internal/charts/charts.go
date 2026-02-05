package charts

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"spends/internal/models"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// GenerateMonthCharts создаёт HTML с графиками за месяц (все 4 графика)
func GenerateMonthCharts(expenses []models.Expense) ([]byte, error) {
	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)

	// 1. Круговая диаграмма: Категории
	page.AddCharts(createCategoryPie(expenses, "Расходы по категориям (месяц)"))

	// 2. Круговая диаграмма: Приоритеты
	page.AddCharts(createPriorityPie(expenses, "Приоритеты (месяц)"))

	// 3. Столбчатая диаграмма: Расходы по дням
	page.AddCharts(createDailyBar(expenses))

	// 4. Линейный график: Накопление
	page.AddCharts(createCumulativeLine(expenses))

	var buf bytes.Buffer
	if err := page.Render(io.MultiWriter(&buf)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateAllTimeCharts создаёт HTML с графиками за всё время (только круговые)
func GenerateAllTimeCharts(expenses []models.Expense) ([]byte, error) {
	page := components.NewPage()
	page.SetLayout(components.PageFlexLayout)

	page.AddCharts(createCategoryPie(expenses, "Расходы по категориям (всё время)"))
	page.AddCharts(createPriorityPie(expenses, "Приоритеты (всё время)"))

	var buf bytes.Buffer
	if err := page.Render(io.MultiWriter(&buf)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Адаптация твоего кода

func createCategoryPie(data []models.Expense, title string) *charts.Pie {
	sums := make(map[string]float64)
	for _, e := range data {
		sums[e.Category] += e.Amount
	}

	items := make([]opts.PieData, 0)
	for k, v := range sums {
		items = append(items, opts.PieData{Name: k, Value: v})
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: title}))
	pie.AddSeries("Категории", items).SetSeriesOptions(
		charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Formatter: "{b}: {c} ₽"}),
	)
	return pie
}

func createPriorityPie(data []models.Expense, title string) *charts.Pie {
	sums := make(map[int]float64)
	for _, e := range data {
		sums[e.Priority] += e.Amount
	}

	items := make([]opts.PieData, 0)
	labels := map[int]string{
		0: "0: База (Зеленый)",
		1: "1: Комфорт (Желтый)",
		2: "2: Утечки (Красный)",
	}

	for k, v := range sums {
		name := labels[k]
		if name == "" {
			name = fmt.Sprintf("Приоритет %d", k)
		}
		items = append(items, opts.PieData{Name: name, Value: v})
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: title}))
	pie.AddSeries("Приоритеты", items).SetSeriesOptions(
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"40%", "75%"}}),
	)
	return pie
}

func createDailyBar(data []models.Expense) *charts.Bar {
	// Сортируем по дате
	sort.Slice(data, func(i, j int) bool {
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	})

	daily := make(map[string]float64)
	var dates []string

	// Агрегация по дням
	for _, e := range data {
		d := e.CreatedAt.Format("02.01")
		if _, exists := daily[d]; !exists {
			dates = append(dates, d)
		}
		daily[d] += e.Amount
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Траты по дням"}))

	var values []opts.BarData
	for _, d := range dates {
		values = append(values, opts.BarData{Value: daily[d]})
	}

	bar.SetXAxis(dates).AddSeries("Сумма", values)
	return bar
}

func createCumulativeLine(data []models.Expense) *charts.Line {
	// Сортируем по дате
	sort.Slice(data, func(i, j int) bool {
		return data[i].CreatedAt.Before(data[j].CreatedAt)
	})

	daily := make(map[string]float64)
	var dates []string

	for _, e := range data {
		d := e.CreatedAt.Format("02.01")
		if _, exists := daily[d]; !exists {
			dates = append(dates, d)
		}
		daily[d] += e.Amount
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Накопительный итог"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true), Trigger: "axis"}),
	)

	var cumValues []opts.LineData
	currentSum := 0.0
	for _, d := range dates {
		currentSum += daily[d]
		cumValues = append(cumValues, opts.LineData{Value: currentSum})
	}

	line.SetXAxis(dates).AddSeries("Всего потрачено", cumValues)
	return line
}
