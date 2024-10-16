package main

import (
	"context"
	"fmt"
	"tech_check/internal/config"
	"tech_check/internal/constants"
	"tech_check/internal/models"
	"tech_check/internal/repositories/mongo_repositories"
	"tech_check/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	mng, err := utils.NewMongo(cfg.Mongo.DB, cfg.Mongo.URL)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()
	dropDB(ctx, mng)

	categoryRepo := mongo_repositories.NewCategory(mng)
	questionRepo := mongo_repositories.NewQuestion(mng)
	for _, category := range getAllCategories() {
		err = categoryRepo.Create(ctx, &category)
		if err != nil {
			panic(err)
		}

		for i := 0; i < 10; i++ {
			question := models.Question{
				Text:       fmt.Sprintf("%s category test question %d", category.Name, i),
				Grade:      constants.GradeJunior,
				CategoryID: category.ID,
			}
			err = questionRepo.Create(ctx, &question)
			if err != nil {
				panic(err)
			}
		}
	}

	var permissionIDs []primitive.ObjectID
	permissionRepo := mongo_repositories.NewPermission(mng)
	for _, permission := range getAllPermissions() {
		err = permissionRepo.Create(ctx, &permission)
		if err != nil {
			panic(err)
		}
		permissionIDs = append(permissionIDs, permission.ID)
	}

	roleRepo := mongo_repositories.NewRole(mng)
	role := getRole(permissionIDs)
	err = roleRepo.Create(ctx, role)
	if err != nil {
		panic(err)
	}

	userRepo := mongo_repositories.NewUser(mng)
	adminUser := getAdminUser(role.ID)
	err = userRepo.Create(ctx, adminUser)
	if err != nil {
		panic(err)
	}

	defaultUser := getDefaultUser()
	err = userRepo.Create(ctx, defaultUser)
	if err != nil {
		panic(err)
	}
}

func dropDB(ctx context.Context, db *mongo.Database) {
	collections, err := db.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	for _, collectionName := range collections {
		err := db.Collection(collectionName).Drop(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func getDefaultUser() *models.User {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &models.User{
		Email:    "default@test.com",
		Name:     "default",
		Password: string(passwordHash),
	}
}

func getAdminUser(roleID primitive.ObjectID) *models.User {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &models.User{
		Email:    "admin@test.com",
		Name:     "admin",
		Password: string(passwordHash),
		RoleIDs:  []primitive.ObjectID{roleID},
	}
}

func getRole(permissionIDs []primitive.ObjectID) *models.Role {
	return &models.Role{
		Name:          "Admin",
		Slug:          "admin",
		PermissionIDs: permissionIDs,
	}
}

func getAllPermissions() []models.Permission {
	return []models.Permission{
		{Name: "Category read", Slug: "category-read"},
		{Name: "Category create", Slug: "category-create"},
		{Name: "Category edit", Slug: "category-edit"},
		{Name: "Category delete", Slug: "category-delete"},

		{Name: "Permission read", Slug: "permission-read"},
		{Name: "Permission create", Slug: "permission-create"},
		{Name: "Permission edit", Slug: "permission-edit"},
		{Name: "Permission delete", Slug: "permission-delete"},

		{Name: "Question read", Slug: "question-read"},
		{Name: "Question create", Slug: "question-create"},
		{Name: "Question edit", Slug: "question-edit"},
		{Name: "Question delete", Slug: "question-delete"},

		{Name: "Role read", Slug: "role-read"},
		{Name: "Role create", Slug: "role-create"},
		{Name: "Role edit", Slug: "role-edit"},
		{Name: "Role delete", Slug: "role-delete"},

		{Name: "User read", Slug: "user-read"},
		{Name: "User create", Slug: "user-create"},
		{Name: "User edit", Slug: "user-edit"},
		{Name: "User delete", Slug: "user-delete"},
	}
}

func getAllCategories() []models.Category {
	return []models.Category{
		{Name: "sql", Slug: "sql", Description: "checking technical sql skills"},
		{Name: "golang", Slug: "golang", Description: "checking technical Go skills"},
		{Name: "php", Slug: "php", Description: "checking technical PHP skills"},
		{Name: "js", Slug: "js", Description: "checking technical JavaScript skills"},
		{Name: "ts", Slug: "ts", Description: "checking technical TypeScript skills"},
		{Name: "python", Slug: "python", Description: "checking technical Python skills"},
		{Name: "vue", Slug: "vue", Description: "checking technical Vue.js skills"},
	}
}
