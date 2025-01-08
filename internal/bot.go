package internal

import (
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"strings"
)

// Карта для хранения языковых предпочтений пользователей
var UserLanguage = make(map[int64]string)

// CreateBot - создает бота и возвращает указатель на него
func CreateBot(botToken string) (*telego.Bot, error) {
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания бота: %v", err)
	}
	return bot, nil
}

// GetUpdates - получает обновления и обрабатывает сообщения от пользователей
func GetUpdates(bot *telego.Bot) {
	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatalf("Ошибка получения обновлений: %v", err)
	}
	defer bot.StopLongPolling()

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			// Если пользователь еще не выбрал язык, предлагаем выбор
			if _, ok := UserLanguage[chatID]; !ok {
				buttons := [][]telego.InlineKeyboardButton{
					{tu.InlineKeyboardButton("English").WithCallbackData("lang_en"), tu.InlineKeyboardButton("Русский").WithCallbackData("lang_ru")},
				}
				msg := tu.Message(tu.ID(chatID), "Choose your language / Выберите язык").WithReplyMarkup(tu.InlineKeyboardGrid(buttons))
				bot.SendMessage(msg)
				continue
			}

			// Получаем город из сообщения и отправляем запрос
			city := strings.TrimSpace(update.Message.Text)
			lang := UserLanguage[chatID]
			weatherInfo, err := GetWeather(city, lang)
			if err != nil {
				weatherInfo = fmt.Sprintf("Error: %v", err)
			}

			// Отправляем ответ
			msg := tu.Message(tu.ID(chatID), weatherInfo)
			bot.SendMessage(msg)
		}

		// Обработка кнопок выбора языка
		if update.CallbackQuery != nil {
			ChatID := update.CallbackQuery.Message.GetChat().ID
			data := update.CallbackQuery.Data
			switch data {
			case "lang_en":
				UserLanguage[ChatID] = "en" // Устанавливаем английский
				bot.SendMessage(tu.Message(tu.ID(ChatID), "Language set to English. Please enter a city name."))
			case "lang_ru":
				UserLanguage[ChatID] = "ru" // Устанавливаем русский
				bot.SendMessage(tu.Message(tu.ID(ChatID), "Язык установлен на русский. Пожалуйста, введите название города."))
			}
			bot.AnswerCallbackQuery(tu.CallbackQuery(update.CallbackQuery.ID))
		}
	}
}
