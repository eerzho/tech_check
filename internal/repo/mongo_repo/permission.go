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

type Permission struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewPermission(db *mongo.Database) *Permission {
	return &Permission{
		maxListCount: 200,
		collection:   db.Collection(def.TablePermissions.String()),
	}
}

func (p *Permission) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Permission, *dto.Pagination, error) {
	const op = "mongo_repo.Permission.List"

	if count > p.maxListCount {
		count = p.maxListCount
	}

	filter := bson.M{}
	for key, value := range filters {
		if key == "name" {
			filter[key] = bson.M{"$regex": value, "$options": "i"}
		} else if key == "slug" {
			filter[key] = value
		}
	}

	sort := bson.D{}
	for key, value := range sorts {
		if key == "created_at" ||
			key == "updated_at" ||
			key == "name" ||
			key == "slug" {
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

	cursor, err := p.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var permissions []model.Permission
	err = cursor.All(ctx, &permissions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := p.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.Pagination{
		Page:  page,
		Count: count,
		Total: int(total),
	}

	return permissions, &pagination, nil
}

func (p *Permission) Create(ctx context.Context, permission *model.Permission) error {
	const op = "mongo_repo.Permission.Create"

	permission.ID = primitive.NewObjectID()
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	_, err := p.collection.InsertOne(ctx, &permission)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Permission) CountBySlug(ctx context.Context, slug string) (int, error) {
	const op = "mongo_repo.Permission.CountBySlug"

	filter := bson.M{"slug": bson.M{"$regex": slug, "$options": "i"}}
	count, err := p.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}

func (p *Permission) GetByID(ctx context.Context, id string) (*model.Permission, error) {
	const op = "mongo_repo.Permission.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var permission model.Permission

	err = p.collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &permission, nil
}

func (p *Permission) Delete(ctx context.Context, id string) error {
	const op = "mongo_repo.Permission.Delete"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}

	result, err := p.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (p *Permission) Update(ctx context.Context, permission *model.Permission) error {
	const op = "mongo_repo.Permission.Update"

	permission.UpdatedAt = time.Now()

	filter := bson.M{"_id": permission.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       permission.Name,
			"updated_at": permission.UpdatedAt,
		},
	}

	result, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
