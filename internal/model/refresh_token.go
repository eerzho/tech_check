package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	IP        string             `bson:"ip"`
	Hash      string             `bson:"hash"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
}
