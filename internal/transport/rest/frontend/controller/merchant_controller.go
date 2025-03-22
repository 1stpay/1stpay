package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/helpers"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/gin-gonic/gin"
)

type MerchantController struct {
	MerchantUsecase usecase.MerchantUsecase
	UserUsecase     usecase.UserUsecase
}

func NewMerchantController(merchantUsecase usecase.MerchantUsecase, userUsecase usecase.UserUsecase) *MerchantController {
	return &MerchantController{
		MerchantUsecase: merchantUsecase,
		UserUsecase:     userUsecase,
	}
}

func (u *MerchantController) MerchantCreate(c *gin.Context) {
	var req restdto.MerchantCreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.CreateMerchant(req, user.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while merchant create"})
		return
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

func (u *MerchantController) MerchantDetail(c *gin.Context) {
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.GetMerchantByUserId(user.ID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
	}
	c.JSON(http.StatusOK, restdto.MerchantDetailResponseDTO{
		ID:             merchant.ID,
		CreatedAt:      merchant.CreatedAt,
		UpdatedAt:      merchant.UpdatedAt,
		UserID:         merchant.UserID,
		Name:           merchant.Name,
		CommissionRate: merchant.CommissionRate,
	})
}

func (u *MerchantController) MerchantUpdate(c *gin.Context) {
	var req restdto.MerchantCreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.UpdateMerchant(req, user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responseData := restdto.MerchantCreateResponseDTO{
		ID:             merchant.ID,
		CreatedAt:      merchant.CreatedAt,
		UpdatedAt:      merchant.UpdatedAt,
		UserID:         merchant.UserID,
		Name:           merchant.Name,
		CommissionRate: merchant.CommissionRate,
	}
	c.JSON(http.StatusOK, responseData)
}

func (u *MerchantController) MerchantTokenList(c *gin.Context) {
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.GetMerchantByUserId(user.ID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
	}
	tokenList, err := u.MerchantUsecase.ListMerchantToken(merchant.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dtoList := make([]restdto.MerchantTokenCreateResponseDTO, 0)
	for _, token := range tokenList {
		dto := restdto.MerchantTokenCreateResponseDTO{
			ID:         token.ID,
			MerchantID: token.MerchantID,
			TokenID:    token.TokenID,
			Active:     true,
			CreatedAt:  token.CreatedAt,
		}
		dtoList = append(dtoList, dto)
	}

	c.JSON(http.StatusOK, dtoList)
}

func (u *MerchantController) MerchantTokenCreate(c *gin.Context) {
	var req restdto.MerchantTokenCreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.GetMerchantByUserId(user.ID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}
	token, err := u.MerchantUsecase.CreateMerchantToken(req, merchant.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	c.JSON(http.StatusOK, restdto.MerchantTokenCreateResponseDTO{
		ID:         token.ID,
		MerchantID: token.MerchantID,
		TokenID:    token.TokenID,
		Active:     true,
		CreatedAt:  token.CreatedAt,
	})
}
