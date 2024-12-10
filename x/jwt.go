package x

import "time"

type JWTToken struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
	IssuedAt     time.Time
}

func GenerateJWTToken(userID, organizationID string) JWTToken {
	return JWTToken{}
}
