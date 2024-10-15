package mongo_repositories

import (
	"context"
	"errors"
	"fmt"
	"tech_check/internal/constants"
	"tech_check/internal/dto"
	"tech_check/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewSession(db *mongo.Database) *Session {
	return &Session{
		maxListCount: 200,
		collection:   db.Collection(constants.TableSessions.String()),
	}
}

func (s *Session) Create(ctx context.Context, session *models.Session) error {
	const op = "mongo_repo.Session.Create"

	session.ID = primitive.NewObjectID()
	session.CreatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, &session)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Session) ExistsActive(ctx context.Context, user *models.User) (bool, error) {
	const op = "mongo_repo.Session.ExistsActive"

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

func (s *Session) GetByID(ctx context.Context, id string) (*models.Session, error) {
	const op = "mongo_repo.Session.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var session models.Session

	err = s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &session, nil
}

func (s *Session) List(ctx context.Context, user *models.User, page, count int) ([]models.Session, *dto.Pagination, error) {
	const op = "mongo_repo.Session.List"

	if count > s.maxListCount {
		count = s.maxListCount
	}

	filter := bson.M{"user_id": user.ID}
	sort := bson.D{{Key: "created_at", Value: -1}}

	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * count))
	findOptions.SetLimit(int64(count))
	findOptions.SetSort(sort)

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var sessions []models.Session
	err = cursor.All(ctx, &sessions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.NewPagination(page, count, int(total))

	return sessions, pagination, nil
}

func (s *Session) Update(ctx context.Context, session *models.Session) error {
	const op = "mongo_repo.Session.Update"

	filter := bson.M{"_id": session.ID}
	update := bson.M{
		"$set": bson.M{
			"summary":     session.Summary,
			"finished_at": session.FinishedAt,
		},
	}
	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
