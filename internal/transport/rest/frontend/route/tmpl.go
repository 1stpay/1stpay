package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(env *config.Env, group *gin.RouterGroup) {
	rates := group.Group("/rates")
	{
		rates.POST("/ping", rest.Ping)
		rates.GET("/ping", rest.Ping)
	}
}
