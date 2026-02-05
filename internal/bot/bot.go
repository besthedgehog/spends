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
	// Команды
	b.bot.Handle("/start", b.handlers.HandleStart)
	b.bot.Handle("/stats", b.handlers.HandleStats)
	b.bot.Handle("/day", b.handlers.HandleDay)
	b.bot.Handle("/month", b.handlers.HandleMonth)

	// Экспорт CSV
	b.bot.Handle("/export_week", b.handlers.HandleExportWeek)
	b.bot.Handle("/export_month", b.handlers.HandleExportMonth)
	b.bot.Handle("/export_all", b.handlers.HandleExportAll)

	// Графики
	b.bot.Handle("/charts_month", b.handlers.HandleChartsMonth)
	b.bot.Handle("/charts_all", b.handlers.HandleChartsAll)

	// Категории
	for btnText, category := range GetCategoryHandlers() {
		b.bot.Handle(btnText, b.handlers.HandleCategory(category))
	}

	// Приоритеты
	for btnText, priority := range GetPriorityHandlers() {
		b.bot.Handle(btnText, b.handlers.HandlePriority(priority))
	}

	// Callback для отмены траты - регистрируем по префиксу
	b.bot.Handle("\fundo_expense", b.handlers.HandleUndoExpense)

	// Текстовые сообщения
	b.bot.Handle(tele.OnText, b.handlers.HandleText)
}

func (b *Bot) Start() {
	log.Println("🤖 Бот запущен...")
	b.bot.Start()
}
