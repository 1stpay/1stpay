package usecase

import (
	"time"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/google/uuid"
)

type PaymentAddressUsecaseInterface interface {
	GenerateAddresses(payment model.Payment, tokens []model.MerchantToken) ([]model.PaymentAddress, error)
}

type PaymentAddressUsecase struct {
	PaymentAddressRepo repository.PaymentAddressRepositoryInterface
}

func NewPaymentAddressUsecase(repo repository.PaymentAddressRepositoryInterface) *PaymentAddressUsecase {
	return &PaymentAddressUsecase{
		PaymentAddressRepo: repo,
	}
}

func (s *PaymentAddressUsecase) GenerateAddresses(payment model.Payment, tokens []model.MerchantToken) ([]model.PaymentAddress, error) {
	var addresses []model.PaymentAddress
	now := time.Now()

	for _, merchantToken := range tokens {
		publicKey := "dummy_public_key_" + merchantToken.TokenID.String()
		privateKey := "dummy_private_key_" + merchantToken.TokenID.String()

		paymentAddress := model.PaymentAddress{
			ID:                 uuid.New(),
			CreatedAt:          now,
			UpdatedAt:          now,
			PaymentID:          payment.ID,
			TokenID:            merchantToken.TokenID,
			PublicKey:          publicKey,
			PrivateKey:         privateKey,
			RequestedAmount:    payment.RequestedAmount,
			PaidAmount:         0,
			RequestedAmountWei: 0,
			PaidAmountWei:      0,
		}

		savedAddress, err := s.PaymentAddressRepo.Create(paymentAddress)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, savedAddress)
	}

	return addresses, nil
}
