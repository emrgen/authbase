package x

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type AccessKey struct {
	ID    string
	Value string
}

// NewAccessKey creates a new access key
func NewAccessKey() AccessKey {
	return AccessKey{
		ID:    uuid.New().String(),
		Value: randomString(32),
	}
}

func ParseAccessKey(key string) (*AccessKey, error) {
	token := ParseToken(key)
	if !token.IsAccessToken() {
		return nil, fmt.Errorf("invalid access key")
	}

	decoded, err := base64DecodeStripped(token.Value)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(decoded, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid access key")
	}

	return &AccessKey{
		ID:    parts[0],
		Value: parts[1],
	}, nil
}

func (a AccessKey) String() string {
	token := base64EncodeStripped(fmt.Sprintf("%s-%s", a.ID, a.Value))
	return NewToken(AccessTokenPrefix, token).String()
}
