package domain

type WeatherData struct {
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Time        []string  `json:"time"`
	Temperature []float64 `json:"temperature_2m"`
}
