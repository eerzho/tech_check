package dto

import (
	"tech_check/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Token struct {
		AToken string `json:"access_token"`
		RToken string `json:"refresh_token"`
	}

	Claims struct {
		IP             string `json:"ip"`
		UserID         string `json:"user_id"`
		RefreshTokenID string `json:"refresh_token_id"`
		jwt.RegisteredClaims
	}
)

func NewToken(aToken, rToken string) *Token {
	return &Token{
		AToken: aToken,
		RToken: rToken,
	}
}

func NewClaims(
	user *models.User,
	refreshToken *models.RefreshToken,
	ip string,
	aTokenExpiresHour int,
) *Claims {
	return &Claims{
		IP:             ip,
		UserID:         user.ID.Hex(),
		RefreshTokenID: refreshToken.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(aTokenExpiresHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
