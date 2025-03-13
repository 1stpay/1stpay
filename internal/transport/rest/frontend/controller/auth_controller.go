package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthUsecase usecase.AuthUsecaseInterface
}

type AuthControllerInterface interface {
	Register(c *gin.Context)
}

func NewAuthController(authUsecase usecase.AuthUsecaseInterface) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req restdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}

	_, token, err := ac.AuthUsecase.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response := restdto.AccessTokenResponse{
		AccessToken: token,
	}

	c.JSON(http.StatusCreated, response)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req restdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong input data"})
		return
	}

	_, token, err := ac.AuthUsecase.Login(req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	response := restdto.AccessTokenResponse{
		AccessToken: token,
	}
	c.JSON(http.StatusOK, response)
}
