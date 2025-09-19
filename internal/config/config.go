package config

import "time"

type (
	Config struct {
		Main       Main       `envPrefix:"MAIN_" env:"inline"`
		Postgresdb Postgresdb `envPrefix:"POSTGRES_" env:"inline"`
	}

	Main struct {
		LogLevel          string `env:"LOG_LEVEL"`
		LogForcePlainText bool   `env:"LOG_FORCE_PLAIN_TEXT"`
	}

	Postgresdb struct {
		Main PostgresCreds `envPrefix:"MAIN_"`
	}

	PostgresCreds struct {
		DBName          string        `env:"DB_NAME,required"`
		User            string        `env:"USER,required"`
		Password        string        `env:"PASSWORD,required"`
		Host            string        `env:"HOST,required"`
		Port            int           `env:"PORT,required"`
		SSLMode         string        `env:"SSL_MODE,required"`
		MaxOpenConns    int           `env:"MAX_OPEN_CONNS,required"`
		MaxIdleConns    int           `env:"MAX_IDLE_CONNS,required"`
		ConnMaxIdleTime time.Duration `env:"CONN_MAX_IDLE_TIME,required"` // Time in seconds // 200s
		ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME,required"`  // Time in seconds // 500s
		MigratePath     string        `env:"MIGRATE_PATH"`
	}
)
