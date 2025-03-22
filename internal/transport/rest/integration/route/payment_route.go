package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
)

func NewPaymentRouter(env *config.Env, group *gin.RouterGroup, deps *config.Dependencies) {
	c := deps.Controllers.IntegrationPaymentController
	rates := group.Group("/payments")
	{
		rates.POST("/", c.Create)
	}
}
