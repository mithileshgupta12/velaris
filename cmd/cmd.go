package cmd

import (
	"log/slog"

	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/helper"
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"github.com/mithileshgupta12/velaris/internal/route"
)

func Execute() {
	cfg := config.NewConfig()

	repositories, err := db.NewDB(&cfg.DB)
	if err != nil {
		helper.LogFatal("failed to connect to database", "err", err)
	}

	slog.Info("Connection to database successful")

	cache, err := cache.NewRedisClient()
	if err != nil {
		helper.LogFatal("failed to connect to cache", "err", err)
	}

	stores := cache.InitStores()

	slog.Info("Connection to cache successful")
	defer cache.Close()

	middlewares := middleware.NewMiddlewares(repositories, stores.SessionStore)

	r := route.NewRouter(cfg.App.FrontendUrl)
	r.RegisterRoutes(repositories, stores, middlewares)
	if err := r.Serve(cfg.App.Port); err != nil {
		helper.LogFatal("failed to start server", "err", err)
	}
}
