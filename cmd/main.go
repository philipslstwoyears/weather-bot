package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"tg_bots/internal"
)

func main() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла: ", err)
	}

	// Получаем токен бота из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Укажите TELEGRAM_BOT_TOKEN в переменных окружения")
	}

	log.Println("Токен бота успешно загружен")

	// Создаем бота
	bot, err := internal.CreateBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем свой Chat ID (замените на нужный ID пользователя)
	//chatID := int64(2039699453) // Введите правильный chatID для отправки приветствия
	//
	//// Отправляем приветственное сообщение
	//welcomeMessage := "Привет! Я готов предоставить информацию о погоде. Введите название города."
	//err = internal.SendWelcomeMessage(bot, chatID, welcomeMessage)
	//if err != nil {
	//	log.Printf("Ошибка отправки приветственного сообщения: %v", err)
	//}

	// Получаем обновления и обрабатываем их
	internal.GetUpdates(bot)
}

// конец
