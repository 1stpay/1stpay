package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/helpers"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/merchant/rest_dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchantController struct {
	MerchantUsecase       usecase.MerchantUsecase
	UserUsecase           usecase.UserUsecase
	MerchantAPIKeyUsecase usecase.MerchantAPIKeyUsecase
}

func NewMerchantController(merchantUsecase usecase.MerchantUsecase, merchantAPIKeyUsecase usecase.MerchantAPIKeyUsecase, userUsecase usecase.UserUsecase) *MerchantController {
	return &MerchantController{
		MerchantUsecase:       merchantUsecase,
		UserUsecase:           userUsecase,
		MerchantAPIKeyUsecase: merchantAPIKeyUsecase,
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

func (u *MerchantController) MerchantAPIKeyCreate(c *gin.Context) {
	var req restdto.CreateAPIKeyRequestDTO
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while merchant check"})
		return
	}

	createdKey, rawKey, err := u.MerchantAPIKeyUsecase.CreateAPIKey(merchant.ID, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := restdto.CreateAPIKeyResponseDTO{
		ID:        createdKey.ID,
		Name:      createdKey.Name,
		APIKey:    rawKey,
		CreatedAt: createdKey.CreatedAt,
		ExpiresAt: createdKey.ExpiresAt,
		IsActive:  createdKey.IsActive,
	}

	c.JSON(http.StatusCreated, response)
}

func (u *MerchantController) MerchantAPIKeyList(c *gin.Context) {
	user, ok := helpers.GetUserOrAbort(c, u.UserUsecase)
	if !ok {
		return
	}
	merchant, err := u.MerchantUsecase.GetMerchantByUserId(user.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while merchant check"})
		return
	}

	keys, err := u.MerchantAPIKeyUsecase.ListAPIKeys(merchant.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dtos []restdto.MerchantAPIKeyResponseDTO
	for _, key := range keys {
		dto := restdto.MerchantAPIKeyResponseDTO{
			ID:        key.ID,
			Name:      key.Name,
			CreatedAt: key.CreatedAt,
			ExpiresAt: key.ExpiresAt,
			IsActive:  key.IsActive,
		}
		dtos = append(dtos, dto)
	}

	c.JSON(http.StatusOK, dtos)
}

func (u *MerchantController) MerchantAPIKeyDeactivate(c *gin.Context) {
	keyIDStr := c.Param("id")
	keyID, err := uuid.Parse(keyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	if err := u.MerchantAPIKeyUsecase.DeactivateAPIKey(keyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deactivated successfully"})
}
