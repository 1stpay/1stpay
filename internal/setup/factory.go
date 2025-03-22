package setup

import (
	"fmt"

	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setup struct {
	db       *gorm.DB
	Usecases *config.Usecases
	Repos    *config.Repos
}

func NewSetup(db *gorm.DB, deps *config.Dependencies) *Setup {
	return &Setup{
		db:       db,
		Usecases: deps.Usecases,
		Repos:    deps.Repos,
	}
}

func (s *Setup) initBlockchain() []model.Blockchain {
	blockachainList := []model.Blockchain{
		{ID: uuid.New(), Name: "Ethereum", IsActive: true, ChainType: enum.EVM},
		{ID: uuid.New(), Name: "Solana", IsActive: true, ChainType: enum.SOLANA},
		{ID: uuid.New(), Name: "Ton", IsActive: true, ChainType: enum.TON},
		{ID: uuid.New(), Name: "Tron", IsActive: true, ChainType: enum.TRON},
	}
	for _, blockchain := range blockachainList {
		_, err := s.Usecases.BlockchainUsecase.Create(blockchain)
		if err != nil {
			panic("Error while test blockchain creation")
		}
	}
	return blockachainList
}
func (s *Setup) initToken(blockchainList []model.Blockchain) []model.Token {
	var tokenList []model.Token
	for _, blockchain := range blockchainList {
		symbol := fmt.Sprintf("USDT-%s", blockchain.Name)
		obj := model.Token{
			ID:           uuid.New(),
			Name:         fmt.Sprintf("USDT-%s", blockchain.Name),
			Symbol:       symbol,
			BlockchainID: blockchain.ID,
			IsNative:     false,
			IsActive:     true,
		}
		tokenList = append(tokenList, obj)
	}

	for _, token := range tokenList {
		_, err := s.Repos.TokenRepo.Create(token)
		if err != nil {
			panic("Error while test blockchain creation")
		}
	}
	return tokenList
}

func (s *Setup) Init() {
	blockchainList := s.initBlockchain()
	s.initToken(blockchainList)
}
