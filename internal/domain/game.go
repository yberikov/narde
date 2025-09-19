package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type Game struct {
	ID            uuid.UUID
	WhitePlayerID uuid.UUID
	BlackPlayerID uuid.UUID

	Status    StatusType
	WinnerID  uuid.UUID
	CreatedAt time.Time
	EndedAt   time.Time
}

type StatusType string

var (
	InProgress StatusType = "in_progress"
	Finished   StatusType = "finished"
)

type Move struct {
	From int `json:"from"`
	To   int `json:"to"`
}
