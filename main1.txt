package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	_ "modernc.org/sqlite"
)

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

type Expense struct {
	ID        int
	UserID    int64
	Name      string
	Amount    float64
	Category  string
	Priority  int
	CreatedAt time.Time
}

var (
	db         *sql.DB
	userStates = make(map[int64]*UserState)
)

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	// Инициализируем БД
	initDB()
	defer db.Close()

	pref := tele.Settings{
		Token:  getToken(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// ===== КЛАВИАТУРА С КАТЕГОРИЯМИ =====
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

	// ===== КЛАВИАТУРА С ПРИОРИТЕТАМИ =====
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

	// ===== КОМАНДА /start =====
	bot.Handle("/start", func(c tele.Context) error {
		userID := c.Sender().ID
		username := c.Sender().Username
		firstName := c.Sender().FirstName

		// Регистрируем или обновляем пользователя
		err := registerUser(userID, username, firstName)
		if err != nil {
			log.Printf("Ошибка регистрации пользователя: %v", err)
		}

		removeKeyboard := &tele.ReplyMarkup{
			RemoveKeyboard: true,
		}
		return c.Send(
			fmt.Sprintf("👋 Привет, %s!\n\nОтправь трату в формате:\n<b>Название Сумма</b>\n\nНапример: Кофе 135", firstName),
			removeKeyboard,
			tele.ModeHTML,
		)
	})

	// ===== КОМАНДА /stats =====
	bot.Handle("/stats", func(c tele.Context) error {
		userID := c.Sender().ID
		stats, err := getStatistics(userID)
		if err != nil {
			return c.Send("❌ Ошибка получения статистики: " + err.Error())
		}
		return c.Send(stats, tele.ModeHTML)
	})

	// ===== КОМАНДА /day (статистика за день) =====
	bot.Handle("/day", func(c tele.Context) error {
		userID := c.Sender().ID
		stats, err := getDayStatistics(userID)
		if err != nil {
			return c.Send("❌ Ошибка: " + err.Error())
		}
		return c.Send(stats, tele.ModeHTML)
	})

	// ===== КОМАНДА /month (статистика за месяц) =====
	bot.Handle("/month", func(c tele.Context) error {
		userID := c.Sender().ID
		stats, err := getMonthStatistics(userID)
		if err != nil {
			return c.Send("❌ Ошибка: " + err.Error())
		}
		return c.Send(stats, tele.ModeHTML)
	})

	// ===== ОБРАБОТКА КАТЕГОРИЙ =====
	categoryHandlers := map[string]string{
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

	for btnText, category := range categoryHandlers {
		category := category
		bot.Handle(btnText, func(c tele.Context) error {
			userID := c.Sender().ID
			state, exists := userStates[userID]

			if !exists || state.State != StateWaitingCategory {
				return c.Send("⚠️ Сначала отправь трату в формате: Название Сумма")
			}

			state.Category = category
			state.State = StateWaitingPriority

			return c.Send(
				fmt.Sprintf("✅ Категория: <b>%s</b>\n\nВыбери приоритет:", category),
				priorityMenu,
				tele.ModeHTML,
			)
		})
	}

	// ===== ОБРАБОТКА ПРИОРИТЕТОВ =====
	priorityHandlers := map[string]int{
		"0️⃣ Приоритет 0": 0,
		"1️⃣ Приоритет 1": 1,
		"2️⃣ Приоритет 2": 2,
	}

	for btnText, priority := range priorityHandlers {
		priority := priority
		bot.Handle(btnText, func(c tele.Context) error {
			userID := c.Sender().ID
			state, exists := userStates[userID]

			if !exists || state.State != StateWaitingPriority {
				return c.Send("⚠️ Сначала выбери категорию")
			}

			state.Priority = priority

			// Сохраняем трату в БД
			err := saveExpense(userID, state)
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
				state.Name,
				state.Amount,
				state.Category,
				state.Priority,
			)

			delete(userStates, userID)

			removeKeyboard := &tele.ReplyMarkup{
				RemoveKeyboard: true,
			}

			return c.Send(message, removeKeyboard, tele.ModeHTML)
		})
	}

	// ===== ОБРАБОТКА ТЕКСТОВЫХ СООБЩЕНИЙ =====
	bot.Handle(tele.OnText, func(c tele.Context) error {
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

		userStates[userID] = &UserState{
			State:  StateWaitingCategory,
			Name:   name,
			Amount: amount,
		}

		return c.Send(
			fmt.Sprintf("📝 Трата: <b>%s</b>\n💰 Сумма: <b>%.2f ₽</b>\n\nВыбери категорию:", name, amount),
			categoryMenu,
			tele.ModeHTML,
		)
	})

	log.Println("🤖 Бот запущен...")
	bot.Start()
}

// ===== ФУНКЦИИ ДЛЯ РАБОТЫ С БД =====

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}

	// Таблица пользователей
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			telegram_id INTEGER UNIQUE NOT NULL,
			username TEXT,
			first_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Таблица расходов
	_, err = db.Exec(`
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
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ База данных инициализирована")
}

func getToken() string {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("❌ Токен не найден! Создай файл .env с TOKEN=your_bot_token")
	}
	return token
}

func registerUser(telegramID int64, username, firstName string) error {
	_, err := db.Exec(`
		INSERT INTO users (telegram_id, username, first_name)
		VALUES (?, ?, ?)
		ON CONFLICT(telegram_id) DO UPDATE SET
			username = excluded.username,
			first_name = excluded.first_name
	`, telegramID, username, firstName)

	return err
}

func saveExpense(userID int64, state *UserState) error {
	_, err := db.Exec(`
		INSERT INTO expenses (user_id, name, amount, category, priority)
		VALUES (?, ?, ?, ?, ?)
	`, userID, state.Name, state.Amount, state.Category, state.Priority)

	return err
}

func getStatistics(userID int64) (string, error) {
	// Общая статистика
	var total float64
	var count int

	err := db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0), COUNT(*)
		FROM expenses
		WHERE user_id = ?
	`, userID).Scan(&total, &count)

	if err != nil {
		return "", err
	}

	// Статистика за сегодня
	var todayTotal float64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND DATE(created_at) = DATE('now')
	`, userID).Scan(&todayTotal)

	if err != nil {
		return "", err
	}

	// Статистика за неделю
	var weekTotal float64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', '-7 days')
	`, userID).Scan(&weekTotal)

	if err != nil {
		return "", err
	}

	// Статистика за месяц
	var monthTotal float64
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
	`, userID).Scan(&monthTotal)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"📊 <b>Статистика</b>\n\n"+
			"📅 Сегодня: %.2f ₽\n"+
			"📆 Неделя: %.2f ₽\n"+
			"📈 Месяц: %.2f ₽\n"+
			"💼 Всего: %.2f ₽\n"+
			"📝 Записей: %d\n\n"+
			"Используй /day или /month для детальной статистики \n Отправь следующую трату",
		todayTotal, weekTotal, monthTotal, total, count,
	), nil
}

func getDayStatistics(userID int64) (string, error) {
	rows, err := db.Query(`
		SELECT name, amount, category, priority, created_at
		FROM expenses
		WHERE user_id = ? AND DATE(created_at) = DATE('now')
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result strings.Builder
	result.WriteString("📅 <b>Траты за сегодня</b>\n\n")

	var total float64
	count := 0

	for rows.Next() {
		var e Expense
		err := rows.Scan(&e.Name, &e.Amount, &e.Category, &e.Priority, &e.CreatedAt)
		if err != nil {
			continue
		}

		result.WriteString(fmt.Sprintf(
			"• %s - <b>%.2f ₽</b>\n  📂 %s | 🎯 %d | ⏰ %s\n\n",
			e.Name,
			e.Amount,
			e.Category,
			e.Priority,
			e.CreatedAt.Format("15:04"),
		))

		total += e.Amount
		count++
	}

	if count == 0 {
		return "📅 За сегодня трат пока нет", nil
	}

	result.WriteString(fmt.Sprintf("💰 <b>Итого: %.2f ₽</b> (%d трат) \n Отправь следующую трату", total, count))
	return result.String(), nil
}

func getMonthStatistics(userID int64) (string, error) {
	// Общая сумма за месяц
	var monthTotal float64
	var count int
	err := db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0), COUNT(*)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
	`, userID).Scan(&monthTotal, &count)

	if err != nil {
		return "", err
	}

	if count == 0 {
		return "📈 За этот месяц трат пока нет", nil
	}

	// Статистика по категориям
	rows, err := db.Query(`
		SELECT category, SUM(amount), COUNT(*)
		FROM expenses
		WHERE user_id = ? AND created_at >= DATE('now', 'start of month')
		GROUP BY category
		ORDER BY SUM(amount) DESC
	`, userID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var result strings.Builder
	result.WriteString(fmt.Sprintf("📈 <b>Статистика за месяц</b>\n\n💰 Всего: <b>%.2f ₽</b> (%d трат)\n\n<b>По категориям:</b>\n", monthTotal, count))

	for rows.Next() {
		var category string
		var amount float64
		var catCount int
		err := rows.Scan(&category, &amount, &catCount)
		if err != nil {
			continue
		}

		percentage := (amount / monthTotal) * 100
		result.WriteString(fmt.Sprintf(
			"📂 %s: %.2f ₽ (%.1f%%, %d трат)\n",
			category,
			amount,
			percentage,
			catCount,
		))
	}

	// Добавляем сообщение в конец
	result.WriteString("\n💡 Добавь следующую трату")

	return result.String(), nil
}
