package helpers

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/middleware"
	"github.com/gin-gonic/gin"
)

func GetUserOrAbort(c *gin.Context, userUsecase usecase.UserUsecase) (model.User, bool) {
	userData, exists := c.Get(middleware.ContextUserKey)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return model.User{}, false
	}
	user, ok := userData.(model.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		return model.User{}, false
	}
	return user, true
}
