package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupMerchantRoutes(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	publicRouter := group.Group("/api/v1")
	NewBlockchainRouter(env, db, publicRouter, deps)
	NewTokenRouter(env, db, publicRouter, deps)
	NewPingRouter(env, publicRouter)

	protectedRouter := publicRouter.Group("")
	protectedRouter.Use(deps.Middleware.JWTAuth)
	NewAuthRouter(env, db, publicRouter, deps)
	NewUserRouter(env, db, protectedRouter, deps)
	NewMerchantRouter(env, db, protectedRouter, deps)
	NewPaymentRouter(env, protectedRouter, deps)

}
