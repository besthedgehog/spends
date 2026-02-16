package bot

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Keyboards struct {
	Category *tele.ReplyMarkup
	Priority *tele.ReplyMarkup
	Menu     *tele.ReplyMarkup // Новое меню
}

func NewKeyboards() *Keyboards {
	// Клавиатура категорий (появляется только после внесения траты)
	categoryMenu := &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true, // Скрывается после выбора
	}

	btnSnacks := categoryMenu.Text("🍿 Снеки")
	btnProducts := categoryMenu.Text("🛒 Продукты")
	btnCafe := categoryMenu.Text("☕ Кафе")           // ← ДОБАВЬ
	btnTaxi := categoryMenu.Text("🚕 Такси")          // ← ДОБАВЬ
	btnTransport := categoryMenu.Text("🚌 Транспорт") // ← ДОБАВЬ
	btnCare := categoryMenu.Text("💆 Уход")
	btnHousehold := categoryMenu.Text("🧼 Бытхимия")
	btnFun := categoryMenu.Text("🎮 Развлечения")
	btnPsych := categoryMenu.Text("🧠 Психолог")
	btnAlcohol := categoryMenu.Text("🍷 Алкоголь")
	btnTobacco := categoryMenu.Text("🚬 Табак")
	btnSubs := categoryMenu.Text("📱 Подписки")

	categoryMenu.Reply(
		categoryMenu.Row(btnSnacks, btnProducts),
		categoryMenu.Row(btnCafe, btnTaxi),
		categoryMenu.Row(btnCare, btnHousehold),
		categoryMenu.Row(btnFun, btnPsych),
		categoryMenu.Row(btnAlcohol, btnTobacco),
		categoryMenu.Row(btnSubs, btnTransport),
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

	// НОВОЕ: Главное меню (постоянная клавиатура)
	mainMenu := &tele.ReplyMarkup{
		ResizeKeyboard: true,
	}

	btnStats := mainMenu.Text("📊 Статистика")
	btnExport := mainMenu.Text("📥 Экспорт")
	btnCharts := mainMenu.Text("📈 Графики")

	mainMenu.Reply(
		mainMenu.Row(btnStats),
		mainMenu.Row(btnExport, btnCharts),
	)

	return &Keyboards{
		Category: categoryMenu,
		Priority: priorityMenu,
		Menu:     mainMenu,
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
	btnUndo := markup.Data("❌ Отменить", "undo_expense", fmt.Sprintf("%d", expenseID))
	markup.Inline(
		markup.Row(btnUndo),
	)
	return markup
}

// НОВОЕ: Inline меню для статистики
func StatisticsMenu() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnDay := markup.Data("📅 День", "stats", "day")
	btnMonth := markup.Data("📆 Месяц", "stats", "month")
	btnAll := markup.Data("📊 Всё время", "stats", "all")

	markup.Inline(
		markup.Row(btnDay, btnMonth),
		markup.Row(btnAll),
	)
	return markup
}

// НОВОЕ: Inline меню для экспорта
func ExportMenu() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnWeek := markup.Data("📅 Неделя", "export", "week")
	btnMonth := markup.Data("📆 Месяц", "export", "month")
	btnAll := markup.Data("📊 Всё время", "export", "all")

	markup.Inline(
		markup.Row(btnWeek, btnMonth),
		markup.Row(btnAll),
	)
	return markup
}

// НОВОЕ: Inline меню для графиков
func ChartsMenu() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnMonth := markup.Data("📆 Месяц", "charts", "month")
	btnAll := markup.Data("📊 Всё время", "charts", "all")

	markup.Inline(
		markup.Row(btnMonth, btnAll),
	)
	return markup
}

// GetCategoryHandlers возвращает мапу кнопок на категории
func GetCategoryHandlers() map[string]string {
	return map[string]string{
		"🍿 Снеки":       "снеки",
		"🛒 Продукты":    "продукты",
		"☕ Кафе":        "кафе",
		"🚕 Такси":       "такси",
		"🚌 Транспорт":   "транспорт",
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
