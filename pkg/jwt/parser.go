package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const unexpectedSigningMethod = "unexpected signing method: %s"

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type Parser struct {
	config *Config
}

func MustParser() *Parser {
	parser := &Parser{}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot parse generator env")
	}
	parser.config = &cfg

	return parser
}

func (p *Parser) ParseToken(ctx context.Context, jwtToken string) (*SessionUser, error) {
	logger := zerolog.Ctx(ctx)

	logger.Debug().Msg("Parse and validate jwt token...")

	var claims Claims
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&claims,
		p.withSignMethod(ctx),
		jwt.WithLeeway(p.config.Leeway),
		jwt.WithIssuer(p.config.Issuer),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		logger.Error().Err(err).Send()

		return nil, ErrInvalidToken
	}

	if !token.Valid {
		logger.Error().Err(ErrInvalidToken).Send()

		return nil, ErrInvalidToken
	}

	userID, err := uuid.FromString(claims.Subject)
	if err != nil {
		logger.Error().Err(err).Send()

		return nil, ErrInvalidToken
	}
	return &SessionUser{userID}, nil
}

func (p *Parser) withSignMethod(ctx context.Context) func(token *jwt.Token) (any, error) {
	logger := zerolog.Ctx(ctx)

	return func(token *jwt.Token) (any, error) {
		claims := token.Claims.(*Claims)

		var key string
		switch claims.Type {
		case Access:
			key = p.config.Access.Key
		case Refresh:
			key = p.config.Refresh.Key
		default:
			return "", ErrInvalidTokenType
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error().Msg("Unexpected signing method")

			return nil, fmt.Errorf(unexpectedSigningMethod, token.Header["alg"])
		}

		return []byte(key), nil
	}
}
