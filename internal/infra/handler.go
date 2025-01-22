package infra

import (
	"encoding/json"
	"net/http"

	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/viacep"
	"github.com/jhonathann10/temperature-system-otel/internal/infra/client/weatherapi"
	"github.com/jhonathann10/temperature-system-otel/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	viaCep  viacep.ViaCepInterface
	weather weatherapi.WeatherAPIInterface
	tracer  trace.Tracer
}

func NewHandler(viaCep viacep.ViaCepInterface, weather weatherapi.WeatherAPIInterface, tracer trace.Tracer) *Handler {
	return &Handler{
		viaCep:  viaCep,
		weather: weather,
		tracer:  tracer,
	}
}

func (h *Handler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	ctx, span := h.tracer.Start(ctx, "span-temperature-system")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	cep := r.URL.Query().Get("cep")

	if isCepInvalid(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPResponseStatusCodeKey.Int(http.StatusUnprocessableEntity),
			semconv.URLFullKey.String(r.URL.String()),
		)
		return
	}

	temperatureUseCase := usecase.NewTemperatureUseCase(h.viaCep, h.weather)
	localidade, err := temperatureUseCase.Execute(ctx, cep)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Message})
		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPResponseStatusCodeKey.Int(err.StatusCode),
			semconv.URLFullKey.String(r.URL.String()),
		)
		return
	}
	span.SetAttributes(
		semconv.HTTPRequestMethodKey.String(r.Method),
		semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK),
		semconv.URLFullKey.String(r.URL.String()),
	)

	errJson := json.NewEncoder(w).Encode(localidade)
	if errJson != nil {
		http.Error(w, errJson.Error(), http.StatusInternalServerError)
		return
	}
}

func isCepInvalid(cep string) bool {
	return len(cep) != 8
}
