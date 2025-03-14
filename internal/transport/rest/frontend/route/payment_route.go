package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
)

func NewPaymentRouter(env *config.Env, group *gin.RouterGroup, deps *config.Dependencies) {
	rates := group.Group("/payments")
	c := deps.Controllers.FrontendPaymentController
	{
		rates.GET("/", rest.Ping)
		rates.POST("/", c.Create)
		rates.GET("/:id/", rest.Ping)
	}
}
