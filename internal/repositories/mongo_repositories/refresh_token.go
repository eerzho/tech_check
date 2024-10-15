package mongo_repositories

import (
	"context"
	"errors"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/models"
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
		collection: db.Collection(constants.TableRefreshTokens.String()),
	}
}

func (r *RefreshToken) Delete(ctx context.Context, user *models.User) error {
	const op = "mongo_repo.RefreshToken.DeleteByUser"

	filter := bson.M{"user_id": user.ID}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, constants.ErrNotFound)
	}

	return nil
}

func (r *RefreshToken) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	const op = "mongo_repo.RefreshToken.Create"

	refreshToken.ID = primitive.NewObjectID()
	refreshToken.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RefreshToken) GetByID(ctx context.Context, user *models.User, id string) (*models.RefreshToken, error) {
	const op = "mongo_repo.RefreshToken.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":     idObj,
		"user_id": user.ID,
	}

	var refreshToken models.RefreshToken
	err = r.collection.FindOne(ctx, filter).Decode(&refreshToken)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &refreshToken, nil
}
