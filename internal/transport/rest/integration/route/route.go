package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupIntegrationRoutes(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	publicRouter := group.Group("/api/v1")
	publicRouter.Use(deps.Middleware.APIKeyAuth)
	NewPaymentRouter(env, publicRouter, deps)
}
