package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Keyboards struct {
	Category *tele.ReplyMarkup
	Priority *tele.ReplyMarkup
}

func NewKeyboards() *Keyboards {
	// Клавиатура категорий
	categoryMenu := &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	btnSnacks := categoryMenu.Text("🍿 Снеки")
	btnProducts := categoryMenu.Text("🛒 Продукты")
	btnCare := categoryMenu.Text("💆 Уход")
	btnHousehold := categoryMenu.Text("🧼 Бытхимия")
	btnFun := categoryMenu.Text("🎮 Развлечения")
	btnPsych := categoryMenu.Text("🧠 Психолог")
	btnAlcohol := categoryMenu.Text("🍷 Алкоголь")
	btnTobacco := categoryMenu.Text("🚬 Табак")
	btnSubs := categoryMenu.Text("📱 Подписки")

	categoryMenu.Reply(
		categoryMenu.Row(btnSnacks, btnProducts),
		categoryMenu.Row(btnCare, btnHousehold),
		categoryMenu.Row(btnFun, btnPsych),
		categoryMenu.Row(btnAlcohol, btnTobacco),
		categoryMenu.Row(btnSubs),
	)

	// Клавиатура приоритетов
	priorityMenu := &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	btn0 := priorityMenu.Text("0️⃣ Приоритет 0")
	btn1 := priorityMenu.Text("1️⃣ Приоритет 1")
	btn2 := priorityMenu.Text("2️⃣ Приоритет 2")

	priorityMenu.Reply(
		priorityMenu.Row(btn0, btn1, btn2),
	)

	return &Keyboards{
		Category: categoryMenu,
		Priority: priorityMenu,
	}
}

func RemoveKeyboard() *tele.ReplyMarkup {
	return &tele.ReplyMarkup{
		RemoveKeyboard: true,
	}
}

// UndoExpenseButton создаёт inline-кнопку для отмены траты
func UndoExpenseButton(expenseID int) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	// Используем уникальный префикс для callback
	btnUndo := markup.Data("❌ Отменить", "undo_expense", fmt.Sprintf("%d", expenseID))
	markup.Inline(
		markup.Row(btnUndo),
	)
	return markup
}

// GetCategoryHandlers возвращает мапу кнопок на категории
func GetCategoryHandlers() map[string]string {
	return map[string]string{
		"🍿 Снеки":       "снеки",
		"🛒 Продукты":    "продукты",
		"💆 Уход":        "уход",
		"🧼 Бытхимия":    "бытхимия",
		"🎮 Развлечения": "развлечения",
		"🧠 Психолог":    "психолог",
		"🍷 Алкоголь":    "алкоголь",
		"🚬 Табак":       "табак",
		"📱 Подписки":    "подписки",
	}
}

// GetPriorityHandlers возвращает мапу кнопок на приоритеты
func GetPriorityHandlers() map[string]int {
	return map[string]int{
		"0️⃣ Приоритет 0": 0,
		"1️⃣ Приоритет 1": 1,
		"2️⃣ Приоритет 2": 2,
	}
}
