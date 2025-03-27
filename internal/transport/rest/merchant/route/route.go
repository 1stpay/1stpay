package route

import (
	_ "github.com/1stpay/1stpay/docs/merchant"
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	publicRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
