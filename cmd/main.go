package main

import (
	"github.com/VoC925/WeatherChecker/internal/api"
	"github.com/VoC925/WeatherChecker/pkg/client"
)

type WeatherData struct {
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Time        []string `json:"time"`
	Temperature float64  `json:"temperature_2m"`
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// // сигнал завершения приложения ctrl+c
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	sender := client.NewSenderToEmail("0742")
	newWeatherPoller := api.NewWPoller(sender)
	newWeatherPoller.Start()
	// горутина слушащая сигнал завершения приложения
	// go func() {
	// 	<-c
	// 	newWeatherPoller.Close()
	// 	cancel()
	// }()
	// <-ctx.Done()
}
