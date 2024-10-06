package mongo_repo

import (
	"context"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionQuestion struct {
	collection *mongo.Collection
}

func NewSessionQuestion(db *mongo.Database) *SessionQuestion {
	return &SessionQuestion{
		collection: db.Collection(def.TableSessionQuestions.String()),
	}
}

func (s *SessionQuestion) Create(ctx context.Context, sessionQuestion *model.SessionQuestion) error {
	const op = "mongo_repo.SessionQuestion.Create"

	sessionQuestion.ID = primitive.NewObjectID()
	sessionQuestion.CreatedAt = time.Now()
	sessionQuestion.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, &sessionQuestion)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
