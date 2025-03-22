package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	user_group := group.Group("/user")
	c := deps.Controllers.MerchantUserController
	{
		user_group.GET("/me/", c.GetProfile)
	}
}
