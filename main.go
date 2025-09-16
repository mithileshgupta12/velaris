package main

import (
	"log"

	"example.com/velaris/internal/api/route"
)

func main() {
	r := route.NewRouter()

	if err := r.Serve(8000); err != nil {
		log.Fatal(err.Error())
	}
}
