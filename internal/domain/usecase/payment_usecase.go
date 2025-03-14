package usecase

import (
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
	CreatePaymentWithWallets(paymentData restdto.InvoiceCreateRestDTO) (model.Payment, error)
}

func (u *PaymentUsecase) CreatePaymentWithWallets(paymentData restdto.InvoiceCreateRestDTO, merchantId uuid.UUID) (model.Payment, error) {
	tx := u.DB.Begin()
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}

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
		tx.Rollback()
		return model.Payment{}, err
	}

	merchantTokens, err := u.MerchantRepo.ListMerchantToken(merchantId.String())
	if err != nil {
		tx.Rollback()
		return model.Payment{}, err
	}

	var paymentAddressList []model.PaymentAddress
	for _, mt := range merchantTokens {
		chainType := mt.Token.Blockchain.ChainType
		provider, err := kms.GetProvider(chainType)
		if err != nil {
			tx.Rollback()
			return model.Payment{}, err
		}

		walletData, err := provider.Create()
		if err != nil {
			tx.Rollback()
			return model.Payment{}, err
		}

		// Формируем запись PaymentAddress
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
		tx.Rollback()
		return model.Payment{}, err
	}
	if err := tx.Commit().Error; err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}
