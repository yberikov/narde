package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
