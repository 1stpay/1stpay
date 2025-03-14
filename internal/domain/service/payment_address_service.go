package service

import (
	"time"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/google/uuid"
)

// PaymentAddressServiceInterface описывает интерфейс сервиса по генерации адресов платежа
type PaymentAddressServiceInterface interface {
	GenerateAddresses(payment model.Payment, tokens []model.MerchantToken) ([]model.PaymentAddress, error)
}

// PaymentAddressService реализует логику генерации кошельков для платежа
type PaymentAddressService struct {
	PaymentAddressRepo repository.PaymentAddressRepositoryInterface
}

// NewPaymentAddressService создаёт новый экземпляр PaymentAddressService
func NewPaymentAddressService(repo repository.PaymentAddressRepositoryInterface) *PaymentAddressService {
	return &PaymentAddressService{
		PaymentAddressRepo: repo,
	}
}

// GenerateAddresses генерирует для платежа кошельки для каждого активного токена мерчанта.
// Здесь генерация ключей упрощена до демо-логики, но в реальном случае можно интегрироваться с блокчейн-сервисом.
func (s *PaymentAddressService) GenerateAddresses(payment model.Payment, tokens []model.MerchantToken) ([]model.PaymentAddress, error) {
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

		// Сохраняем адрес через репозиторий
		savedAddress, err := s.PaymentAddressRepo.Create(paymentAddress)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, savedAddress)
	}

	return addresses, nil
}
