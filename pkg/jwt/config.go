package jwt

import "time"

type Config struct {
	Issuer string `env:"TOKEN_ISSUER,required"`

	Leeway time.Duration `env:"LEEWAY,required"`

	Access  TokenConfig `envPrefix:"ACCESS_TOKEN_"`
	Refresh TokenConfig `envPrefix:"REFRESH_TOKEN_"`
}

type TokenConfig struct {
	Key string        `env:"KEY,required"`
	Ttl time.Duration `env:"TTL,required"`
}
