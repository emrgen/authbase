package x

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	// RefreshTokenDuration is the duration for the refresh token
	// TODO: this should be configurable in the future
	RefreshTokenDuration = 7 * 24 * time.Hour
	// AccessTokenDuration is the duration for the access token
	// TODO: this should be configurable in the future
	AccessTokenDuration = 15 * time.Minute
)

func jwtSecret() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		logrus.Error("jwt is not set")
		panic("JWT_SECRET is not set")
	}

	return secretKey
}

type Claims struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ClientID  string    `json:"client_id"`
	ProjectID string    `json:"project_id"`
	AccountID string    `json:"account_id"`
	Audience  string    `json:"aud"`
	Jti       string    `json:"jti"`
	ExpireAt  time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
	Provider  string    `json:"provider"` // google, github, etc
	Scopes    []string  `json:"scopes"`
}

type JWTToken struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
	IssuedAt     time.Time
}

// GenerateJWTToken generates a JWT token for the user
func GenerateJWTToken(claims Claims) (*JWTToken, error) {
	claim := jwt.MapClaims{
		"username":   claims.Username,
		"email":      claims.Email,
		"account_id": claims.AccountID,
		"project_id": claims.ProjectID,
		"client_id":  claims.ClientID,
		"exp":        claims.ExpireAt.Unix(),
		"iat":        time.Now().Unix(),
		"jti":        claims.Jti,
		"provider":   "authbase",
		"scopes":     claims.Scopes,
	}

	// Generate the access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(jwtSecret()))
	if err != nil {
		return nil, err
	}

	// Generate the refresh token
	claim["exp"] = time.Now().Add(RefreshTokenDuration).Unix()
	claim["iat"] = time.Now().Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	refreshToken, err := token.SignedString([]byte(jwtSecret()))
	if err != nil {
		return nil, err
	}

	return &JWTToken{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpireAt:     claims.ExpireAt,
	}, nil
}

// VerifyJWTToken verifies the JWT token
func VerifyJWTToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret()), nil
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

	if time.Now().After(expireAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	issuedAt, err := claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, fmt.Errorf("jti not found")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}

	projectID, ok := claims["project_id"].(string)
	if !ok {
		return nil, fmt.Errorf("project_id not found")
	}

	clientID, ok := claims["client_id"].(string)
	if !ok {
		return nil, fmt.Errorf("client_id not found")
	}

	provider, ok := claims["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("provider not found")
	}

	scopes, ok := claims["scopes"].([]string)
	if !ok {
		scopes = []string{}
	}

	return &Claims{
		AccountID: userID,
		ProjectID: projectID,
		ClientID:  clientID,
		Jti:       jti,
		Provider:  provider,
		ExpireAt:  expireAt.Time,
		IssuedAt:  issuedAt.Time,
		Scopes:    scopes,
	}, nil
}
