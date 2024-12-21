package x

import (
	"context"
	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrUserNotFoundInContext
	}

	return userID, nil
}
