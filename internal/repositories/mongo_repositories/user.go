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

type User struct {
	maxListCount int
	collection   *mongo.Collection
}

func NewUser(db *mongo.Database) *User {
	return &User{
		maxListCount: 200,
		collection:   db.Collection(constants.TableUsers.String()),
	}
}

func (u *User) List(ctx context.Context, page, count int, filters, sorts map[string]string) ([]models.User, *dto.Pagination, error) {
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

	var users []models.User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	total, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	pagination := dto.NewPagination(page, count, int(total))

	return users, pagination, nil
}

func (u *User) Create(ctx context.Context, user *models.User) error {
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

func (u *User) GetByID(ctx context.Context, id string) (*models.User, error) {
	const op = "mongo_repo.User.GetById"

	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filter := bson.M{"_id": idObj}
	var user models.User

	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) Update(ctx context.Context, user *models.User) error {
	const op = "mongo_repo.User.Update"

	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"avatar":     user.Avatar,
			"updated_at": user.UpdatedAt,
			"role_ids":   user.RoleIDs,
		},
	}

	result, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("%s: %w", op, constants.ErrNotFound)
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
		return fmt.Errorf("%s: %w", op, constants.ErrNotFound)
	}

	return nil
}

func (u *User) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const op = "mongo_repo.User.ExistsByEmail"

	filter := bson.M{"email": email}

	count, err := u.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "mongo_repo.User.GetByEmail"

	filter := bson.M{"email": email}
	var user models.User

	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%s: %w", op, constants.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (u *User) HasPermission(ctx context.Context, user *models.User, permissionSlug string) (bool, error) {
	const op = "mongo_repo.User.HasPermission"

	if len(user.RoleIDs) == 0 {
		return false, nil
	}

	roleFilter := bson.M{
		"_id": bson.M{"$in": user.RoleIDs},
	}
	rolesCursor, err := u.collection.Database().
		Collection(constants.TableRoles.String()).
		Find(ctx, roleFilter)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer rolesCursor.Close(ctx)

	var roles []models.Role
	err = rolesCursor.All(ctx, &roles)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var permissionIDs []primitive.ObjectID
	for _, role := range roles {
		permissionIDs = append(permissionIDs, role.PermissionIDs...)
	}

	if len(permissionIDs) == 0 {
		return false, nil
	}

	permissionFilter := bson.M{
		"_id":  bson.M{"$in": permissionIDs},
		"slug": permissionSlug,
	}
	count, err := u.collection.Database().
		Collection(constants.TablePermissions.String()).
		CountDocuments(ctx, permissionFilter)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}
