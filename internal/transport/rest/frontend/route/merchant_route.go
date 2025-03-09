package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/repository"
	rest "github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMerchantRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	rates := group.Group("/merchants")
	merchantRepo := repository.NewMerchantRepository(db)
	userRepo := repository.NewUserRepository(db)
	merchUc := usecase.NewMerchantUsecase(merchantRepo)
	userUc := usecase.NewUserUsecase(userRepo)
	c := rest.NewMerchantController(merchUc, userUc)
	{
		rates.GET("/", rest.Ping)
		rates.POST("/", c.CreateMerchant)
		rates.GET("/:merchant_id/", rest.Ping)
		rates.PUT("/:merchant_id/", rest.Ping)
		rates.GET("/:merchant_id/blockchains/", rest.Ping)
		rates.POST("/:merchant_id/blockchains/", rest.Ping)
		rates.GET("/:merchant_id/tokens/", rest.Ping)
		rates.POST("/:merchant_id/tokens/", rest.Ping)
	}
}
