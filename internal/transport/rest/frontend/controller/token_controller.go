package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
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
