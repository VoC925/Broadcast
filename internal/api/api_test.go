package api

import (
	"os"
	"testing"

	"github.com/VoC925/WeatherChecker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockSender struct{}

func (m *MockSender) Send(data []byte) error {
	return nil
}

const (
	apiWeatherResponseFile = "responseWeatherApi.txt"
	apiIPResponseFile      = "responseIPApi.txt"
)

func getWeatherData() *domain.WeatherData {
	return &domain.WeatherData{
		Latitude:  55.75,
		Longitude: 37.625,
	}
}

func TestDecodeWeatherDataMethod(t *testing.T) {
	// мок-заглушка интерфейса Sender
	mockSender := new(MockSender)
	// экземпляр структуры сервиса погоды
	newWPoller := NewWPoller("62.181.51.102", mockSender)
	// Reader на основе файла ответа от API
	// Проверка существования файла
	_, err := os.Stat(apiWeatherResponseFile)
	require.NoError(t, err)
	// открытие файла для чтения
	responseReader, err := os.OpenFile(apiWeatherResponseFile, os.O_RDONLY, 0644)
	require.NoError(t, err)
	data, err := newWPoller.decodeWeatherData(responseReader)
	require.NoError(t, err)
	expected := getWeatherData()
	// проверка равенства значений
	assert.Equal(t, expected.Latitude, data.Latitude)
	assert.Equal(t, expected.Longitude, data.Longitude)
	// проверка равенства типов
	assert.IsType(t, "string", data.Time[0])
	assert.IsType(t, 1.0, data.Temperature[0])
	// проверка длины слайсов
	assert.Equal(t, len(data.Time), len(data.Temperature))
}

type testIPCoordinates struct {
	addressIP string
	coord     *coordinates
}

func getTestindDataIP() []testIPCoordinates {
	return []testIPCoordinates{
		{
			addressIP: "62.181.51.102",
			coord: &coordinates{
				Lat: 55.6784,
				Lon: 37.2652,
			},
		},
	}
}

func TestDecodeIPDataMethod(t *testing.T) {
	// мок-заглушка интерфейса Sender
	mockSender := new(MockSender)
	// тестовые данные
	testData := getTestindDataIP()
	// Проверка существования файла
	_, err := os.Stat(apiIPResponseFile)
	require.NoError(t, err)

	for _, data := range testData {
		// экземпляр структуры сервиса погоды
		newWPoller := NewWPoller(data.addressIP, mockSender)
		// метод
		coord, err := newWPoller.getCoordinatesByIP()
		require.NoError(t, err)
		// проверка равенства значений
		assert.Equal(t, data.coord.Lat, coord.Lat)
		assert.Equal(t, data.coord.Lon, coord.Lon)
	}
}
