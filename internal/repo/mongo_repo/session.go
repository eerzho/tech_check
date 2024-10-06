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

type Session struct {
	collection *mongo.Collection
}

func NewSession(db *mongo.Database) *Session {
	return &Session{
		collection: db.Collection(def.TableSessions.String()),
	}
}

func (s *Session) Create(ctx context.Context, session *model.Session) error {
	const op = "mongo_repo.Session.Create"

	session.ID = primitive.NewObjectID()
	session.CreatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, &session)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Session) IsExistsActive(ctx context.Context, user *model.User) (bool, error) {
	const op = "mongo_repo.Session.IsExistsActive"

	filter := bson.M{
		"user_id":     user.ID,
		"finished_at": nil,
	}
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (s *Session) GetByID(ctx context.Context, id string) (*model.Session, error) {
	const op = "mongo_repo.Session.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var session model.Session

	err = s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &session, nil
}
