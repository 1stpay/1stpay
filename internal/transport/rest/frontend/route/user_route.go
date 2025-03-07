package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *config.Env, group *gin.RouterGroup) {
	rates := group.Group("/user")
	{
		rates.GET("/me/", rest.Ping)
	}
}
