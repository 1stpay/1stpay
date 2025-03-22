package route

import (
	"github.com/1stpay/1stpay/internal/config"
	rest "github.com/1stpay/1stpay/internal/transport/rest/merchant/controller"
	"github.com/gin-gonic/gin"
)

func NewPingRouter(env *config.Env, group *gin.RouterGroup) {
	group.GET("/ping", rest.Ping)
}
