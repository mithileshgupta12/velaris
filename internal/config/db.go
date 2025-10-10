package config

import "flag"

type DBFlags struct {
	Host     string
	User     string
	Password string
	Name     string
	PORT     int
	SSLMode  string
}

func (dbf *DBFlags) Load() *DBFlags {
	flag.StringVar(&dbf.Host, "db-host", "localhost", "Database host")
	flag.StringVar(&dbf.Name, "db-name", "velaris", "Database name")
	flag.StringVar(&dbf.User, "db-username", "postgres", "Database username")
	flag.StringVar(&dbf.Password, "db-password", "password", "Database password")
	flag.IntVar(&dbf.PORT, "db-port", 5432, "Database port")
	flag.StringVar(&dbf.SSLMode, "db-sslmode", "disable", "Database SSL mode")

	return dbf
}
