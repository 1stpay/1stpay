package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupFrontendRoutes(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	publicRouter := group.Group("/api/v1")
	NewPaymentRouter(env, publicRouter, deps)
}
