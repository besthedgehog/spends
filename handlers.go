package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

type BotHandlers struct {
	repo       Repository
	keyboards  *Keyboards
	userStates map[int64]*UserState
}

func NewBotHandlers(repo Repository, keyboards *Keyboards) *BotHandlers {
	return &BotHandlers{
		repo:       repo,
		keyboards:  keyboards,
		userStates: make(map[int64]*UserState),
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
		RemoveKeyboard(),
		tele.ModeHTML,
	)
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
			"Используй /day или /month для детальной статистики\n\n"+
			"💡 Отправь следующую трату",
		stats.TodayTotal, stats.WeekTotal, stats.MonthTotal, stats.TotalSum, stats.Count,
	)

	return c.Send(message, tele.ModeHTML)
}

func (h *BotHandlers) HandleDay(c tele.Context) error {
	userID := c.Sender().ID
	expenses, err := h.repo.GetExpensesByDate(userID, "now")
	if err != nil {
		return c.Send("❌ Ошибка: " + err.Error())
	}

	if len(expenses) == 0 {
		return c.Send("📅 За сегодня трат пока нет\n\n💡 Отправь следующую трату")
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
	return c.Send(result.String(), tele.ModeHTML)
}

func (h *BotHandlers) HandleMonth(c tele.Context) error {
	userID := c.Sender().ID
	categoryStats, monthTotal, count, err := h.repo.GetMonthCategoryStats(userID)
	if err != nil {
		return c.Send("❌ Ошибка: " + err.Error())
	}

	if count == 0 {
		return c.Send("📈 За этот месяц трат пока нет\n\n💡 Отправь следующую трату")
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
	return c.Send(result.String(), tele.ModeHTML)
}

func (h *BotHandlers) HandleText(c tele.Context) error {
	text := c.Text()
	userID := c.Sender().ID

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

	h.userStates[userID] = &UserState{
		State:  StateWaitingCategory,
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

		if !exists || state.State != StateWaitingCategory {
			return c.Send("⚠️ Сначала отправь трату в формате: Название Сумма")
		}

		state.Category = category
		state.State = StateWaitingPriority

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

		if !exists || state.State != StateWaitingPriority {
			return c.Send("⚠️ Сначала выбери категорию")
		}

		state.Priority = priority

		// Сохраняем и получаем ID траты
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
				"Отправь следующую трату или /stats для статистики",
			state.Name, state.Amount, state.Category, state.Priority,
		)

		delete(h.userStates, userID)

		// Отправляем с кнопкой отмены
		return c.Send(message, UndoExpenseButton(expenseID), tele.ModeHTML)
	}
}

// Новый обработчик - отмена траты
func (h *BotHandlers) HandleUndoExpense(c tele.Context) error {
	log.Println("🔍 HandleUndoExpense вызван")

	// Проверяем, что это наш callback
	if c.Callback() == nil {
		log.Println("⚠️ Callback is nil")
		return nil
	}

	callback := c.Callback()
	log.Printf("📋 Callback Data: %s, Unique: %s", callback.Data, callback.Unique)

	// Проверяем префикс (он в Unique!)
	if callback.Unique != "undo_expense" {
		log.Printf("⚠️ Неверный префикс: %s", callback.Unique)
		return nil
	}

	// Получаем ID траты (он в Data!)
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

	// Проверяем, что трата принадлежит пользователю
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

	// Удаляем трату
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

	// Отвечаем на callback
	c.Respond(&tele.CallbackResponse{
		Text:      "✅ Трата удалена",
		ShowAlert: false,
	})

	// Удаляем старое сообщение с кнопкой
	err = c.Delete()
	if err != nil {
		log.Printf("⚠️ Не удалось удалить сообщение: %v", err)
	}

	// Отправляем новое сообщение об отмене
	cancelMessage := fmt.Sprintf(
		"🗑 <b>Трата отменена!</b>\n\n"+
			"📝 %s - %.2f ₽\n"+
			"📂 %s | 🎯 %d\n\n"+
			"💡 Отправь следующую трату",
		expense.Name, expense.Amount, expense.Category, expense.Priority,
	)

	return c.Send(cancelMessage, tele.ModeHTML)
}
