package route

import (
	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/controller"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	user_group := group.Group("/user")
	repo := repository.NewUserRepository(db)
	uc := usecase.NewUserUsecase(repo)
	c := controller.NewUserController(uc)
	{
		user_group.GET("/me/", c.GetProfile)
	}
}
