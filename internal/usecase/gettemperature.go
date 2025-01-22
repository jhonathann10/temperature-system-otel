package usecase

import (
	"context"
	"net/http"

	"github.com/jhonathann10/temperature-system-otel/internal/entity"
	"github.com/jhonathann10/temperature-system-otel/internal/httperror"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/viacep"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/weatherapi"
)

type WeatherDTO struct {
	Localidade string  `json:"localidade"`
	Celsius    float64 `json:"temp_c,omitempty"`
	Fahrenheit float64 `json:"temp_f,omitempty"`
	Kelvin     float64 `json:"temp_k,omitempty"`
}

type GetTemperatureUseCase struct {
	viaCepInterface  viacep.ViaCepInterface
	weatherInterface weatherapi.WeatherAPIInterface
}

func NewTemperatureUseCase(viaCepInterface viacep.ViaCepInterface, weatherInterface weatherapi.WeatherAPIInterface) *GetTemperatureUseCase {
	return &GetTemperatureUseCase{
		viaCepInterface:  viaCepInterface,
		weatherInterface: weatherInterface,
	}
}

func (g *GetTemperatureUseCase) Execute(ctx context.Context, cep string) (*WeatherDTO, *httperror.HttpError) {
	localidade, err := g.viaCepInterface.GetAddressByCEP(ctx, cep)
	if err != nil {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
		}
	}
	if localidade.Localidade == "" {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusNotFound,
			Message:    "can not find zipcode",
		}
	}

	weather, errWeatherByCity := g.weatherInterface.GetWeatherByCity(ctx, localidade.Localidade)
	if errWeatherByCity != nil {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    errWeatherByCity.Error(),
		}
	}

	newWeather, errWeather := entity.NewWeather(localidade.Localidade, weather.Current.TempCelsius)
	if errWeather != nil {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    errWeather.Error(),
		}
	}

	return toWeatherDTO(newWeather), nil
}

func toWeatherDTO(newWeather *entity.Weather) *WeatherDTO {
	return &WeatherDTO{
		Localidade: newWeather.Localidade,
		Celsius:    newWeather.Celsius,
		Fahrenheit: newWeather.Fahrenheit,
		Kelvin:     newWeather.Kelvin,
	}
}
