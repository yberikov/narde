package service

import (
	"context"
	"github.com/gofrs/uuid"
	"narde/internal/domain"
)

type (
	AuthRepository interface {
		CreateUser(ctx context.Context, user *domain.User) error
		GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
		GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	}

	GameRepository interface {
	}

	GameStateCacher interface {
		GetGameStateByID(ctx context.Context, gameID uuid.UUID) (*domain.GameState, error)
		SaveGameState(ctx context.Context, state *domain.GameState) error
	}

	TokenGenerator interface {
		Tokens(ctx context.Context, userID uuid.UUID) (string, string, error)
	}
)
