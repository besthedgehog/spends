package bot

import (
	"log"
	"spends/internal/repository"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot      *tele.Bot
	handlers *BotHandlers
}

func NewBot(token string, repo repository.Repository) (*Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	keyboards := NewKeyboards()
	handlers := NewBotHandlers(repo, keyboards)

	b := &Bot{
		bot:      bot,
		handlers: handlers,
	}

	b.registerHandlers()
	return b, nil
}

func (b *Bot) registerHandlers() {
	log.Println("📝 Регистрация handlers...")

	// 1. КОМАНДЫ (самые специфичные)
	b.bot.Handle("/start", b.handlers.HandleStart)
	log.Println("✅ Зарегистрирована команда /start")

	// 2. КНОПКИ ГЛАВНОГО МЕНЮ (перед OnText!)
	b.bot.Handle("📊 Статистика", b.handlers.HandleStatsButton)
	b.bot.Handle("📥 Экспорт", b.handlers.HandleExportButton)
	b.bot.Handle("📈 Графики", b.handlers.HandleChartsButton)
	log.Println("✅ Зарегистрированы кнопки меню")

	// 3. КАТЕГОРИИ (перед OnText!)
	categoryHandlers := GetCategoryHandlers()
	for btnText, category := range categoryHandlers {
		b.bot.Handle(btnText, b.handlers.HandleCategory(category))
	}
	log.Printf("✅ Зарегистрировано категорий: %d", len(categoryHandlers))

	// 4. ПРИОРИТЕТЫ (перед OnText!)
	priorityHandlers := GetPriorityHandlers()
	for btnText, priority := range priorityHandlers {
		b.bot.Handle(btnText, b.handlers.HandlePriority(priority))
	}
	log.Printf("✅ Зарегистрировано приоритетов: %d", len(priorityHandlers))

	// 5. CALLBACK КНОПКИ (inline)
	b.bot.Handle("\fstats", b.handlers.HandleInlineCallback)
	b.bot.Handle("\fexport", b.handlers.HandleInlineCallback)
	b.bot.Handle("\fcharts", b.handlers.HandleInlineCallback)
	b.bot.Handle("\fundo_expense", b.handlers.HandleUndoExpense)
	log.Println("✅ Зарегистрированы inline callbacks")

	// 6. ОБЩИЙ ОБРАБОТЧИК ТЕКСТА (В САМОМ КОНЦЕ!!!)
	b.bot.Handle(tele.OnText, b.handlers.HandleText)
	log.Println("✅ Зарегистрирован OnText (последним)")

	log.Println("✨ Все handlers зарегистрированы!")
}

func (b *Bot) Start() {
	log.Println("🤖 Бот запущен и готов к работе!")
	b.bot.Start()
}
