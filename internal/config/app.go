package config

import "flag"

type AppFlags struct {
	Port        int
	FrontendUrl string
}

func (af *AppFlags) Load() *AppFlags {
	flag.IntVar(&af.Port, "app-port", 8000, "Port number for the application server")
	flag.StringVar(&af.FrontendUrl, "frontend-url", "http://localhost:8000", "Frontend URL for CORS and redirects")

	return af
}
