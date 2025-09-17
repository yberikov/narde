package service

import (
	"context"
	"github.com/gofrs/uuid"
)

type (
	TokenGenerator interface {
		Tokens(ctx context.Context, userID uuid.UUID) (string, string, error)
	}
)
