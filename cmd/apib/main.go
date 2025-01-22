package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jhonathann10/temperature-system-otel/internal/infra"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/providerotel"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/webserver"
	"go.opentelemetry.io/otel"
)

const (
	serverPort = ":8081"
)

func main() {
	tp := providerotel.InitProvider()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("temperature-system-otel-b")

	server := webserver.NewWebServer(serverPort)

	handler := infra.NewHandlerB(tracer)
	server.AddHandler("/temperature", handler.Handle)
	fmt.Println("Starting web server on port", serverPort)

	server.Start()
}
