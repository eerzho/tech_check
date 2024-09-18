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

type User struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewUser(db *mongo.Database) *User {
	return &User{
		maxListCount: 200,
		collection:   db.Collection(def.TableUsers.String()),
	}
}

func (u *User) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]model.User, *dto.Pagination, error) {
	const op = "mongo_repo.User.List"

	if count > u.maxListCount {
		count = u.maxListCount
	}

	filter := bson.M{}
	for key, value := range filters {
		if key == "name" || key == "email" {
			filter[key] = bson.M{"$regex": value, "$options": "i"}
		}
	}

	sort := bson.D{}
	for key, value := range sorts {
		if key == "created_at" || key == "updated_at" {
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

	cursor, err := u.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}
	defer cursor.Close(ctx)

	var users []model.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.Pagination{
		Page:  page,
		Count: count,
		Total: int(total),
	}

	return users, &pagination, nil
}

func (u *User) Create(ctx context.Context, user *model.User) error {
	const op = "mongo_repo.User.Create"

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := u.collection.InsertOne(ctx, &user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *User) GetByID(ctx context.Context, id string) (*model.User, error) {
	const op = "mongo_repo.User.GetById"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var user model.User

	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) Update(ctx context.Context, user *model.User) error {
	const op = "mongo_repo.User.Update"

	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
			"role_ids":   user.RoleIDs,
		},
	}

	result, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (u *User) Delete(ctx context.Context, id string) error {
	const op = "mongo_repo.User.Delete"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}

	result, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, def.ErrNotFound)
	}

	return nil
}

func (u *User) IsExistsEmail(ctx context.Context, email string) (bool, error) {
	const op = "mongo_repo.User.IsExistsEmail"

	filter := bson.M{"email": email}

	count, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "mongo_repo.User.GetByEmail"

	filter := bson.M{"email": email}
	var user model.User

	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, def.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) RemoveRole(ctx context.Context, roleID string) error {
	const op = "mongo_repo.User.RemoveRole"

	roleIDObj, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{
		"role_ids": roleIDObj,
	}

	update := bson.M{
		"$pull": bson.M{
			"role_ids": roleIDObj,
		},
	}

	_, err = u.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
