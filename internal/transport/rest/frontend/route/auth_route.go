package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
)

func NewAuthRouter(env *config.Env, group *gin.RouterGroup) {
	rates := group.Group("/auth")
	{
		rates.GET("/register/", rest.Ping)
		rates.POST("/login/", rest.Ping)
		rates.POST("/logout/", rest.Ping)
	}
}
