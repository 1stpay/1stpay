package usecase

import (
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/integration/restdto"
	"github.com/google/uuid"
)

type PaymentUsecase struct {
	PaymentRepo repository.PaymentRepositoryInterface
}

type PaymentUsecaseInterface interface {
	ListActive() ([]model.Payment, error)
	Create(payment restdto.InvoiceCreateRestDTO) (model.Payment, error)
}

func NewPaymentUsecase(repo repository.PaymentRepositoryInterface) *PaymentUsecase {
	return &PaymentUsecase{
		PaymentRepo: repo,
	}
}

func (u *PaymentUsecase) Create(paymentData restdto.InvoiceCreateRestDTO) (model.Payment, error) {
	now := time.Now()
	AMLStatus := enum.PaymentAMLStatusPending
	paymentStatus := enum.PaymentStatusPending
	payment := model.Payment{
		ID:               uuid.New(),
		MerchantID:       uuid.New(),
		RequestedAmount:  paymentData.RequestedAmount,
		PaidAmount:       0,
		CommissionAmount: 0,
		ExpiresAt:        &now,
		AMLStatus:        &AMLStatus,
		Status:           paymentStatus,
		InvoiceEmail:     paymentData.Email,
	}
	payment, err := u.PaymentRepo.Create(payment)
	if err != nil {
		return model.Payment{}, err
	}
	return payment, err
}
