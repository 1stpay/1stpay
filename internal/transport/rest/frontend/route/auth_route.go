package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	repo := repository.NewUserRepository(db)
	uc := usecase.NewAuthUsecase(repo)
	c := controller.NewAuthController(uc)
	auth_group := group.Group("/auth")
	{
		auth_group.POST("/register/", c.Register)
		auth_group.POST("/login/", c.Login)
	}
}
