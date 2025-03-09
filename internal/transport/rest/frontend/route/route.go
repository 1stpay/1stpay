package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupFrontendRoutes(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	publicRouter := group.Group("/api/v1")
	protectedRouter := publicRouter.Group("")
	protectedRouter.Use(middleware.JWTAuthMiddleware("hehe"))
	NewAuthRouter(env, db, publicRouter)
	NewUserRouter(env, db, protectedRouter)
	NewMerchantRouter(env, db, protectedRouter)
	NewPaymentRouter(env, protectedRouter)

	NewPingRouter(env, publicRouter)
}
