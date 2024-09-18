package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text       string             `bson:"text" json:"text"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}
