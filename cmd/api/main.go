package main

import (
	"fmt"
	"log/slog"

	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/helper"
)

func main() {
	cfg := config.NewConfig()

	repositories, err := db.NewDB(&cfg.DB)
	if err != nil {
		helper.LogFatal(fmt.Sprintf("failed to connect to database: %v", err))
	}

	slog.Info("Connection to database successful")

	cache, err := cache.NewRedisClient()
	if err != nil {
		helper.LogFatal(fmt.Sprintf("failed to connect to cache: %v", err))
	}

	stores := cache.InitStores()

	slog.Info("Connection to cache successful")
	defer cache.Close()

	middlewares := middleware.NewMiddlewares(repositories, stores.SessionStore)

	r := route.NewRouter(cfg.App.FrontendUrl)
	r.RegisterRoutes(repositories, stores, middlewares)
	if err := r.Serve(cfg.App.Port); err != nil {
		helper.LogFatal(fmt.Sprintf("failed to start server: %v", err))
	}
}
