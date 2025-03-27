package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/middleware"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUsecase usecase.UserUsecase
}

type UserControllerInterfase interface {
	GetProfile(c *gin.Context)
}

func NewUserController(userUsecase usecase.UserUsecase) *UserController {
	return &UserController{
		UserUsecase: userUsecase,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the profile of the authenticated user.
// @Tags User
// @Produce json
// @Success 200 {object} restdto.UserMeResponse "User profile data"
// @Failure 400 {object} map[string]string "Invalid user or invalid user type"
// @Router /merchant/api/v1/user/me/ [get]
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
