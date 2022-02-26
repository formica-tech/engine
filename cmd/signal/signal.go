package main

import (
	"database/sql"
	"github.com/hamzali/formica-engine/cmd/signal/rest"
	"github.com/hamzali/formica-engine/impl/entitysignalrepo"
	"github.com/hamzali/formica-engine/internal"
	"github.com/hamzali/formica-engine/usecases"
	"net/http"
)

func main() {
	loggers := internal.InitLogger()

	dsn := internal.ReadConfigWithDefault(
		"postgresql.dsn",
		"host=178.62.208.93 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		loggers.Error().Fatalf("db connection error: %v\n", err)
	}

	su := usecases.NewSignalUseCases(entitysignalrepo.NewSql(db))
	handler := rest.SetupRestServer(su)

	port := internal.ReadConfigWithDefault("port", "8080")
	loggers.Info().Printf("started listening on port %s...\n", port)
	err = http.ListenAndServe(
		":"+port,
		handler,
	)
	if err != nil {
		loggers.Error().Fatalf("server startup error: %v\n", err)
	}
}
