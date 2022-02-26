package usecases

import (
	"context"
	"github.com/hamzali/formica-engine/domain"
)

type SignalUseCases interface {
	BatchSave(ctx context.Context, signals []domain.EntitySignal) error
}

type signalUseCases struct {
	entitySignalRepository domain.EntitySignalRepository
}

func NewSignalUseCases(
	entitySignalRepository domain.EntitySignalRepository,

) *signalUseCases {
	return &signalUseCases{
		entitySignalRepository: entitySignalRepository,
	}
}

func (app *signalUseCases) BatchSave(_ context.Context, signals []domain.EntitySignal) error {
	return app.entitySignalRepository.BatchSave(signals)
}
