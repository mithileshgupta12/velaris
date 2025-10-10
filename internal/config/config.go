package config

import "flag"

type Config struct {
	DB DBFlags
}

func NewConfig() *Config {
	c := &Config{}
	c.DB.Load()

	flag.Parse()

	return c
}
