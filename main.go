package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// Api для получения данных погоды
	endpoint = "https://api.open-meteo.com/v1/forecast" //latitude=55.7522&longitude=37.6156&hourly=temperature_2m
)

var (
	// врменнной промежуток между которым происходит запрос к API
	IntervalReload = time.Second * 2
)

type WeatherData struct {
	//
	Elevation float64 `json:"elevation"`
	// информация "дата"-"температура"
	Hourly map[string]any `json:"hourly"`
}

// структура приложения
type WPoller struct {
	// канал закрывающий приложение
	closeCh chan struct{}
}

// конструктор WPoller
func NewWPoller() *WPoller {
	return &WPoller{
		closeCh: make(chan struct{}),
	}
}

// start запускает приложение
func (w *WPoller) start() {
	ticker := time.NewTicker(IntervalReload)
	fmt.Println("Weather checking start")
loop:
	for {
		select {
		case <-ticker.C:
			// Получение сырых данных
			data, err := GetWeatherByCoordiantes(55.7522, 37.6156)
			if err != nil {
				log.Fatal(err)
			}
			// обработка сырых данных
			if err := w.handleData(data); err != nil {
				log.Fatal(err)
			}
			fmt.Println(data)
		case <-w.closeCh:
			break loop
		}
	}
	fmt.Println()
}

// close закрывает канал closeCh
func (w *WPoller) close() {
	close(w.closeCh)
}

// handleData обрабатывает сырые данные
func (w *WPoller) handleData(data *WeatherData) error {
	// обработка
	// ...
	return nil
}

func main() {
	NewWPoller().start()
}

func GetWeatherByCoordiantes(lat, long float64) (*WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)
	// Request
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	// запрос к API
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Десериализация данных
	var data WeatherData
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	// fmt.Printf(`%T`, data.Hourly["time"])
	return &data, nil
}
