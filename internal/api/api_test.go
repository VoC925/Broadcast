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
	apiResponseFile = "responseApi.txt"
)

func getWeatherData() *domain.WeatherData {
	return &domain.WeatherData{
		Latitude:  55.75,
		Longitude: 37.625,
	}
}

func TestDecodeDataMethod(t *testing.T) {
	// мок-заглушка интерфейса Sender
	mockSender := new(MockSender)
	// экземпляр структуры сервиса погоды
	newWPoller := NewWPoller(mockSender)
	// Reader на основе файла ответа от API
	// Проверка существования файла
	_, err := os.Stat(apiResponseFile)
	require.NoError(t, err)
	// открытие файла для чтения
	responseReader, err := os.OpenFile(apiResponseFile, os.O_RDONLY, 0644)
	require.NoError(t, err)
	data, err := newWPoller.decodeData(responseReader)
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
