package config

import "flag"

type Config struct {
	DB  DBFlags
	App AppFlags
}

func NewConfig() *Config {
	c := &Config{}
	c.DB.Load()
	c.App.Load()

	flag.Parse()

	return c
}
