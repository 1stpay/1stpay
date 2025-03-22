package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	c := deps.Controllers.MerchantAuthController
	auth_group := group.Group("/auth")
	{
		auth_group.POST("/register/", c.Register)
		auth_group.POST("/login/", c.Login)
	}
}
