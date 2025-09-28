package cache

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"narde/internal/cache/model"
	"narde/internal/domain"
)

type (
	gameStateCache struct {
		rdb *redis.Client
	}
)

func NewGameStateCache(rdb *redis.Client) *gameStateCache {
	return &gameStateCache{
		rdb: rdb,
	}
}

func (c gameStateCache) getKeyForStates() string {
	return "game_states"
}

func (c gameStateCache) GetGameState(ctx context.Context, gameID uuid.UUID) (*domain.GameState, error) {
	logger := zerolog.Ctx(ctx).With().Str("gameID", gameID.String()).Logger()

	logger.Debug().Msg("Fetching last state from cache...")
	result, err := c.rdb.HGet(ctx, c.getKeyForStates(), gameID.String()).Result()
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}

	var state model.State
	if err = json.Unmarshal([]byte(result), &state); err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}

	return state.ToDomain(), nil
}

func (c gameStateCache) SaveGameState(ctx context.Context, item *domain.GameState) error {
	logger := zerolog.Ctx(ctx)

	stateData, err := json.Marshal(model.NewStateFromDomain(item))
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	}

	if err := c.rdb.HSet(ctx, c.getKeyForStates(), item.GameID, stateData).Err(); err != nil {
		logger.Error().Err(err).Send()
		return err
	}
	return nil
}
