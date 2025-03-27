package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthUsecase usecase.AuthUsecase
}

func NewAuthController(authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
	}
}

// Register godoc
// @Summary User registration
// @Description Registers new user and returns access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body restdto.RegisterRequest true "Registration data"
// @Success 201 {object} restdto.AccessTokenResponse "Returns success token"
// @Failure 400 {object} gin.H "Incorrect data"
// @Router /merchant/api/v1/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req restdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect data"})
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

// Login godoc
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает токен доступа
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body restdto.LoginRequest true "Данные для входа"
// @Success 200 {object} restdto.AccessTokenResponse "Возвращает токен доступа"
// @Failure 400 {object} gin.H "Неверные входные данные"
// @Failure 403 {object} gin.H "Доступ запрещён"
// @Router /merchant/api/v1/auth/login [post]
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
