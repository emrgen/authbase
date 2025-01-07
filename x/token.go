package x

import "fmt"

const (
	AccessTokenPrefix  = "aba"
	RefreshTokenPrefix = "abr"
	AccessKeyPrefix    = "abk"
)

var ErrInvalidToken = fmt.Errorf("not a valid access token")

type Token struct {
	Kind  string
	Value string
}

func ParseToken(token string) *Token {
	return &Token{
		Kind:  token[:3],
		Value: token[3:],
	}
}

func NewToken(kind, value string) *Token {
	return &Token{
		Kind:  kind,
		Value: value,
	}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s_%s", t.Kind, t.Value)
}

func (t *Token) IsAccessToken() bool {
	return t.Kind == AccessTokenPrefix
}

func (t *Token) IsRefreshToken() bool {
	return t.Kind == RefreshTokenPrefix
}

func (t *Token) IsAccessKey() bool {
	return t.Kind == AccessKeyPrefix
}
