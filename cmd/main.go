package main

import (
	"os"
	"os/signal"
	"sync"

	"github.com/VoC925/WeatherChecker/internal/api"
	"github.com/VoC925/WeatherChecker/pkg/client"
)

func main() {
	// сигнал завершения приложения ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// sender := client.NewHermesSender("ewg.covaleov1999@yandex.ru", "Broadcast")
	sender := client.NewSenderToEmail("0712")
	// экземпляр App прогноза погоды
	newWeatherPoller := api.NewWPoller("62.181.13.102", sender) // 62.181.13.102 62.181.51.102

	// горутина запускающая приложение
	go func() {
		defer wg.Done()
		newWeatherPoller.Start()
	}()

	// горутина слушащая сигнал завершения приложения
	go func() {
		defer wg.Done()
		<-c
		newWeatherPoller.Close()
	}()

	wg.Wait()
}
