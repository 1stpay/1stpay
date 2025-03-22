package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMerchantRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup, deps *config.Dependencies) {
	rates := group.Group("/merchant")
	c := deps.Controllers.MerchantMerchantController
	{
		rates.POST("/", c.MerchantCreate)
		rates.GET("/me/", c.MerchantDetail)
		rates.PUT("/me/", c.MerchantUpdate)
		rates.GET("/me/tokens/", c.MerchantTokenList)
		rates.POST("/me/tokens/", c.MerchantTokenCreate)
		rates.POST("/me/api-key/", c.MerchantAPIKeyCreate)
		rates.GET("/me/api-key/", c.MerchantAPIKeyList)
		rates.POST("/me/api-key/:id/", c.MerchantAPIKeyDeactivate)
	}
}
