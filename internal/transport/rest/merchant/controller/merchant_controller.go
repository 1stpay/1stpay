package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/helpers"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
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

// MerchantCreate godoc
// @Summary Create a new merchant
// @Description Creates a new merchant using provided data and returns the merchant details.
// @Tags Merchant
// @Accept json
// @Produce json
// @Param merchant body restdto.MerchantCreateRequestDTO true "Merchant creation payload"
// @Success 201 {object} restdto.MerchantCreateResponseDTO "Merchant created successfully"
// @Failure 400 {object} map[string]string "Invalid request or error during creation"
// @Router /merchant/api/v1/merchant/ [post]
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

// MerchantDetail godoc
// @Summary Get merchant details
// @Description Retrieves details of the merchant associated with the current user.
// @Tags Merchant
// @Produce json
// @Success 200 {object} restdto.MerchantDetailResponseDTO "Merchant details"
// @Failure 404 {object} map[string]string "Merchant not found"
// @Router /merchant/api/v1/merchant/me/ [get]
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

// MerchantUpdate godoc
// @Summary Update merchant information
// @Description Updates the merchant information for the current user.
// @Tags Merchant
// @Accept json
// @Produce json
// @Param merchant body restdto.MerchantCreateRequestDTO true "Merchant update payload"
// @Success 200 {object} restdto.MerchantCreateResponseDTO "Merchant updated successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /merchant/api/v1/merchant/me/ [put]
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

// MerchantTokenList godoc
// @Summary List merchant tokens
// @Description Retrieves a list of tokens associated with the current merchant.
// @Tags Merchant, Token
// @Produce json
// @Success 200 {array} restdto.MerchantTokenCreateResponseDTO "List of merchant tokens"
// @Failure 404 {object} map[string]string "Merchant not found"
// @Router /merchant/api/v1/merchant/me/tokens/ [get]
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

// MerchantTokenCreate godoc
// @Summary Create a merchant token
// @Description Creates a new token for the current merchant.
// @Tags Merchant, Token
// @Accept json
// @Produce json
// @Param token body restdto.MerchantTokenCreateRequestDTO true "Merchant token creation payload"
// @Success 200 {object} restdto.MerchantTokenCreateResponseDTO "Token created successfully"
// @Failure 400 {object} map[string]string "Invalid request or merchant not found"
// @Router /merchant/api/v1/merchant/me/tokens/ [post]
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

// MerchantAPIKeyCreate godoc
// @Summary Create a new API key for a merchant
// @Description Generates and stores a new API key for the current merchant. The raw API key is returned only once.
// @Tags Merchant, APIKey
// @Accept json
// @Produce json
// @Param apiKey body restdto.CreateAPIKeyRequestDTO true "API key creation payload"
// @Success 201 {object} restdto.CreateAPIKeyResponseDTO "API key created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /merchant/api/v1/merchant/me/api-key/ [post]
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

// MerchantAPIKeyList godoc
// @Summary List merchant API keys
// @Description Retrieves all API keys associated with the current merchant.
// @Tags Merchant, APIKey
// @Produce json
// @Success 200 {array} restdto.MerchantAPIKeyResponseDTO "List of API keys"
// @Failure 400 {object} map[string]string "Merchant not found"
// @Router /merchant/api/v1/merchant/me/api-key/ [get]
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

// MerchantAPIKeyDeactivate godoc
// @Summary Deactivate an API key
// @Description Deactivates a merchant's API key specified by its ID.
// @Tags Merchant, APIKey
// @Produce json
// @Param id path string true "API key ID"
// @Success 200 {object} map[string]string "API key deactivated successfully"
// @Failure 400 {object} map[string]string "Invalid API key ID"
// @Failure 500 {object} map[string]string "Error deactivating API key"
// @Router /merchant/api/v1/merchant/me/api-key/{id}/ [post]
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
