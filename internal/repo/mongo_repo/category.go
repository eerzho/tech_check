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

type Category struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewCategory(db *mongo.Database) *Category {
	return &Category{
		maxListCount: 200,
		collection:   db.Collection(def.TableCategories.String()),
	}
}

func (c *Category) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Category, *dto.Pagination, error) {
	const op = "mongo_repo.Category.List"

	if count > c.maxListCount {
		count = c.maxListCount
	}

	filter := bson.M{}
	for key, value := range filters {
		if key == "name" || key == "description" {
			filter[key] = bson.M{"$regex": value, "$options": "i"}
		} else if key == "slug" {
			filter[key] = value
		}
	}

	sort := bson.D{}
	for key, value := range sorts {
		if key == "created_at" ||
			key == "updated_at" ||
			key == "slug" ||
			key == "name" {
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

	cursor, err := c.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var categories []model.Category
	err = cursor.All(ctx, &categories)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.Pagination{
		Page:  page,
		Count: count,
		Total: int(total),
	}

	return categories, &pagination, nil
}

func (c *Category) Create(ctx context.Context, category *model.Category) error {
	const op = "mongo_repo.Category.Create"

	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := c.collection.InsertOne(ctx, &category)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Category) CountBySlug(ctx context.Context, slug string) (int, error) {
	const op = "mongo_repo.Category.CountBySlug"

	filter := bson.M{"slug": bson.M{"$regex": slug, "$options": "i"}}
	count, err := c.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}

func (c *Category) GetByID(ctx context.Context, id string) (*model.Category, error) {
	const op = "mongo_repo.Category.GetById"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var category model.Category

	err = c.collection.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &category, nil
}

func (c *Category) Update(ctx context.Context, category *model.Category) error {
	const op = "mongo_repo.Category.Update"

	category.UpdatedAt = time.Now()

	filter := bson.M{"_id": category.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        category.Name,
			"description": category.Description,
			"updated_at":  category.UpdatedAt,
		},
	}

	result, err := c.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (c *Category) Delete(ctx context.Context, id string) error {
	const op = "mongo_repo.Category.Delete"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}

	result, err := c.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}
