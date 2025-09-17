package jwt

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type Generator struct {
	config *Config
	parser *Parser
}

var (
	ErrInvalidTokenType = errors.New("invalid token type")
)

func MustGenerator() *Generator {
	generator := &Generator{parser: MustParser()}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot parse generator env")
	}
	generator.config = &cfg

	return generator
}

func (r *Generator) Tokens(ctx context.Context, userID uuid.UUID) (string, string, error) {
	logger := zerolog.Ctx(ctx)

	accessToken, err := r.generateTokenByType(Access, userID)
	if err != nil {
		logger.Error().Err(err).Send()
		return "", "", err
	}

	refreshToken, err := r.generateTokenByType(Refresh, userID)
	if err != nil {
		logger.Error().Err(err).Send()
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (r *Generator) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	logger := zerolog.Ctx(ctx)

	user, err := r.parser.ParseToken(ctx, refreshToken)
	if err != nil {
		logger.Error().Err(err).Send()
		return "", "", err
	}
	return r.Tokens(ctx, user.ID)
}

func (r *Generator) generateTokenByType(tokenType Type, userID uuid.UUID) (string, error) {
	now := time.Now().UTC()

	var tokenConfig TokenConfig

	switch tokenType {
	case Access:
		tokenConfig = r.config.Access
	case Refresh:
		tokenConfig = r.config.Refresh
	default:
		return "", ErrInvalidTokenType
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    r.config.Issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenConfig.Ttl)),
		},
	})

	return t.SignedString([]byte(tokenConfig.Key))
}
