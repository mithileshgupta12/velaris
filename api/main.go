package main

import (
	"log"

	"github.com/mithileshgupta12/velaris/routes"
)

func main() {
	router := routes.NewRouter().Init()

	if err := router.Serve(":8000"); err != nil {
		log.Fatal(err)
	}
}
