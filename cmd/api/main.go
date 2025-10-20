package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func main() {
	lgr := logger.NewLogger(logger.FormatJSON)

	cfg := config.NewConfig()

	database, err := db.NewDB(&cfg.DB)
	if err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("Failed to connect to database: %v", err), nil)
	}

	if err := database.Ping(context.Background()); err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("Failed to ping database: %v", err), nil)
	}

	lgr.Log(logger.INFO, "Connection to database successful", nil)
	defer database.Close()

	r := route.NewRouter(lgr, database)
	r.RegisterRoutes()
	log.Fatal(r.Serve(8000))
}
