package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/merchant/controller"
	"github.com/gin-gonic/gin"
)

func NewPaymentRouter(env *config.Env, group *gin.RouterGroup) {
	rates := group.Group("/payments")
	{
		rates.POST("/", rest.Ping)
		rates.GET("/", rest.Ping)
		rates.GET("/:id/", rest.Ping)
	}
}
