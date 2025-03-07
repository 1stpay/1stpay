package transport

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/route"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(env *config.Env, db *gorm.DB, gin *gin.Engine) {
	frontendGroup := gin.Group("/frontend")
	route.SetupFrontendRoutes(env, db, frontendGroup)

}
