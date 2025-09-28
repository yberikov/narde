package model

import (
	"github.com/gofrs/uuid"
	"narde/internal/domain"
	"time"
)

type (
	State struct {
		GameID        uuid.UUID `json:"gameId"`
		WhitePlayerID uuid.UUID `json:"whitePlayerId"`
		BlackPlayerID uuid.UUID `json:"blackPlayerId"`

		Board [26]int `json:"board"`

		Turn       string    `json:"turn"`
		Dice       [2]int    `json:"dice"`
		LastMoveAt time.Time `json:"lastMoveAt"`
	}
)

func NewStateFromDomain(item *domain.GameState) State {
	state := State{
		GameID:        item.GameID,
		WhitePlayerID: item.WhitePlayerID,
		BlackPlayerID: item.BlackPlayerID,
		Board:         item.Board,
		Turn:          string(item.Turn),
		Dice:          item.Dice,
		LastMoveAt:    item.LastMoveAt,
	}

	return state
}

func (c State) ToDomain() *domain.GameState {
	state := &domain.GameState{
		GameID:        c.GameID,
		WhitePlayerID: c.WhitePlayerID,
		BlackPlayerID: c.BlackPlayerID,
		Board:         c.Board,
		Turn:          domain.TurnType(c.Turn),
		Dice:          c.Dice,
		LastMoveAt:    c.LastMoveAt,
	}

	return state
}
