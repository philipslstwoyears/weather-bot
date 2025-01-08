package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// преобразую температура в фаренгейты для англичан
func celsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}

const weatherAPIKey = "90c64572496d9b5e97108df48307cefc" // Замените на ваш API-ключ

func GetWeather(city, lang string) (string, error) {
	// Формирование запроса к API
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s&lang=%s", city, weatherAPIKey, lang)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось получить погоду для города %s", city)
	}

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return "", err
	}

	// В зависимости от языка, форматируем ответ
	if lang == "en" {
		// Преобразуем температуру в Фаренгейты
		tempFahrenheit := celsiusToFahrenheit(weather.Main.Temp)
		return fmt.Sprintf("Weather in %s: %s, Temperature %.1f°F", city, weather.Weather[0].Description, tempFahrenheit), nil
	}

	// Если русский язык, показываем температуру в Цельсиях
	return fmt.Sprintf("Погода в %s: %s, температура %.1f°C", city, weather.Weather[0].Description, weather.Main.Temp), nil
}
