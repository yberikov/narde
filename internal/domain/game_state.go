package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type GameState struct {
	GameID        uuid.UUID
	WhitePlayerID uuid.UUID
	BlackPlayerID uuid.UUID

	Board      [24]int
	Turn       TurnType
	Dice       [2]int
	LastMoveAt time.Time
}

type TurnType string

var (
	WhiteTurnType TurnType = "white"
	BlackTurnType TurnType = "black"
)
