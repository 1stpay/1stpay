package transport

import (
	"github.com/1stpay/1stpay/internal/config"
	frontendRoute "github.com/1stpay/1stpay/internal/transport/rest/frontend/route"
	integrationRoute "github.com/1stpay/1stpay/internal/transport/rest/integration/route"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(env *config.Env, db *gorm.DB, gin *gin.Engine, deps *config.Dependencies) {
	frontendGroup := gin.Group("/frontend")
	frontendRoute.SetupFrontendRoutes(env, db, frontendGroup, deps)

	integrationGroup := gin.Group("/integration")
	integrationRoute.SetupIntegrationRoutes(env, db, integrationGroup, deps)

}
