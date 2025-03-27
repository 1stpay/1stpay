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
	"github.com/1stpay/1stpay/internal/transport/rest/common/restdto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentUsecase struct {
	PaymentRepo        repository.PaymentRepository
	PaymentAddressRepo repository.PaymentAddressRepository
	MerchantRepo       repository.MerchantRepository
	PriceService       price_service.PriceService
	DB                 *gorm.DB
}

type PaymentUsecase interface {
	CreatePaymentWithWallets(paymentData restdto.PaymentCreateRestDTO, merchantId uuid.UUID) (model.Payment, error)
	GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error)
	ConfirmInvoice(inv model.Payment, addr model.PaymentAddress, balance *big.Int) error
}

func NewPaymentUsecase(
	db *gorm.DB,
	paymentRepo repository.PaymentRepository,
	paymentAddressRepo repository.PaymentAddressRepository,
	merchantRepo repository.MerchantRepository,
	priceService price_service.PriceService,
) PaymentUsecase {
	return &paymentUsecase{
		PaymentRepo:        paymentRepo,
		PaymentAddressRepo: paymentAddressRepo,
		MerchantRepo:       merchantRepo,
		PriceService:       priceService,
		DB:                 db,
	}
}

func (u *paymentUsecase) CreatePaymentWithWallets(paymentData restdto.PaymentCreateRestDTO, merchantId uuid.UUID) (model.Payment, error) {
	tx := u.DB.Begin()
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().Add(time.Duration(time.Hour))
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

func (u *paymentUsecase) GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error) {
	payment, paypaymentAddressList, err := u.PaymentRepo.GetPaymentWithAddresses(paymentID)
	if err != nil {
		return model.Payment{}, []model.PaymentAddress{}, err
	}
	return payment, paypaymentAddressList, nil
}

func (u *paymentUsecase) ConfirmInvoice(inv model.Payment, addr model.PaymentAddress, balance *big.Int) error {
	decimals := addr.Token.Decimals
	factor := math.Pow10(decimals)

	fBalance := new(big.Float).SetInt(balance)
	paidAmountFloat, _ := new(big.Float).Quo(fBalance, big.NewFloat(factor)).Float64()

	return u.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.PaymentAddress{}).
			Where("id = ?", addr.ID).
			Updates(map[string]interface{}{
				"paid_amount":     paidAmountFloat,
				"paid_amount_wei": balance.Int64(),
			}).Error; err != nil {
			return fmt.Errorf("error updating payment address: %w", err)
		}

		if err := tx.Model(&model.Payment{}).
			Where("id = ?", inv.ID).
			Updates(map[string]interface{}{
				"status":        enum.PaymentStatusCompleted,
				"used_token_id": addr.Token.ID,
			}).Error; err != nil {
			return fmt.Errorf("error updating payment: %w", err)
		}
		if err := tx.Model(&model.MerchantToken{}).
			Where("merchant_id = ? AND token_id = ?", inv.MerchantID, addr.Token.ID).
			Update("balance", gorm.Expr("balance + ?", paidAmountFloat)).Error; err != nil {
			return fmt.Errorf("error updating merchant token balance: %w", err)
		}
		return nil
	})
}
