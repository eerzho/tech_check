package mongo_repo

import (
	"context"
	"errors"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshToken struct {
	collection *mongo.Collection
}

func NewRefreshToken(db *mongo.Database) *RefreshToken {
	return &RefreshToken{
		collection: db.Collection(def.TableRefreshTokens.String()),
	}
}

func (r *RefreshToken) DeleteByUser(ctx context.Context, user *model.User) error {
	const op = "mongo_repo.RefreshToken.DeleteByUser"

	filter := bson.M{"user_id": user.ID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (r *RefreshToken) Create(ctx context.Context, refreshToken *model.RefreshToken) error {
	const op = "mongo_repo.RefreshToken.Create"

	refreshToken.ID = primitive.NewObjectID()
	refreshToken.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RefreshToken) GetByUserAndID(ctx context.Context, user *model.User, id string) (*model.RefreshToken, error) {
	const op = "mongo_repo.RefreshToken.GetByUserAndID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     idObj,
		"user_id": user.ID,
	}

	var refreshToken model.RefreshToken
	err = r.collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}
