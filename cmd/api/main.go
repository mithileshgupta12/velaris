package main

import (
	"context"
	"fmt"

	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func main() {
	lgr := logger.NewLogger(logger.FormatJSON)

	cfg := config.NewConfig()

	database, err := db.NewDB(&cfg.DB)
	if err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to connect to database: %v", err), nil)
	}

	if err := database.Ping(context.Background()); err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to ping database: %v", err), nil)
	}

	lgr.Log(logger.INFO, "Connection to database successful", nil)
	defer database.Close()

	cache, err := cache.NewRedisClient()
	if err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to connect to cache: %v", err), nil)
	}

	stores := cache.InitStores()

	lgr.Log(logger.INFO, "Connection to cache successful", nil)
	defer cache.Close()

	r := route.NewRouter(lgr, database.Queries, stores)
	r.RegisterRoutes()
	if err := r.Serve(8000); err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to start server: %v", err), nil)
	}
}
