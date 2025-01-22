package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	TempCelsius float64 `json:"temp_c"`
}

type WeatherAPIClient struct {
	BaseURL string
	Token   string
}

func NewWeatherAPIClient(baseURL string, token string) *WeatherAPIClient {
	return &WeatherAPIClient{
		BaseURL: baseURL,
		Token:   token,
	}
}

func (w *WeatherAPIClient) GetWeatherByCity(ctx context.Context, city string) (*Weather, error) {
	weather := &Weather{}
	encodeCity := url.QueryEscape(city)
	uri := fmt.Sprintf("%s/current.json?key=%s&q=%s", w.BaseURL, w.Token, encodeCity)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(weather)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
