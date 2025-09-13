package main

import (
	"fmt"
	"log"

	"github.com/mithileshgupta12/velaris/routes"
)

func main() {
	router := routes.NewRouter().Init()

	fmt.Println("Hello, World!")

	if err := router.Serve(":8000"); err != nil {
		log.Fatal(err)
	}
}
