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

type Permission struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewPermission(db *mongo.Database) *Permission {
	return &Permission{
		maxListCount: 200,
		collection:   db.Collection(constants.TablePermissions.String()),
	}
}

func (p *Permission) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.Permission, *dto.Pagination, error) {
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

	var permissions []models.Permission
	err = cursor.All(ctx, &permissions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := p.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.NewPagination(page, count, int(total))

	return permissions, pagination, nil
}

func (p *Permission) Create(ctx context.Context, permission *models.Permission) error {
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

func (p *Permission) GetByID(ctx context.Context, id string) (*models.Permission, error) {
	const op = "mongo_repo.Permission.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var permission models.Permission

	err = p.collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &permission, nil
}
