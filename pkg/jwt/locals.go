package jwt

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

var errorInvalidTokenFmt = fiber.NewError(fiber.StatusBadRequest, "INVALID_TOKEN_FORMAT")

const (
	UserLocalFiberKey = "userID"
	AuthTypeBearer    = "Bearer"
)

func GetUserIdFromCtx(c *fiber.Ctx) (uuid.UUID, error) {
	currentUserID, ok := c.Locals(UserLocalFiberKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errorInvalidTokenFmt
	}

	return currentUserID, nil
}
