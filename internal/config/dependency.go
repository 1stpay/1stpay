package config

import (
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/middleware"
	"gorm.io/gorm"
)

type Dependencies struct {
	Repos       *Repos
	Usecases    *Usecases
	Controllers *Controllers
	Middleware  *Middleware
	Services    *Services
}

func NewDependencies(db *gorm.DB, env *Env) *Dependencies {
	repos := NewRepositories(db)

	services := NewServices(repos, env)
	usecases := NewUsecases(db, repos, services)

	controllers := NewControllers(usecases)
	mw := &Middleware{
		middleware.JWTAuthMiddleware(env.JwtSecret, usecases.UserUsecase),
	}

	return &Dependencies{
		Repos:       repos,
		Usecases:    usecases,
		Controllers: controllers,
		Middleware:  mw,
		Services:    services,
	}
}
