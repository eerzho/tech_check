package main

import (
	"context"
	"fmt"
	"tech_check/internal/app"
	"tech_check/internal/model"
	"tech_check/internal/repo/mongo_repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	app := app.MustNew()

	ctx := context.TODO()
	dropAllCollections(ctx, app.Mng)

	permissionRepo := mongo_repo.NewPermission(app.Mng)
	roleRepo := mongo_repo.NewRole(app.Mng)
	userRepo := mongo_repo.NewUser(app.Mng)
	categoryRepo := mongo_repo.NewCategory(app.Mng)
	questionRepo := mongo_repo.NewQuestion(app.Mng)

	categories := getAllCategories()
	for _, category := range categories {
		err := categoryRepo.Create(ctx, &category)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 10; i++ {
			question := model.Question{
				Text:       fmt.Sprintf("test question for %s number %d", category.Name, i),
				CategoryID: category.ID,
			}
			err := questionRepo.Create(ctx, &question)
			if err != nil {
				panic(err)
			}
		}
	}

	permissions := getAllPermissions()
	adminPermissions := make([]primitive.ObjectID, len(permissions))
	for index, permission := range permissions {
		err := permissionRepo.Create(ctx, &permission)
		if err != nil {
			panic(err)
		}
		adminPermissions[index] = permission.ID
	}

	adminRole := getAdminRole()
	adminRole.PermissionIDs = adminPermissions
	err := roleRepo.Create(ctx, adminRole)
	if err != nil {
		panic(err)
	}

	adminUser := getAdminUser()
	adminUser.RoleIDs = []primitive.ObjectID{adminRole.ID}
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

func dropAllCollections(ctx context.Context, db *mongo.Database) {
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

func getAllCategories() []model.Category {
	return []model.Category{
		{Name: "sql", Slug: "sql", Description: "checking technical sql skills"},
		{Name: "golang", Slug: "golang", Description: "checking technical Go skills"},
		{Name: "php", Slug: "php", Description: "checking technical PHP skills"},
		{Name: "js", Slug: "js", Description: "checking technical JavaScript skills"},
		{Name: "ts", Slug: "ts", Description: "checking technical TypeScript skills"},
		{Name: "python", Slug: "python", Description: "checking technical Python skills"},
		{Name: "vue", Slug: "vue", Description: "checking technical Vue.js skills"},
	}
}

func getDefaultUser() *model.User {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &model.User{
		Email:    "default@test.com",
		Name:     "default",
		Password: string(passwordHash),
	}
}

func getAdminUser() *model.User {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &model.User{
		Email:    "admin@test.com",
		Name:     "admin",
		Password: string(passwordHash),
	}
}

func getAdminRole() *model.Role {
	return &model.Role{
		Name: "Admin",
		Slug: "admin",
	}
}

func getAllPermissions() []model.Permission {
	permissionsReadPermission := model.Permission{
		Name: "Permission Read",
		Slug: "permission-read",
	}
	permissionsCreatePermission := model.Permission{
		Name: "Permission Create",
		Slug: "permission-create",
	}
	permissionsDeletePermission := model.Permission{
		Name: "Permission Delete",
		Slug: "permission-delete",
	}
	permissionsUpdatePermission := model.Permission{
		Name: "Permission Update",
		Slug: "permission-update",
	}

	usersReadPermission := model.Permission{
		Name: "User Read",
		Slug: "user-read",
	}
	usersCreatePermission := model.Permission{
		Name: "User Create",
		Slug: "user-create",
	}
	usersDeletePermission := model.Permission{
		Name: "User Delete",
		Slug: "user-delete",
	}
	usersUpdatePermission := model.Permission{
		Name: "User Update",
		Slug: "user-update",
	}

	rolesReadPermission := model.Permission{
		Name: "Role Read",
		Slug: "role-read",
	}
	rolesCreatePermission := model.Permission{
		Name: "Role Create",
		Slug: "role-create",
	}
	rolesDeletePermission := model.Permission{
		Name: "Role Delete",
		Slug: "role-delete",
	}
	rolesUpdatePermission := model.Permission{
		Name: "Role Update",
		Slug: "role-update",
	}

	return []model.Permission{
		permissionsReadPermission,
		permissionsCreatePermission,
		permissionsDeletePermission,
		permissionsUpdatePermission,
		usersReadPermission,
		usersCreatePermission,
		usersDeletePermission,
		usersUpdatePermission,
		rolesReadPermission,
		rolesCreatePermission,
		rolesDeletePermission,
		rolesUpdatePermission,
	}
}
