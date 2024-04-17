package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Time        []string `json:"time"`
	Temperature float64  `json:"temperature_2m"`
}

// структура приложения
type WPoller struct {
	// интерфейс, куда отправлять данные погоды
	sender Sender
	// канал закрывающий приложение
	closeCh chan struct{}
}

// конструктор WPoller
func NewWPoller(sender Sender) *WPoller {
	return &WPoller{
		sender:  sender,
		closeCh: make(chan struct{}),
	}
}

// start запускает приложение
func (w *WPoller) start() {
	ticker := time.NewTicker(IntervalReload)
	fmt.Println("------------Weather checking start------------")
loop:
	for {
		select {
		case <-ticker.C:
			// Получение сырых данных
			weatherData, err := w.GetWeatherByCoordiantes(55.7522, 37.6156)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(weatherData)
		case <-w.closeCh:
			break loop
		}
	}
	fmt.Println("------------Weather checking finished------------")
}

// close закрывает канал closeCh
func (w *WPoller) close() {
	close(w.closeCh)
}

// handleData обрабатывает сырые данные
func (w *WPoller) decodeData(r io.Reader) (*WeatherData, error) {
	var (
		wData WeatherData
		m     map[string]*json.RawMessage
	)
	// декодирование в мапу m
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	for key, value := range m {
		switch key {
		case "latitude":
			if err := json.Unmarshal(*value, &wData.Latitude); err != nil {
				return nil, err
			}
		case "longitude":
			if err := json.Unmarshal(*value, &wData.Longitude); err != nil {
				return nil, err
			}
		case "hourly":
			if err := w.desirializeTimeAndTemperature(&wData, value); err != nil {
				return nil, err
			}
		default:
			continue
		}
	}
	return &wData, nil
}

func (w *WPoller) desirializeTimeAndTemperature(wData *WeatherData, data *json.RawMessage) error {
	var m map[string][]any
	if err := json.Unmarshal(*data, &m); err != nil {
		return err
	}
	for key, value := range m {
		switch key {
		case "time":
			sliceOfDate := make([]string, 0, len(value))
			for index, elem := range value {
				sliceOfDate[index] = elem.(string)
			}
			wData.Time = sliceOfDate
		case "temperature_2m":
			sliceOfDate := make([]string, 0, len(value))
			for index, elem := range value {
				sliceOfDate[index] = elem.(string)
			}
			wData.Time = sliceOfDate
		}
	}
	return nil
}

func (w *WPoller) GetWeatherByCoordiantes(lat, long float64) (*WeatherData, error) {
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

	weatherD, err := w.decodeData(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return weatherD, nil
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	sender := NewSenderToEmail("0742")
	newWeatherPoller := NewWPoller(sender)
	newWeatherPoller.start()
}
