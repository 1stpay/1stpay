package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
)

func NewMerchantRouter(env *config.Env, group *gin.RouterGroup) {
	rates := group.Group("/merchants")
	{
		rates.GET("/", rest.Ping)
		rates.POST("/", rest.Ping)
		rates.GET("/:merchant_id/", rest.Ping)
		rates.PUT("/:merchant_id/", rest.Ping)
		rates.GET("/:merchant_id/blockchains/", rest.Ping)
		rates.POST("/:merchant_id/blockchains/", rest.Ping)
		rates.GET("/:merchant_id/tokens/", rest.Ping)
		rates.POST("/:merchant_id/tokens/", rest.Ping)
	}
}
