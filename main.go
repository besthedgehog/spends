package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	// Получаем токен
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("❌ Токен не найден! Создай файл .env с TOKEN=your_bot_token")
	}

	// Инициализируем репозиторий
	repo, err := NewSQLiteRepository("./expenses.db")
	if err != nil {
		log.Fatal("❌ Ошибка инициализации БД:", err)
	}
	defer repo.Close()

	log.Println("✅ База данных инициализирована")

	// Создаём и запускаем бота
	bot, err := NewBot(token, repo)
	if err != nil {
		log.Fatal("❌ Ошибка создания бота:", err)
	}

	bot.Start()
}
