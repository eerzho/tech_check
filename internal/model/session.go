package model

import (
	"tech_check/internal/def"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	CategoryID primitive.ObjectID `bson:"category_id" json:"category_id"`
	Grade      def.GradeName      `bson:"grade" json:"grade"`
	Summary    string             `bson:"summary" json:"summary"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	FinishedAt *time.Time         `bson:"finished_at" json:"finished_at"`
}
