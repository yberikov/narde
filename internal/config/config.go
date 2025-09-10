package config

type (
	Config struct {
		Main Main `envPrefix:"MAIN_" env:"inline"`
	}

	Main struct {
		LogLevel          string `env:"LOG_LEVEL"`
		LogForcePlainText bool   `env:"LOG_FORCE_PLAIN_TEXT"`
	}
)
