package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jhonathann10/temperature-system-otel/internal/httperror"
)

type LocalidadeCEP struct {
	Localidade string `json:"localidade"`
}

type ViaCEPClient struct {
	BaseURL string
}

func NewViaCEPClient(baseURL string) *ViaCEPClient {
	return &ViaCEPClient{
		BaseURL: baseURL,
	}
}

func (v *ViaCEPClient) GetAddressByCEP(ctx context.Context, cep string) (*LocalidadeCEP, *httperror.HttpError) {
	localidade := &LocalidadeCEP{}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/json/", v.BaseURL, cep), nil)
	if err != nil {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(localidade)
	if err != nil {
		return nil, &httperror.HttpError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return localidade, nil
}
