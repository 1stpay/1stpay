package factory

import (
	"encoding/json"
	"fmt"

	"github.com/1stpay/1stpay/internal/config"
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/model"
	commonRestDto "github.com/1stpay/1stpay/internal/transport/rest/common/restdto"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestFactory struct {
	db       *gorm.DB
	Usecases *config.Usecases
	Repos    *config.Repos
	Deps     *config.Dependencies
}

func NewTestFactory(db *gorm.DB, deps *config.Dependencies) *TestFactory {
	return &TestFactory{
		db:       db,
		Usecases: deps.Usecases,
		Repos:    deps.Repos,
		Deps:     deps,
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
	configMap := map[string]string{"price_service_key": "USDT"}
	configBytes, err := json.Marshal(configMap)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling config: %v", err))
	}
	var tokenList []model.Token
	for _, blockchain := range blockchainList {
		symbol := fmt.Sprintf("USDT-%s", blockchain.Name)
		obj := model.Token{
			ID:              uuid.New(),
			Name:            fmt.Sprintf("USDT-%s", blockchain.Name),
			Symbol:          symbol,
			BlockchainID:    blockchain.ID,
			ContractAddress: "0x00",
			Decimals:        18,
			IsNative:        false,
			IsActive:        true,
			Config:          configBytes,
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

func (f *TestFactory) CreateMerchantTokenList(merchant model.Merchant, tokenList []model.Token) []model.MerchantToken {
	var merchantTokenList []model.MerchantToken
	for _, token := range tokenList {
		obj := model.MerchantToken{
			ID:         uuid.New(),
			MerchantID: merchant.ID,
			TokenID:    token.ID,
			Balance:    0,
			IsActive:   true,
		}
		merchantTokenList = append(merchantTokenList, obj)
	}

	for _, token := range merchantTokenList {
		_, err := f.Repos.MerchantRepo.CreateMerchantToken(token)
		if err != nil {
			panic(fmt.Sprintf("Error while creating token: %v", err))
		}
	}
	return merchantTokenList
}

func (f *TestFactory) CreatePayment(merchant model.Merchant) model.Payment {
	payment, err := f.Usecases.PaymentUsecase.CreatePaymentWithWallets(
		commonRestDto.PaymentCreateRestDTO{RequestedAmount: 100},
		merchant.ID,
	)
	if err != nil {
		panic("Error while test blockchain creation")
	}
	return payment
}
