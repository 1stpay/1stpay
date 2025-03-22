package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/gin-gonic/gin"
)

type authController struct {
	AuthUsecase usecase.AuthUsecase
}

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

func NewAuthController(authUsecase usecase.AuthUsecase) AuthController {
	return &authController{
		AuthUsecase: authUsecase,
	}
}

func (ac *authController) Register(c *gin.Context) {
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

func (ac *authController) Login(c *gin.Context) {
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
