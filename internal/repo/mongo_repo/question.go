package mongo_repo

import (
	"context"
	"errors"
	"fmt"
	"tech_check/internal/def"
	"tech_check/internal/dto"
	"tech_check/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Question struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewQuestion(db *mongo.Database) *Question {
	return &Question{
		maxListCount: 200,
		collection:   db.Collection(def.TableQuestions.String()),
	}
}

func (q *Question) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Question, *dto.Pagination, error) {
	const op = "mongo_repo.Question.List"

	if count > q.maxListCount {
		count = q.maxListCount
	}

	filter := bson.M{}
	for key, value := range filters {
		if key == "text" {
			filter[key] = bson.M{"$regex": value, "$options": "i"}
		} else if key == "grade" {
			_, err := def.ValidateGradeName(value)
			if err == nil {
				filter[key] = value
			}
		}
	}

	sort := bson.D{}
	for key, value := range sorts {
		if key == "created_at" ||
			key == "updated_at" {
			if value == "asc" {
				sort = append(sort, bson.E{Key: key, Value: 1})
			} else if value == "desc" {
				sort = append(sort, bson.E{Key: key, Value: -1})
			}
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * count))
	findOptions.SetLimit(int64(count))
	findOptions.SetSort(sort)

	cursor, err := q.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var questions []model.Question
	err = cursor.All(ctx, &questions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := q.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.Pagination{
		Page:  page,
		Count: count,
		Total: int(total),
	}

	return questions, &pagination, nil
}

func (q *Question) Create(ctx context.Context, question *model.Question) error {
	const op = "mongo_repo.Question.Create"

	question.ID = primitive.NewObjectID()
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	_, err := q.collection.InsertOne(ctx, &question)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (q *Question) CountBySlug(ctx context.Context, slug string) (int, error) {
	const op = "mongo_repo.Question.CountBySlug"

	filter := bson.M{"slug": bson.M{"$regex": slug, "$options": "i"}}
	count, err := q.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}

func (q *Question) GetByID(ctx context.Context, id string) (*model.Question, error) {
	const op = "mongo_repo.Question.GetById"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var question model.Question

	err = q.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &question, nil
}

func (q *Question) Update(ctx context.Context, question *model.Question) error {
	const op = "mongo_repo.Question.Update"

	question.UpdatedAt = time.Now()

	filter := bson.M{"_id": question.ID}
	update := bson.M{
		"$set": bson.M{
			"text":       question.Text,
			"grade":      question.Grade,
			"updated_at": question.UpdatedAt,
		},
	}

	result, err := q.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (q *Question) Delete(ctx context.Context, id string) error {
	const op = "mongo_repo.Question.Delete"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}

	result, err := q.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (q *Question) GetRandom(ctx context.Context, category *model.Category, grade def.GradeName, count int) ([]model.Question, error) {
	const op = "mongo_repo.Question.GetRandom"

	filter := bson.M{
		"grade":       grade,
		"category_id": category.ID,
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sample", Value: bson.M{"size": count}}},
	}
	cursor, err := q.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var questions []model.Question
	err = cursor.All(ctx, &questions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return questions, nil
}
