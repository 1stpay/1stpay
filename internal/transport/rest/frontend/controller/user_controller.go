package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
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
	user, err := uc.UserUsecase.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user data in context"})
		return
	}

	c.JSON(http.StatusOK, restdto.UserMeResponse{
		Id:    user.ID.String(),
		Email: user.Email,
	})
}
