package bot

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"spends/internal/charts"
	"spends/internal/export"
	"spends/internal/models"
	"spends/internal/repository"

	tele "gopkg.in/telebot.v4"
)

type BotHandlers struct {
	repo       repository.Repository
	keyboards  *Keyboards
	userStates map[int64]*models.UserState
}

func NewBotHandlers(repo repository.Repository, keyboards *Keyboards) *BotHandlers {
	return &BotHandlers{
		repo:       repo,
		keyboards:  keyboards,
		userStates: make(map[int64]*models.UserState),
	}
}

func (h *BotHandlers) HandleStart(c tele.Context) error {
	userID := c.Sender().ID
	username := c.Sender().Username
	firstName := c.Sender().FirstName

	err := h.repo.RegisterUser(userID, username, firstName)
	if err != nil {
		return c.Send("❌ Ошибка регистрации: " + err.Error())
	}

	return c.Send(
		fmt.Sprintf("👋 Привет, %s!\n\nОтправь трату в формате:\n<b>Название Сумма</b>\n\nНапример: Кофе 135", firstName),
		h.keyboards.Menu,
		tele.ModeHTML,
	)
}

func (h *BotHandlers) HandleStatsButton(c tele.Context) error {
	log.Println("🎯 HandleStatsButton ВЫЗВАН!")
	if err := c.Send("Выбери период:", StatisticsMenu()); err != nil {
		return err
	}

	return c.Send("Отправь следующую трату", h.keyboards.Menu)

}

func (h *BotHandlers) HandleExportButton(c tele.Context) error {

	if err := c.Send("Выбери период для экспорта:", ExportMenu()); err != nil {
		return err
	}

	return c.Send("Отправь следующую трату", h.keyboards.Menu)
}

func (h *BotHandlers) HandleChartsButton(c tele.Context) error {

	return h.GetAllCharts(c)
}

// НОВЫЕ обработчики inline кнопок

func (h *BotHandlers) HandleInlineCallback(c tele.Context) error {
	if c.Callback() == nil {
		return nil
	}

	callback := c.Callback()

	switch callback.Unique {
	case "stats":
		return h.handleStatsCallback(c, callback.Data)
	case "export":
		return h.handleExportCallback(c, callback.Data)
	case "charts":
		return h.handleChartsCallback(c, callback.Data)
	}

	return nil
}

// ФИК 2: Убрали неиспользуемую переменную
func (h *BotHandlers) handleStatsCallback(c tele.Context, period string) error {
	c.Respond(&tele.CallbackResponse{Text: "⏳ Загружаю..."})

	switch period {
	case "day":
		return h.HandleDay(c)
	case "month":
		return h.HandleMonth(c)
	case "all":
		return h.HandleStats(c)
	}
	return nil
}

func (h *BotHandlers) handleExportCallback(c tele.Context, period string) error {
	c.Respond(&tele.CallbackResponse{Text: "⏳ Создаю файл..."})

	switch period {
	case "week":
		return h.HandleExportWeek(c)
	case "month":
		return h.HandleExportMonth(c)
	case "all":
		return h.HandleExportAll(c)
	}
	return nil
}

func (h *BotHandlers) handleChartsCallback(c tele.Context, period string) error {
	c.Respond(&tele.CallbackResponse{Text: "⏳ Создаю графики..."})

	return h.GetAllCharts(c)
	// switch period {
	// case "month":
	// 	return h.HandleChartsMonth(c)
	// case "all":
	// 	return h.HandleChartsAll(c)
	// }
	// return nil
}

func (h *BotHandlers) HandleStats(c tele.Context) error {
	userID := c.Sender().ID
	stats, err := h.repo.GetStatistics(userID)
	if err != nil {
		return c.Send("❌ Ошибка получения статистики: " + err.Error())
	}

	message := fmt.Sprintf(
		"📊 <b>Статистика</b>\n\n"+
			"📅 Сегодня: %.2f ₽\n"+
			"📆 Неделя: %.2f ₽\n"+
			"📈 Месяц: %.2f ₽\n"+
			"💼 Всего: %.2f ₽\n"+
			"📝 Записей: %d\n\n"+
			"💡 Отправь следующую трату",
		stats.TodayTotal, stats.WeekTotal, stats.MonthTotal, stats.TotalSum, stats.Count,
	)

	return c.Send(message, h.keyboards.Menu, tele.ModeHTML) // ← Добавили меню
}

func (h *BotHandlers) HandleDay(c tele.Context) error {
	userID := c.Sender().ID
	expenses, err := h.repo.GetExpensesByDate(userID, "now")
	if err != nil {
		return c.Send("❌ Ошибка: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📅 За сегодня трат пока нет\n\n💡 Отправь следующую трату", h.keyboards.Menu) // ← Добавили меню
	}

	var result strings.Builder
	result.WriteString("📅 <b>Траты за сегодня</b>\n\n")

	var total float64
	for _, e := range expenses {
		result.WriteString(fmt.Sprintf(
			"• %s - <b>%.2f ₽</b>\n  📂 %s | 🎯 %d | ⏰ %s\n\n",
			e.Name, e.Amount, e.Category, e.Priority, e.CreatedAt.Format("15:04"),
		))
		total += e.Amount
	}

	result.WriteString(fmt.Sprintf("💰 <b>Итого: %.2f ₽</b> (%d трат)\n\n💡 Отправь следующую трату", total, len(expenses)))
	return c.Send(result.String(), h.keyboards.Menu, tele.ModeHTML) // ← Добавили меню
}

func (h *BotHandlers) HandleMonth(c tele.Context) error {
	userID := c.Sender().ID
	categoryStats, monthTotal, count, err := h.repo.GetMonthCategoryStats(userID)
	if err != nil {
		return c.Send("❌ Ошибка: " + err.Error())
	}

	if count == 0 {
		return c.Send("📈 За этот месяц трат пока нет\n\n💡 Отправь следующую трату", h.keyboards.Menu) // ← Добавили меню
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📈 <b>Статистика за месяц</b>\n\n💰 Всего: <b>%.2f ₽</b> (%d трат)\n\n<b>По категориям:</b>\n", monthTotal, count))

	for _, cs := range categoryStats {
		percentage := (cs.Amount / monthTotal) * 100
		result.WriteString(fmt.Sprintf(
			"📂 %s: %.2f ₽ (%.1f%%, %d трат)\n",
			cs.Category, cs.Amount, percentage, cs.Count,
		))
	}

	result.WriteString("\n💡 Добавь следующую трату")
	return c.Send(result.String(), h.keyboards.Menu, tele.ModeHTML) // ← Добавили меню
}

func (h *BotHandlers) HandleText(c tele.Context) error {
	text := c.Text()
	userID := c.Sender().ID

	// КРИТИЧЕСКИ ВАЖНО: Игнорируем кнопки меню
	menuButtons := []string{"📊 Статистика", "📥 Экспорт", "📈 Графики"}
	for _, btn := range menuButtons {
		if text == btn {
			log.Printf("⚠️ OnText проигнорировал кнопку меню: %s", text)
			return nil // Пропускаем - обработается специфичным handler
		}
	}

	// Игнорируем категории
	for btnText := range GetCategoryHandlers() {
		if text == btnText {
			log.Printf("⚠️ OnText проигнорировал категорию: %s", text)
			return nil
		}
	}

	// Игнорируем приоритеты
	for btnText := range GetPriorityHandlers() {
		if text == btnText {
			log.Printf("⚠️ OnText проигнорировал приоритет: %s", text)
			return nil
		}
	}

	log.Printf("📝 OnText обрабатывает: %s", text)

	parts := strings.Fields(text)
	if len(parts) < 2 {
		return c.Send(
			"⚠️ Неверный формат. Используй:\n<b>Название Сумма</b>\n\nНапример: Кофе 135",
			tele.ModeHTML,
		)
	}

	amountStr := parts[len(parts)-1]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.Send("⚠️ Не могу распознать сумму. Убедись, что последнее слово - это число")
	}

	name := strings.Join(parts[:len(parts)-1], " ")

	h.userStates[userID] = &models.UserState{
		State:  models.StateWaitingCategory,
		Name:   name,
		Amount: amount,
	}

	return c.Send(
		fmt.Sprintf("📝 Трата: <b>%s</b>\n💰 Сумма: <b>%.2f ₽</b>\n\nВыбери категорию:", name, amount),
		h.keyboards.Category,
		tele.ModeHTML,
	)
}

func (h *BotHandlers) HandleCategory(category string) tele.HandlerFunc {
	return func(c tele.Context) error {
		userID := c.Sender().ID
		state, exists := h.userStates[userID]

		if !exists || state.State != models.StateWaitingCategory {
			return c.Send("⚠️ Сначала отправь трату в формате: Название Сумма")
		}

		state.Category = category
		state.State = models.StateWaitingPriority

		return c.Send(
			fmt.Sprintf("✅ Категория: <b>%s</b>\n\nВыбери приоритет:", category),
			h.keyboards.Priority,
			tele.ModeHTML,
		)
	}
}

func (h *BotHandlers) HandlePriority(priority int) tele.HandlerFunc {
	return func(c tele.Context) error {
		userID := c.Sender().ID
		state, exists := h.userStates[userID]

		if !exists || state.State != models.StateWaitingPriority {
			return c.Send("⚠️ Сначала выбери категорию")
		}

		state.Priority = priority

		expenseID, err := h.repo.SaveExpense(userID, state)
		if err != nil {
			return c.Send("❌ Ошибка при сохранении: " + err.Error())
		}

		message := fmt.Sprintf(
			"💾 <b>Трата сохранена!</b>\n\n"+
				"📝 Название: %s\n"+
				"💰 Сумма: %.2f ₽\n"+
				"📂 Категория: %s\n"+
				"🎯 Приоритет: %d\n\n"+
				"Отправь следующую трату или используй меню ⬇️",
			state.Name, state.Amount, state.Category, state.Priority,
		)

		delete(h.userStates, userID)

		if err := c.Send(message, UndoExpenseButton(expenseID), tele.ModeHTML); err != nil {
			return err
		}

		return c.Send("Используй меню ⬇️", h.keyboards.Menu)
	}
}

func (h *BotHandlers) HandleUndoExpense(c tele.Context) error {
	log.Println("🔍 HandleUndoExpense вызван")

	if c.Callback() == nil {
		log.Println("⚠️ Callback is nil")
		return nil
	}

	callback := c.Callback()
	log.Printf("📋 Callback Data: %s, Unique: %s", callback.Data, callback.Unique)

	if callback.Unique != "undo_expense" {
		log.Printf("⚠️ Неверный префикс: %s", callback.Unique)
		return nil
	}

	expenseIDStr := callback.Data
	log.Printf("🆔 Expense ID string: %s", expenseIDStr)

	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		log.Printf("❌ Ошибка парсинга ID: %v", err)
		c.Respond(&tele.CallbackResponse{
			Text:      "❌ Ошибка: неверный ID",
			ShowAlert: false,
		})
		return nil
	}

	log.Printf("✅ Parsed expense ID: %d", expenseID)

	userID := c.Sender().ID

	expense, err := h.repo.GetExpense(expenseID)
	if err != nil {
		log.Printf("❌ Ошибка получения траты: %v", err)
		c.Respond(&tele.CallbackResponse{
			Text:      "❌ Трата не найдена",
			ShowAlert: false,
		})
		return nil
	}

	log.Printf("📊 Expense found: %+v", expense)

	if expense.UserID != userID {
		log.Printf("❌ User mismatch: %d != %d", expense.UserID, userID)
		c.Respond(&tele.CallbackResponse{
			Text:      "❌ Это не твоя трата",
			ShowAlert: true,
		})
		return nil
	}

	err = h.repo.DeleteExpense(expenseID, userID)
	if err != nil {
		log.Printf("❌ Ошибка удаления: %v", err)
		c.Respond(&tele.CallbackResponse{
			Text:      "❌ Ошибка удаления",
			ShowAlert: false,
		})
		return nil
	}

	log.Println("✅ Трата успешно удалена")

	c.Respond(&tele.CallbackResponse{
		Text:      "✅ Трата удалена",
		ShowAlert: false,
	})

	err = c.Delete()
	if err != nil {
		log.Printf("⚠️ Не удалось удалить сообщение: %v", err)
	}

	cancelMessage := fmt.Sprintf(
		"🗑 <b>Трата отменена!</b>\n\n"+
			"📝 %s - %.2f ₽\n"+
			"📂 %s | 🎯 %d\n\n"+
			"💡 Отправь следующую трату",
		expense.Name, expense.Amount, expense.Category, expense.Priority,
	)

	// Отправляем сообщение
	if err := c.Send(cancelMessage, tele.ModeHTML); err != nil {
		return err
	}

	// ВАЖНО: Возвращаем главное меню
	return c.Send("Используй меню:", h.keyboards.Menu)
}

// ===== ЭКСПОРТ CSV =====

func (h *BotHandlers) HandleExportWeek(c tele.Context) error {
	userID := c.Sender().ID

	end := time.Now()
	start := end.AddDate(0, 0, -7)

	expenses, err := h.repo.GetExpensesByPeriod(userID, start, end)
	if err != nil {
		return c.Send("❌ Ошибка получения данных: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📭 За последнюю неделю трат нет")
	}

	csvData, err := export.ExportToCSV(expenses)
	if err != nil {
		return c.Send("❌ Ошибка создания CSV: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.FromReader(bytes.NewReader(csvData)),
		FileName: fmt.Sprintf("expenses_week_%s.csv", time.Now().Format("2006-01-02")),
		MIME:     "text/csv",
	}

	if err := c.Send(doc); err != nil {
		return err
	}

	return c.Send("Отправь следующую трату", h.keyboards.Menu)

}

func (h *BotHandlers) HandleExportMonth(c tele.Context) error {
	userID := c.Sender().ID

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := now

	expenses, err := h.repo.GetExpensesByPeriod(userID, start, end)
	if err != nil {
		return c.Send("❌ Ошибка получения данных: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📭 За этот месяц трат нет")
	}

	csvData, err := export.ExportToCSV(expenses)
	if err != nil {
		return c.Send("❌ Ошибка создания CSV: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.FromReader(bytes.NewReader(csvData)),
		FileName: fmt.Sprintf("expenses_month_%s.csv", time.Now().Format("2006-01")),
		MIME:     "text/csv",
	}

	if err := c.Send(doc); err != nil {
		return err
	}

	return c.Send("Отправь следующую трату", h.keyboards.Menu)

}

func (h *BotHandlers) HandleExportAll(c tele.Context) error {
	userID := c.Sender().ID

	expenses, err := h.repo.GetAllExpenses(userID)
	if err != nil {
		return c.Send("❌ Ошибка получения данных: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📭 У тебя пока нет трат")
	}

	csvData, err := export.ExportToCSV(expenses)
	if err != nil {
		return c.Send("❌ Ошибка создания CSV: " + err.Error())
	}

	doc := &tele.Document{
		File:     tele.FromReader(bytes.NewReader(csvData)),
		FileName: fmt.Sprintf("expenses_all_%s.csv", time.Now().Format("2006-01-02")),
		MIME:     "text/csv",
	}

	if err := c.Send(doc); err != nil {
		return err
	}

	return c.Send("Отправь следующую трату", h.keyboards.Menu)

}

func (h *BotHandlers) GetAllCharts(c tele.Context) error {

	ChartsNames := struct {
		CumulativeChart string
		CategoryPie     string
		PriorityPie     string
		DailyBarChart   string
	}{
		CumulativeChart: "CumulativeChart.png",
		CategoryPie:     "CategoryPie.png",
		PriorityPie:     "PriorityPie.png",
		DailyBarChart:   "DailyBarChart.png",
	}

	userID := c.Sender().ID

	expenses, err := h.repo.GetAllExpenses(userID)
	if err != nil {
		return c.Send("❌ Ошибка получения данных: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📭 У тебя пока нет трат")
	}

	{
		img1 := charts.CreateCumulativeChart(expenses)
		img2 := charts.CreateCategoryPie(expenses)
		img3 := charts.CreatePriorityPie(expenses)
		img4 := charts.CreateDailyBarChart(expenses)

		charts.RenderPNG(ChartsNames.CumulativeChart, img1)
		charts.RenderPNG(ChartsNames.CategoryPie, img2)
		charts.RenderPNG(ChartsNames.PriorityPie, img3)
		charts.RenderPNG(ChartsNames.DailyBarChart, img4)

		c.Send(&tele.Photo{File: tele.FromReader(bytes.NewReader(img1)), Caption: "📈 Накопление"})
		c.Send(&tele.Photo{File: tele.FromReader(bytes.NewReader(img2)), Caption: "🥧 Категории"})
		c.Send(&tele.Photo{File: tele.FromReader(bytes.NewReader(img3)), Caption: "🎯 Приоритеты"})
		c.Send(&tele.Photo{File: tele.FromReader(bytes.NewReader(img4)), Caption: "📊 По дням"})
	}

	defer func() {
		os.Remove(ChartsNames.CumulativeChart)
		os.Remove(ChartsNames.CategoryPie)
		os.Remove(ChartsNames.PriorityPie)
		os.Remove(ChartsNames.DailyBarChart)
	}()

	return c.Send("Отправь следующую трату", h.keyboards.Menu)
}
