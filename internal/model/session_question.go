package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionQuestion struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SessionID primitive.ObjectID `bson:"session_id" json:"session_id"`
	Text      string             `bson:"text" json:"text"`
	Answer    string             `bson:"answer" json:"answer"`
	Summary   string             `bson:"summary" json:"summary"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
