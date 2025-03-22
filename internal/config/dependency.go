package config

import (
	integrationMiddleware "github.com/1stpay/1stpay/internal/transport/rest/integration/middleware"
	merchantMiddleware "github.com/1stpay/1stpay/internal/transport/rest/merchant/middleware"
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
		merchantMiddleware.JWTAuthMiddleware(env.JwtSecret, usecases.UserUsecase),
		integrationMiddleware.APIKeyAuthMiddleware(usecases.MerchantAPIKeyUsecase),
	}

	return &Dependencies{
		Repos:       repos,
		Usecases:    usecases,
		Controllers: controllers,
		Middleware:  mw,
		Services:    services,
	}
}
