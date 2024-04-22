package api

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/VoC925/WeatherChecker/internal/domain"
	"github.com/VoC925/WeatherChecker/pkg/client"
	"github.com/VoC925/WeatherChecker/pkg/utils"
	"github.com/sirupsen/logrus"
)

const (
	// Api для получения данных погоды
	endpointBroadcast = "https://api.open-meteo.com/v1/forecast" // `?latitude=55.7522&longitude=37.6156&hourly=temperature_2m`
	endpointIP        = "http://ip-api.com"                      // `/json/62.181.51.102`
)

var (
	// временнной промежуток между которым происходит запрос к API
	IntervalReload = time.Second * 1
)

// WPoller структура приложения, получающая температуру и отправляющая sender.
type WPoller struct {
	addressIP string        // IP адрес пользователя
	sender    client.Sender // интерфейс, куда отправлять данные погоды
	closeCh   chan struct{} // канал закрывающий приложение
}

// NewWPoller - конструктор структуры WPoller.
func NewWPoller(addr string, sender client.Sender) *WPoller {
	return &WPoller{
		addressIP: addr,
		sender:    sender,
		closeCh:   make(chan struct{}),
	}
}

// Start() - метод, который запускает приложение с обновлением по time.Ticker
// входные параметры - координаты пользователя: Longitude, Latitude
func (w *WPoller) Start() {
	ticker := time.NewTicker(IntervalReload)
	defer ticker.Stop()

	fmt.Println("------------Weather checking start------------")

	logrus.WithFields(logrus.Fields{
		"IP": w.addressIP,
	}).Info("App started")

	go func() {

		// Получение координат по IP адресу
		coordinates, err := w.getCoordinatesByIP()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"api": "IP",
			}).Errorf(err.Error())
		}

		for {
			// Получение сырых данных
			weatherData, err := w.GetWeatherByCoordiantes(coordinates.Lat, coordinates.Lon)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"api": "Weather",
				}).Errorf(err.Error())

			}
		loop:
			for {
				if !weatherData.GetTemperatureForNowMoment() {
					break loop
				}
				// отправка погоды sender
				if err := w.sender.Send([]byte(weatherData.String())); err != nil {
					logrus.Errorf(err.Error())
					continue
				}
				<-ticker.C
			}
		}
	}()

	<-w.closeCh

	logrus.Infof("App finished")
	fmt.Println("------------Weather checking finished------------")
}

// close - метод, который закрывает канал closeCh.
func (w *WPoller) Close() {
	close(w.closeCh)
}

// decodeWeatherData - кастомный десериализатор
func (w *WPoller) decodeWeatherData(r io.Reader) (*domain.WeatherData, error) {
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
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpointBroadcast, lat, long)
	// запрос к API
	resp, err := utils.GetRequest(uri)
	if err != nil {
		return nil, err
	}
	// обработка данных из io.Reader с чтением и получением структуры domain.WeatherData
	weatherD, err := w.decodeWeatherData(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return weatherD, nil
}

// структура для хранения координат Longitude, Latitude
type coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// getCoordinatesByIP метод, который возвращает координаты по IP адресу
func (w *WPoller) getCoordinatesByIP() (*coordinates, error) {
	uri := fmt.Sprintf("%s/json/%s", endpointIP, w.addressIP)
	// запрос к API
	resp, err := utils.GetRequest(uri)
	if err != nil {
		return nil, err
	}

	coord := new(coordinates)

	if err := json.NewDecoder(resp.Body).Decode(coord); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return coord, nil
}
