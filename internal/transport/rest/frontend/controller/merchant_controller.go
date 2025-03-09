package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/gin-gonic/gin"
)

type MerchantController struct {
	MerchantUsecase usecase.MerchantUsecaseInterface
	UserUsecase     usecase.UserUsecaseInterface
}

type MerchantControllerInterface interface {
	CreateMerchant(c *gin.Context)
}

func NewMerchantController(merchantUsecase usecase.MerchantUsecaseInterface, userUsecase usecase.UserUsecaseInterface) *MerchantController {
	return &MerchantController{
		MerchantUsecase: merchantUsecase,
		UserUsecase:     userUsecase,
	}
}

func (u *MerchantController) CreateMerchant(c *gin.Context) {
	var req restdto.MerchantCreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, err := u.UserUsecase.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user"})
	}
	merchant, err := u.MerchantUsecase.Create(req, user.ID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error while merchant create"})
	}
	c.JSON(http.StatusCreated, restdto.MerchantCreateResponseDTO{
		ID:             merchant.ID,
		CreatedAt:      merchant.CreatedAt,
		UpdatedAt:      merchant.UpdatedAt,
		UserID:         merchant.UserID,
		Name:           merchant.Name,
		CommissionRate: merchant.CommissionRate,
	})
}
