package app

import (
	"log/slog"
	"tech_check/internal/config"
	"tech_check/internal/repositories/mongo_repositories"
	"tech_check/internal/services"
	"tech_check/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	App struct {
		Cfg   *config.Config
		Lg    *slog.Logger
		Mng   *mongo.Database
		Srvcs *srvcs
	}

	repos struct {
		User            *mongo_repositories.User
		Role            *mongo_repositories.Role
		Permission      *mongo_repositories.Permission
		RefreshToken    *mongo_repositories.RefreshToken
		Category        *mongo_repositories.Category
		Question        *mongo_repositories.Question
		Session         *mongo_repositories.Session
		SessionQuestion *mongo_repositories.SessionQuestion
	}

	srvcs struct {
		User            *services.User
		Role            *services.Role
		Permission      *services.Permission
		RefreshToken    *services.RefreshToken
		Auth            *services.Auth
		Category        *services.Category
		Question        *services.Question
		Session         *services.Session
		SessionQuestion *services.SessionQuestion
	}
)

func MustNew() *App {
	cfg := mustSetupConfig()
	lg := mustSetupLogger(cfg)
	mng := mustSetupMongo(cfg)

	repos := setupRepositories(mng)
	srvcs := setupServices(cfg, repos)

	return &App{
		Cfg:   cfg,
		Lg:    lg,
		Mng:   mng,
		Srvcs: srvcs,
	}
}

func setupRepositories(mng *mongo.Database) *repos {
	user := mongo_repositories.NewUser(mng)
	role := mongo_repositories.NewRole(mng)
	permission := mongo_repositories.NewPermission(mng)
	refreshToken := mongo_repositories.NewRefreshToken(mng)
	category := mongo_repositories.NewCategory(mng)
	question := mongo_repositories.NewQuestion(mng)
	session := mongo_repositories.NewSession(mng)
	sessionQuestion := mongo_repositories.NewSessionQuestion(mng)

	return &repos{
		User:            user,
		Role:            role,
		Permission:      permission,
		RefreshToken:    refreshToken,
		Category:        category,
		Question:        question,
		Session:         session,
		SessionQuestion: sessionQuestion,
	}
}

func setupServices(cfg *config.Config, repos *repos) *srvcs {
	permission := services.NewPermission(repos.Permission)
	role := services.NewRole(repos.Role, permission)
	user := services.NewUser(repos.User, role)
	refreshToken := services.NewRefreshToken(repos.RefreshToken)
	auth := services.NewAuth(cfg.Google.ClientID, cfg.JWT.Secret, user, refreshToken)
	category := services.NewCategory(repos.Category)
	question := services.NewQuestion(repos.Question, category)
	sessionQuestion := services.NewSessionQuestion(repos.SessionQuestion)
	session := services.NewSession(repos.Session, category, question, sessionQuestion)

	return &srvcs{
		User:            user,
		Role:            role,
		Permission:      permission,
		RefreshToken:    refreshToken,
		Auth:            auth,
		Category:        category,
		Question:        question,
		Session:         session,
		SessionQuestion: sessionQuestion,
	}
}

func mustSetupConfig() *config.Config {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	return cfg
}

func mustSetupLogger(cfg *config.Config) *slog.Logger {
	lg := utils.NewLogger(cfg.Log.Level, cfg.Log.Format)

	return lg
}

func mustSetupMongo(cfg *config.Config) *mongo.Database {
	mng, err := utils.NewMongo(cfg.Mongo.DB, cfg.Mongo.URL)
	if err != nil {
		panic(err)
	}

	return mng
}
