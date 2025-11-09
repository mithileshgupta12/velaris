package main

import (
	"fmt"

	"github.com/mithileshgupta12/velaris/internal/api/middleware"
	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/cache"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
	"github.com/mithileshgupta12/velaris/internal/pkg/logger"
)

func main() {
	lgr := logger.NewLogger()

	cfg := config.NewConfig()

	repositories, err := db.NewDB(&cfg.DB)
	if err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to connect to database: %v", err), nil)
	}

	lgr.Log(logger.INFO, "Connection to database successful", nil)

	cache, err := cache.NewRedisClient()
	if err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to connect to cache: %v", err), nil)
	}

	stores := cache.InitStores()

	lgr.Log(logger.INFO, "Connection to cache successful", nil)
	defer cache.Close()

	middlewares := middleware.NewMiddlewares(lgr, repositories, stores.SessionStore)

	r := route.NewRouter(lgr, cfg.App.FrontendUrl)
	r.RegisterRoutes(repositories, stores, middlewares)
	if err := r.Serve(cfg.App.Port); err != nil {
		lgr.Log(logger.FATAL, fmt.Sprintf("failed to start server: %v", err), nil)
	}
}
