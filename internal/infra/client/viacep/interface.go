package viacep

import (
	"context"

	"github.com/jhonathann10/temperature-system-otel/internal/httperror"
)

type ViaCepInterface interface {
	GetAddressByCEP(ctx context.Context, cep string) (*LocalidadeCEP, *httperror.HttpError)
}
