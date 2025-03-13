package factory

import (
	"fmt"

	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/model"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestFactory struct {
	db       *gorm.DB
	Usecases *config.Usecases
	Repos    *config.Repos
}

func NewTestFactory(db *gorm.DB, deps *config.Dependencies) *TestFactory {
	return &TestFactory{
		db:       db,
		Usecases: deps.Usecases,
		Repos:    deps.Repos,
	}
}

func (f *TestFactory) CreateUser() (model.User, string) {
	uniqueEmail := fmt.Sprintf("testuser_%s@example.com", uuid.New().String())
	registerData := restdto.RegisterRequest{
		Email:    uniqueEmail,
		Password: "Secret",
	}
	user, accessToken, err := f.Usecases.AuthUsecase.Register(registerData)
	if err != nil {
		panic("Error while test user creation")
	}
	return user, accessToken
}

func (f *TestFactory) CreateMerchant(userId string) model.Merchant {
	merchantData := restdto.MerchantCreateRequestDTO{
		Name: "Test",
	}
	user, err := f.Usecases.MerchantUsecase.CreateMerchant(merchantData, userId)
	if err != nil {
		panic("Error while test user creation")
	}
	return user
}

func (f *TestFactory) CreateBlockchainList() []model.Blockchain {
	blockachainList := []model.Blockchain{
		{ID: uuid.New(), Name: "Ethereum", IsActive: true, ChainType: enum.EVM},
		{ID: uuid.New(), Name: "Solana", IsActive: true, ChainType: enum.SOLANA},
		{ID: uuid.New(), Name: "Ton", IsActive: true, ChainType: enum.TON},
		{ID: uuid.New(), Name: "Tron", IsActive: true, ChainType: enum.TRON},
	}
	for _, blockchain := range blockachainList {
		_, err := f.Usecases.BlockchainUsecase.Create(blockchain)
		if err != nil {
			panic("Error while test blockchain creation")
		}
	}
	return blockachainList
}

func (f *TestFactory) CreateTokenList(blockchainList []model.Blockchain) []model.Token {
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
		_, err := f.Repos.TokenRepo.Create(token)
		if err != nil {
			panic("Error while test blockchain creation")
		}
	}
	return tokenList
}
