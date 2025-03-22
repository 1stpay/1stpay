package usecase

import (
	"errors"
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/domain/service/kms"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/integration/restdto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentUsecase struct {
	PaymentRepo        repository.PaymentRepositoryInterface
	PaymentAddressRepo repository.PaymentAddressRepositoryInterface
	MerchantRepo       repository.MerchantRepositoryInterface
	DB                 *gorm.DB
}

type PaymentUsecaseInterface interface {
	CreatePaymentWithWallets(paymentData restdto.InvoiceCreateRestDTO, merchantId uuid.UUID) (model.Payment, error)
	GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error)
}

func NewPaymentUsecase(db *gorm.DB, paymentRepo repository.PaymentRepositoryInterface, paymentAddressRepo repository.PaymentAddressRepositoryInterface, merchantRepo repository.MerchantRepositoryInterface) *PaymentUsecase {
	return &PaymentUsecase{
		PaymentRepo:        paymentRepo,
		PaymentAddressRepo: paymentAddressRepo,
		MerchantRepo:       merchantRepo,
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

		paymentAddressList = append(paymentAddressList, model.PaymentAddress{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			PaymentID:  payment.ID,
			TokenID:    mt.Token.ID,
			PublicKey:  walletData.Address,
			PrivateKey: walletData.PrivateKey,
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
