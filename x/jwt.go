package x

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = ""

func init() {
	secretKey = os.Getenv("JWT_SECRET")
	if secretKey == "" {
		panic("JWT_SECRET is not set")
	}
}

type Claims struct {
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	OrganizationID string    `json:"org_id"`
	UserID         string    `json:"user_id"`
	Permission     uint32    `json:"permission"`
	Audience       string    `json:"aud"`
	Jti            string    `json:"jti"`
	ExpireAt       time.Time `json:"exp"`
	IssuedAt       time.Time `json:"iat"`
	Provider       string    `json:"provider"` // google, github, etc
}

type JWTToken struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
	IssuedAt     time.Time
}

// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(organizationID, userID, jti string, exp time.Duration) (*JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"org_id":   organizationID,
		"exp":      time.Now().Add(exp).Unix(),
		"iat":      time.Now().Unix(),
		"jti":      jti,
		"provider": "authbase",
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		AccessToken: tokenString,
		ExpireAt:    time.Now().Add(exp),
	}, nil
}

// VerifyJWTToken verifies the JWT token
func VerifyJWTToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	expireAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	issuedAt, err := claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}

	return &Claims{
		UserID:         claims["user_id"].(string),
		OrganizationID: claims["org_id"].(string),
		Jti:            claims["jti"].(string),
		Provider:       claims["provider"].(string),
		ExpireAt:       expireAt.Time,
		IssuedAt:       issuedAt.Time,
	}, nil
}
