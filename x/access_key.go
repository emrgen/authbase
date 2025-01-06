package x

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type AccessKey struct {
	ID    uuid.UUID
	Value string
}

// NewAccessKey creates a new access key
func NewAccessKey() AccessKey {
	token := generateSecureToken(32)
	return AccessKey{
		ID:    uuid.New(),
		Value: token,
	}
}

func IsAccessKey(key string) bool {
	return strings.HasPrefix(key, AccessTokenPrefix)
}

// ParseAccessKey parses an access key
func ParseAccessKey(key string) (*AccessKey, error) {
	token := ParseToken(key)
	if !token.IsAccessToken() {
		return nil, ErrInvalidToken
	}

	accessKey := token.Value[33:]
	accessKeyID := token.Value[1:33]

	id, err := uuidFromStripped(accessKeyID)
	if err != nil {
		return nil, err
	}

	return &AccessKey{
		ID:    id,
		Value: accessKey,
	}, nil
}

// String returns the string representation of the access key
func (a AccessKey) String() string {
	token := fmt.Sprintf("%s%s", uuidStripped(a.ID), a.Value)
	return NewToken(AccessTokenPrefix, token).String()
}

func uuidStripped(id uuid.UUID) string {
	return strings.ReplaceAll(id.String(), "_", "")
}

func uuidFromStripped(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(strings.Join([]string{uuidStr[:8], uuidStr[8:12], uuidStr[12:16], uuidStr[16:20], uuidStr[20:]}, "-"))
}
