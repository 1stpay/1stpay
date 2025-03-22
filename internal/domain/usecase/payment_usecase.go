package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/domain/service/kms"
	"github.com/1stpay/1stpay/internal/infrastructure/price_service"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/integration/restdto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	PaymentRepo        repository.PaymentRepositoryInterface
	PaymentAddressRepo repository.PaymentAddressRepository
	MerchantRepo       repository.MerchantRepository
	PriceService       price_service.PriceService
	DB                 *gorm.DB
}

type PaymentUsecaseInterface interface {
	CreatePaymentWithWallets(paymentData restdto.InvoiceCreateRestDTO, merchantId uuid.UUID) (model.Payment, error)
	GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error)
}

func NewPaymentUsecase(
	db *gorm.DB,
	paymentRepo repository.PaymentRepositoryInterface,
	paymentAddressRepo repository.PaymentAddressRepository,
	merchantRepo repository.MerchantRepository,
	priceService price_service.PriceService,
) *PaymentUsecase {
	return &PaymentUsecase{
		PaymentRepo:        paymentRepo,
		PaymentAddressRepo: paymentAddressRepo,
		MerchantRepo:       merchantRepo,
		PriceService:       priceService,
		DB:                 db,
	}
}

func (u *PaymentUsecase) CreatePaymentWithWallets(paymentData restdto.InvoiceCreateRestDTO, merchantId uuid.UUID) (model.Payment, error) {
	tx := u.DB.Begin()
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now()
	amlStatus := enum.PaymentAMLStatusPending
	paymentStatus := enum.PaymentStatusPending

	payment := model.Payment{
		ID:               uuid.New(),
		MerchantID:       merchantId,
		RequestedAmount:  paymentData.RequestedAmount,
		PaidAmount:       0,
		CommissionAmount: 0,
		ExpiresAt:        &now,
		AMLStatus:        &amlStatus,
		Status:           paymentStatus,
		InvoiceEmail:     paymentData.Email,
	}

	payment, err := u.PaymentRepo.CreateTx(tx, payment)
	if err != nil {
		return model.Payment{}, err
	}

	merchantTokens, err := u.MerchantRepo.ListMerchantToken(merchantId.String())
	if len(merchantTokens) == 0 {
		return model.Payment{}, errors.New("setup tokens you work with")
	}
	if err != nil {
		return model.Payment{}, err
	}

	var paymentAddressList []model.PaymentAddress
	for _, mt := range merchantTokens {
		chainType := mt.Token.Blockchain.ChainType
		provider, err := kms.GetProvider(chainType)
		if err != nil {
			return model.Payment{}, err
		}

		walletData, err := provider.Create()
		if err != nil {
			return model.Payment{}, err
		}
		var tokenCfg map[string]string
		if err := json.Unmarshal(mt.Token.Config, &tokenCfg); err != nil {
			return model.Payment{}, fmt.Errorf("failed to parse config for token %s: %w", mt.ID, err)
		}
		priceServiceKey, ok := tokenCfg["price_service_key"]
		if !ok || priceServiceKey == "" {
			return model.Payment{}, fmt.Errorf("failed to query price for token %s: %w", mt.ID, err)
		}
		assetPrice, err := u.PriceService.GetPrice(priceServiceKey)
		if err != nil {
			return model.Payment{}, fmt.Errorf("failed to query price for token %s: %w", mt.ID, err)
		}
		requestedAmountForAsset := paymentData.RequestedAmount / assetPrice
		factor := math.Pow10(mt.Token.Decimals)
		bf := new(big.Float).SetFloat64(requestedAmountForAsset * factor)
		requestedAmountForAssetWei, _ := bf.Int64()

		paymentAddressList = append(paymentAddressList, model.PaymentAddress{
			ID:                 uuid.New(),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			PaymentID:          payment.ID,
			TokenID:            mt.Token.ID,
			PublicKey:          walletData.Address,
			PrivateKey:         walletData.PrivateKey,
			RequestedAmount:    requestedAmountForAsset,
			RequestedAmountWei: int(requestedAmountForAssetWei),
		})

	}
	_, err = u.PaymentAddressRepo.BulkCreateTx(tx, paymentAddressList)
	if err != nil {
		return model.Payment{}, err
	}
	if err := tx.Commit().Error; err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (u *PaymentUsecase) GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error) {
	payment, paypaymentAddressList, err := u.PaymentRepo.GetPaymentWithAddresses(paymentID)
	if err != nil {
		return model.Payment{}, []model.PaymentAddress{}, err
	}
	return payment, paypaymentAddressList, nil
}
