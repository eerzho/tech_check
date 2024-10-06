package app

import (
	"log/slog"
	"tech_check/internal/config"
	"tech_check/internal/repo/mongo_repo"
	"tech_check/internal/srvc"
	"tech_check/internal/util"

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
		User            *mongo_repo.User
		Role            *mongo_repo.Role
		Permission      *mongo_repo.Permission
		RefreshToken    *mongo_repo.RefreshToken
		Category        *mongo_repo.Category
		Question        *mongo_repo.Question
		Session         *mongo_repo.Session
		SessionQuestion *mongo_repo.SessionQuestion
	}

	srvcs struct {
		User            *srvc.User
		Role            *srvc.Role
		Permission      *srvc.Permission
		RefreshToken    *srvc.RefreshToken
		Auth            *srvc.Auth
		Category        *srvc.Category
		Question        *srvc.Question
		Session         *srvc.Session
		SessionQuestion *srvc.SessionQuestion
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
	user := mongo_repo.NewUser(mng)
	role := mongo_repo.NewRole(mng)
	permission := mongo_repo.NewPermission(mng)
	refreshToken := mongo_repo.NewRefreshToken(mng)
	category := mongo_repo.NewCategory(mng)
	question := mongo_repo.NewQuestion(mng)
	session := mongo_repo.NewSession(mng)
	sessionQuestion := mongo_repo.NewSessionQuestion(mng)

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
	permission := srvc.NewPermission(repos.Permission)
	role := srvc.NewRole(repos.Role, permission)
	user := srvc.NewUser(repos.User, role)
	refreshToken := srvc.NewRefreshToken(repos.RefreshToken)
	auth := srvc.NewAuth(cfg.Google.ClientID, cfg.JWT.Secret, user, refreshToken)
	category := srvc.NewCategory(repos.Category)
	question := srvc.NewQuestion(repos.Question, category)
	sessionQuestion := srvc.NewSessionQuestion(repos.SessionQuestion)
	session := srvc.NewSession(repos.Session, category, question, sessionQuestion)

	return &srvcs{
		User:         user,
		Role:         role,
		Permission:   permission,
		RefreshToken: refreshToken,
		Auth:         auth,
		Category:     category,
		Question:     question,
		Session:      session,
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
	lg := util.NewLogger(cfg.Log.Level, cfg.Log.Format)

	return lg
}

func mustSetupMongo(cfg *config.Config) *mongo.Database {
	mng, err := util.NewMongo(cfg.Mongo.DB, cfg.Mongo.URL)
	if err != nil {
		panic(err)
	}

	return mng
}
