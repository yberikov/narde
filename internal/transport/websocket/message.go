package websocket

import (
	"encoding/json"
	"narde/internal/domain"
)

type (
	Message struct {
		Type    MessageType     `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	MakeMovePayload struct {
		Moves []domain.Move `json:"moves"`
	}
)

type MessageType string

var (
	MessageTypeFindMatch MessageType = "find_match"
	MessageTypeGameStart MessageType = "game_start"
	MessageTypeMakeMove  MessageType = "make_move"
	MessageTypeGameState MessageType = "game_state"
)
