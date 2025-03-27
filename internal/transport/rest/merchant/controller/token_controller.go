package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v6"
)

type TokenController struct {
	TokenUsecase usecase.TokenUsecase
}

func NewTokenController(usecase usecase.TokenUsecase) *TokenController {
	return &TokenController{
		TokenUsecase: usecase,
	}
}

// ListActive godoc
// @Summary List active tokens
// @Description Retrieves a list of all active tokens.
// @Tags Token
// @Produce json
// @Success 200 {array} restdto.TokenListResponseDTO "List of active tokens"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /merchant/api/v1/token/ [get]
func (con *TokenController) ListActive(c *gin.Context) {
	tokenList, err := con.TokenUsecase.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dtoList := make([]restdto.TokenListResponseDTO, 0)
	for _, token := range tokenList {
		dto := restdto.TokenListResponseDTO{
			IsNative:     token.IsNative,
			ID:           token.ID,
			Name:         token.Name,
			Symbol:       token.Symbol,
			BlockchainID: token.BlockchainID,
			Logo:         null.StringFromPtr(token.Logo),
			IsActive:     token.IsActive,
		}
		dtoList = append(dtoList, dto)
	}

	c.JSON(http.StatusOK, dtoList)
}
