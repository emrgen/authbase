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
	Username   string `json:"username"`
	Org        string `json:"org"`
	Permission uint32 `json:"permission"`
}

type JWTToken struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
	IssuedAt     time.Time
}

// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(userID, organizationID string) (*JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userID,
		"org":      organizationID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		AccessToken: tokenString,
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

	return &Claims{
		Username: claims["username"].(string),
		Org:      claims["org"].(string),
	}, nil
}
