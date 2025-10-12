package main

import (
	"context"
	"log"

	"github.com/mithileshgupta12/velaris/internal/api/route"
	"github.com/mithileshgupta12/velaris/internal/config"
	"github.com/mithileshgupta12/velaris/internal/db"
)

func main() {
	cfg := config.NewConfig()

	database, err := db.NewDB(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := database.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	log.Println("Connection to database successful")
	defer database.Close()

	r := route.NewRouter(database)
	r.RegisterRoutes()
	log.Fatal(r.Serve(8000))
}
