package infra

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

type HandlerB struct {
	tracer trace.Tracer
}

func NewHandlerB(tracer trace.Tracer) *HandlerB {
	return &HandlerB{
		tracer: tracer,
	}
}

func (h *HandlerB) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "span-temperature-system-otel-b")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	cep := r.URL.Query().Get("cep")

	time.Sleep(1 * time.Second)
	resp, err := makeRequest(ctx, "http://localhost:8080/temperature?cep="+cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "internal server error"})
		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPResponseStatusCodeKey.Int(http.StatusUnprocessableEntity),
			semconv.URLFullKey.String(r.URL.String()),
		)
		return
	}
	defer resp.Body.Close()

	span.SetAttributes(
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK),
		semconv.URLFullKey.String(r.URL.String()),
	)

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}

}

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	return client.Do(req)
}
