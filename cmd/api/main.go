package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jhonathann10/temperature-system-otel/configs"
	"github.com/jhonathann10/temperature-system-otel/internal/infra"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/viacep"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/weatherapi"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/providerotel"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/webserver"
	"go.opentelemetry.io/otel"
)

const (
	baseURLViaCEP  = "https://viacep.com.br/ws/"
	baseURLWeather = "http://api.weatherapi.com/v1"
	serverPort     = ":8080"
)

func main() {
	tp := providerotel.InitProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("temperature-system")

	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	server := webserver.NewWebServer(serverPort)
	cepHandler := viacep.NewViaCEPClient(baseURLViaCEP)
	weatherHandler := weatherapi.NewWeatherAPIClient(baseURLWeather, config.KeyAPIWeatherApi)
	handler := infra.NewHandler(cepHandler, weatherHandler, tracer)
	server.AddHandler("/temperature", handler.GetTemperature)
	fmt.Println("Starting web server client on port", serverPort)

	server.Start()
}
