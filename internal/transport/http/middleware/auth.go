package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"narde/pkg/jwt"
	"strings"
)

type (
	Auth struct {
		parser *jwt.Parser
	}
)

func NewAuth(parser *jwt.Parser) *Auth {
	return &Auth{parser: parser}
}

func (m Auth) Handle(c *fiber.Ctx) error {
	ctx := c.Context()
	logger := zerolog.Ctx(ctx)
	authorization := c.Get(fiber.HeaderAuthorization)

	if len(authorization) <= 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid jwt token")
	}

	jwtToken, err := m.parseAndTrimHeader(ctx, authorization)
	if err != nil {
		return err
	}

	sessionUser, err := m.parser.ParseToken(ctx, jwtToken)
	if err != nil {
		logger.Error().Err(err).Send()
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid jwt token")
	}

	c.Locals(jwt.UserLocalFiberKey, sessionUser.ID)

	return c.Next()
}

func (m Auth) parseAndTrimHeader(ctx context.Context, authorization string) (string, error) {
	logger := zerolog.Ctx(ctx)

	logger.Debug().Msg("Check jwt token passed...")

	l := len(jwt.AuthTypeBearer)
	if strings.EqualFold(authorization[:l], jwt.AuthTypeBearer) {
		return strings.TrimSpace(authorization[l:]), nil
	}

	logger.Warn().Msg("Token is not pass in Authorization header")

	return "", fiber.NewError(fiber.StatusUnauthorized, "Token is not pass")
}
