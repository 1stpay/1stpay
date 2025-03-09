package usecase

import (
	"errors"
	"fmt"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/google/uuid"
)

type MerchantUsecase struct {
	MerchantRepo repository.MerchantRepositoryInterface
}

type MerchantUsecaseInterface interface {
	Create(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error)
}

func NewMerchantUsecase(merchantRepo repository.MerchantRepositoryInterface) *MerchantUsecase {
	return &MerchantUsecase{
		MerchantRepo: merchantRepo,
	}
}

func (u *MerchantUsecase) Create(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return model.Merchant{}, fmt.Errorf("invalid user id format: %w", err)
	}

	existingMerchant, err := u.MerchantRepo.GetByUserId(userId)
	if err == nil {
		return existingMerchant, errors.New("user with this email already exists")
	}
	merchant := model.Merchant{
		UserID:         userUUID,
		Name:           merchantData.Name,
		CommissionRate: 0.5,
	}
	return u.MerchantRepo.Create(merchant)
}
