package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/VoC925/WeatherChecker/pkg/utils"
)

// структура данных погоды
type WeatherData struct {
	// координаты, полученные по IP
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// слайс дат и соответствующих температур, полученных из запроса к IP
	Time        []string  `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
	// текущая дата и температура
	CurTime string
	CurTemp float64
}

// Реализация интерфейса Stringer
func (d *WeatherData) String() string {
	date := d.CurTime[:10]
	time := d.CurTime[11:]
	return fmt.Sprintf("time: %s %sh | temperature: %.2f", date, time, d.CurTemp)
}

// GetTemperatureForNowMoment - метод, заполняющий поля CurTime и CurTemp
// на основе данных в слайсах Time и Temperature
func (d *WeatherData) GetTemperatureForNowMoment() bool {
	// "2024-04-19T00:00" - формат времени от API
	// "2006-01-02T15:04:05Z07:00" = RFC3339
	// index := new(int)
	now := time.Now().Format(time.RFC3339)
	currentTime := now[:13]

	// поиск с проверкой: есть ли текущее время в слайсе CurTime
	index := utils.BinarySearch(d.Time, strings.Join([]string{currentTime, ":00"}, ""))
	if index == -1 {
		return false
	}

	d.CurTime = currentTime
	d.CurTemp = d.Temperature[index]

	return true
}
