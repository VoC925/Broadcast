package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/VoC925/WeatherChecker/internal/domain"
	"github.com/VoC925/WeatherChecker/pkg/client"
	"github.com/VoC925/WeatherChecker/pkg/utils"
)

const (
	// Api для получения данных погоды
	endpoint = "https://api.open-meteo.com/v1/forecast" //latitude=55.7522&longitude=37.6156&hourly=temperature_2m
)

var (
	// временнной промежуток между которым происходит запрос к API
	IntervalReload = time.Second * 2
)

// WPoller структура приложения, получающая температуру и отправляющая sender.
type WPoller struct {
	sender  client.Sender // интерфейс, куда отправлять данные погоды
	closeCh chan struct{} // канал закрывающий приложение
}

// NewWPoller - конструктор структуры WPoller.
func NewWPoller(sender client.Sender) *WPoller {
	return &WPoller{
		sender:  sender,
		closeCh: make(chan struct{}),
	}
}

// start - метод, который запускает приложение с обновлением по time.Ticker.
func (w *WPoller) Start() {
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

// close - метод, который закрывает канал closeCh.
func (w *WPoller) Close() {
	close(w.closeCh)
}

// handleData обрабатывает сырых данных
func (w *WPoller) decodeData(r io.Reader) (*domain.WeatherData, error) {
	var (
		wData  domain.WeatherData
		m      = make(map[string]*json.RawMessage)
		mSlice = make(map[string][]any)
	)
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
			if err := json.Unmarshal(*value, &mSlice); err != nil {
				return nil, err
			}
			for k, v := range mSlice {
				switch k {
				case "time":
					timeSlice, err := utils.ConvertSliceAnyToSliceString(v)
					if err != nil {
						return nil, err
					}
					wData.Time = timeSlice
				case "temperature_2m":
					tempSlice, err := utils.ConvertSliceAnyToSliceFloat(v)
					if err != nil {
						return nil, err
					}
					wData.Temperature = tempSlice
				default:
					continue
				}
			}
		default:
			continue
		}
	}
	return &wData, nil
}

// GetWeatherByCoordiantes - метод, который возвращает данные погоды, на основе
// вводимых данных координат lat, long.
func (w *WPoller) GetWeatherByCoordiantes(lat, long float64) (*domain.WeatherData, error) {
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
	// обработка данных из io.Reader с чтением и получением структуры domain.WeatherData
	weatherD, err := w.decodeData(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return weatherD, nil
}
