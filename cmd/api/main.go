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

	db, err := db.NewDB(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	log.Println("Connection to database successful")
	defer db.Close()

	r := route.NewRouter()
	r.RegisterRoutes()
	log.Fatal(r.Serve(8000))
}
