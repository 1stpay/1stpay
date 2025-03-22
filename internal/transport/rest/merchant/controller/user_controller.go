package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/middleware"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/merchant/rest_dto"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase usecase.UserUsecaseInterface
}

type UserControllerInterfase interface {
	GetProfile(c *gin.Context)
}

func NewUserController(userUsecase usecase.UserUsecaseInterface) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userData, exists := c.Get(middleware.ContextUserKey)

	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}
	user, ok := userData.(model.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		return
	}

	c.JSON(http.StatusOK, restdto.UserMeResponse{
		Id:    user.ID.String(),
		Email: user.Email,
	})
}
