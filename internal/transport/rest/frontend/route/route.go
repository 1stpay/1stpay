package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupFrontendRoutes(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	publicRouter := group.Group("/api/v1")
	NewPingRouter(env, publicRouter)
	NewAuthRouter(env, publicRouter)
	NewMerchantRouter(env, publicRouter)
	NewPaymentRouter(env, publicRouter)
	NewUserRouter(env, publicRouter)
}
