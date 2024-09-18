package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string               `bson:"name" json:"name"`
	Slug          string               `bson:"slug" json:"slug"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
	PermissionIDs []primitive.ObjectID `bson:"permission_ids" json:"permission_ids"`
}
