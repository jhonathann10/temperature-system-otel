package weatherapi

import "context"

type WeatherAPIInterface interface {
	GetWeatherByCity(ctx context.Context, city string) (*Weather, error)
}
