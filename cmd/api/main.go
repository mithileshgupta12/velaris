package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
)

func main() {
	cfg := config.NewConfig()

	repositories, err := db.NewDB(&cfg.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to database: %v", err))
		os.Exit(1)
	}

	slog.Info("Connection to database successful")

	cache, err := cache.NewRedisClient()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to cache: %v", err))
		os.Exit(1)
	}

	stores := cache.InitStores()

	slog.Info("Connection to cache successful")
	defer cache.Close()

	middlewares := middleware.NewMiddlewares(repositories, stores.SessionStore)

	r := route.NewRouter(cfg.App.FrontendUrl)
	r.RegisterRoutes(repositories, stores, middlewares)
	if err := r.Serve(cfg.App.Port); err != nil {
		slog.Info(fmt.Sprintf("failed to start server: %v", err))
		os.Exit(1)
	}
}
