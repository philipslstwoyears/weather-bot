package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"tg_bots/dto"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const weatherAPIKey = "90c64572496d9b5e97108df48307cefc" // Замените на ваш API-ключ

func getWeather(city string) (string, error) {
	// Формирование запроса к API
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=ru", city, weatherAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить погоду для города %s", city)
	}

	var weather dto.WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return "", err
	}

	description := weather.Weather[0].Description
	temp := weather.Main.Temp
	return fmt.Sprintf("Погода в %s: %s, температура %.1f°C", city, description, temp), nil
}

func main() {
	botToken := "7609550937:AAFe6JZrFYGuqqRJ2qidIcZtpq__Pl5DjmQ"
	//botToken := os.Getenv("7609550937:AAFe6JZrFYGuqqRJ2qidIcZtpq__Pl5DjmQ")
	if botToken == "" {
		log.Fatal("Укажите TELEGRAM_BOT_TOKEN в переменных окружения")
	}

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatalf("Ошибка получения обновлений: %v", err)
	}
	defer bot.StopLongPolling()

	log.Println("Бот запущен...")

	for update := range updates {
		if update.Message != nil {
			city := strings.TrimSpace(update.Message.Text)
			weatherInfo, err := getWeather(city)
			if err != nil {
				weatherInfo = fmt.Sprintf("Ошибка: %v", err)
			}

			msg := tu.Message(tu.ID(update.Message.Chat.ID), weatherInfo)
			_, _ = bot.SendMessage(msg)
		}
	}
}
