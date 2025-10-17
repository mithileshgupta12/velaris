package main

import (
	"context"
	"log"

	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func main() {
	lgr := logger.NewLogger()

	cfg := config.NewConfig()

	database, err := db.NewDB(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := database.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	lgr.Log(logger.FormatJSON, logger.INFO, "Connection to database successful")
	defer database.Close()

	r := route.NewRouter(lgr, database)
	r.RegisterRoutes()
	log.Fatal(r.Serve(8000))
}
