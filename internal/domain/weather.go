package domain

import (
	"time"
)

type WeatherData struct {
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Time        []string `json:"time"`
	CurTime     string
	Temperature []float64 `json:"temperature_2m"`
	CurTemp     float64
}

func (d *WeatherData) GetTemperatureForNowMoment() {
	// "2024-04-19T00:00" - формат времени от API
	// "2006-01-02T15:04:05Z07:00" = RFC3339
	// index := new(int)
	var index int
	now := time.Now().Format(time.RFC3339)
	currentTime := now[:13]
	for i, elem := range d.Time {
		if currentTime == elem[:13] {
			index = i
			break
		}
	}
	d.CurTime = now
	d.CurTemp = d.Temperature[index]
}
