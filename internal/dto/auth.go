package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Token struct {
		AToken string `json:"access_token"`
		RToken string `json:"refresh_token"`
	}

	Claims struct {
		IP             string             `json:"ip"`
		UserID         primitive.ObjectID `json:"user_id"`
		RefreshTokenID primitive.ObjectID `json:"refresh_token_id"`
		jwt.RegisteredClaims
	}
)
