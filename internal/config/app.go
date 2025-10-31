package config

type AppFlags struct {
	Port        int
	FrontendUrl string
}

func (af *AppFlags) Load() *AppFlags {
	return af
}
