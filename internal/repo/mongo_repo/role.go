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

type Role struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewRole(db *mongo.Database) *Role {
	return &Role{
		maxListCount: 200,
		collection:   db.Collection(def.TableRoles.String()),
	}
}

func (r *Role) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.Role, *dto.Pagination, error) {
	const op = "mongo_repo.Role.List"

	if count > r.maxListCount {
		count = r.maxListCount
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

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var roles []model.Role
	err = cursor.All(ctx, &roles)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.Pagination{
		Page:  page,
		Count: count,
		Total: int(total),
	}

	return roles, &pagination, nil
}

func (r *Role) Create(ctx context.Context, role *model.Role) error {
	const op = "mongo_repo.Role.Create"

	role.ID = primitive.NewObjectID()
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, &role)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Role) CountBySlug(ctx context.Context, slug string) (int, error) {
	const op = "mongo_repo.Role.CountBySlug"

	filter := bson.M{"slug": bson.M{"$regex": slug, "$options": "i"}}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(count), nil
}

func (r *Role) GetByID(ctx context.Context, id string) (*model.Role, error) {
	const op = "mongo_repo.Role.GetByID"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var role model.Role

	err = r.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

func (r *Role) Delete(ctx context.Context, id string) error {
	const op = "mongo_repo.Role.Delete"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (r *Role) Update(ctx context.Context, role *model.Role) error {
	const op = "mongo_repo.Role.Update"

	role.UpdatedAt = time.Now()

	filter := bson.M{"_id": role.ID}
	update := bson.M{
		"$set": bson.M{
			"name":           role.Name,
			"updated_at":     role.UpdatedAt,
			"permission_ids": role.PermissionIDs,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}
