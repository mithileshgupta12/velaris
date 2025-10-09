package main

import (
	"log"

	"github.com/mithileshgupta12/velaris/internal/api/route"
)

func main() {
	r := route.NewRouter()
	r.RegisterRoutes()
	log.Fatal(r.Serve(8000))
}
