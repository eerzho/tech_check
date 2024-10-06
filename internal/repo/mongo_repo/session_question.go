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

type SessionQuestion struct {
	collection *mongo.Collection
}

func NewSessionQuestion(db *mongo.Database) *SessionQuestion {
	return &SessionQuestion{
		collection: db.Collection(def.TableSessionQuestions.String()),
	}
}

func (s *SessionQuestion) Create(ctx context.Context, question *model.SessionQuestion) error {
	const op = "mongo_repo.SessionQuestion.Create"

	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, &question)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SessionQuestion) List(ctx context.Context, session *model.Session) ([]model.SessionQuestion, error) {
	const op = "mongo_repo.SessionQuestion.List"

	filter := bson.M{"session_id": session.ID}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var questions []model.SessionQuestion
	err = cursor.All(ctx, &questions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}

func (s *SessionQuestion) GetByID(ctx context.Context, session *model.Session, id string) (*model.SessionQuestion, error) {
	const op = "mongo_repo.SessionQuestion.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{
		"_id":        idObj,
		"session_id": session.ID,
	}
	var question model.SessionQuestion
	err = s.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}

func (s *SessionQuestion) Update(ctx context.Context, question *model.SessionQuestion) error {
	const op = "mongo_repo.SessionQuestion.Update"

	question.UpdatedAt = time.Now()
	filter := bson.M{"_id": question.ID}
	update := bson.M{
		"$set": bson.M{
			"answer":    question.Answer,
			"summary":   question.Summary,
			"update_at": question.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}
