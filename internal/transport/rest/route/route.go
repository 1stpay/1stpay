package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(env *config.Env, db *gorm.DB, gin *gin.Engine) {
	publicRouter := gin.Group("/v1")
	NewPingRouter(env, publicRouter)
}
