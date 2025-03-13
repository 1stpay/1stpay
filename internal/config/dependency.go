package config

import (
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/middleware"
	"gorm.io/gorm"
)

type Dependencies struct {
	Repos       *Repos
	Usecases    *Usecases
	Controllers *Controllers
	Middleware  *Middleware
}

func NewDependencies(db *gorm.DB, env *Env) *Dependencies {
	repos := NewRepositories(db)

	usecases := NewUsecases(repos)

	controllers := NewControllers(usecases)
	mw := &Middleware{
		middleware.JWTAuthMiddleware(env.JwtSecret, usecases.UserUsecase),
	}

	return &Dependencies{
		Repos:       repos,
		Usecases:    usecases,
		Controllers: controllers,
		Middleware:  mw,
	}
}
