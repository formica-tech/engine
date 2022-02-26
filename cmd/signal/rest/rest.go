package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hamzali/formica-engine/usecases"
	"net/http"
	"time"
)

func SetupRestServer(signalUseCases usecases.SignalUseCases) http.Handler {
	baseRouter := chi.NewRouter()
	baseRouter.Use(middleware.RequestID)
	baseRouter.Use(middleware.RealIP)
	baseRouter.Use(middleware.Logger)
	baseRouter.Use(middleware.Recoverer)
	baseRouter.Use(middleware.Timeout(60 * time.Second))

	h := Handler{
		signalUseCases: signalUseCases,
	}
	baseRouter.Post("/signal", h.BatchInsertSdkSignal)

	return baseRouter
}
