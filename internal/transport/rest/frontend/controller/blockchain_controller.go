package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v6"
)

type BlockchainController struct {
	BlockchainUsecase usecase.BlockchainUsecase
}

func NewBlockchainController(usecase usecase.BlockchainUsecase) *BlockchainController {
	return &BlockchainController{
		BlockchainUsecase: usecase,
	}
}

func (con *BlockchainController) ListActive(c *gin.Context) {
	blockchainList, err := con.BlockchainUsecase.ListActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dtoList := make([]restdto.BlockchainListResponseDTO, 0)
	for _, blockchain := range blockchainList {
		dto := restdto.BlockchainListResponseDTO{
			ID:        blockchain.ID,
			Name:      blockchain.Name,
			Logo:      null.StringFromPtr(blockchain.Logo),
			IsActive:  blockchain.IsActive,
			ChainType: blockchain.ChainType,
		}
		dtoList = append(dtoList, dto)
	}

	c.JSON(http.StatusOK, dtoList)
}
