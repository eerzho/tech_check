package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email     string               `bson:"email" json:"email"`
	Name      string               `bson:"name" json:"name"`
	Password  string               `bson:"password" json:"-"`
	Avatar    string               `bson:"avatar" json:"avatar"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
	RoleIDs   []primitive.ObjectID `bson:"role_ids" json:"role_ids"`
}
