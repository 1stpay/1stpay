package usecase

import (
	"errors"
	"fmt"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/google/uuid"
)

type merchantUsecase struct {
	MerchantRepo repository.MerchantRepository
}

type MerchantUsecase interface {
	CreateMerchant(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error)
	UpdateMerchant(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error)
	GetMerchantByUserId(id string) (model.Merchant, error)
	CreateMerchantToken(merchantTokenData restdto.MerchantTokenCreateRequestDTO, merchantId string) (model.MerchantToken, error)
	ListMerchantToken(merchantId string) ([]model.MerchantToken, error)
}

func NewMerchantUsecase(merchantRepo repository.MerchantRepository) MerchantUsecase {
	return &merchantUsecase{
		MerchantRepo: merchantRepo,
	}
}

func (u *merchantUsecase) CreateMerchant(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return model.Merchant{}, fmt.Errorf("invalid user id format: %w", err)
	}

	existingMerchant, err := u.MerchantRepo.GetMerchantByUserId(userId)
	if err == nil {
		return existingMerchant, errors.New("user with this email already exists")
	}
	merchant := model.Merchant{
		UserID:         userUUID,
		Name:           merchantData.Name,
		CommissionRate: 0.5,
	}
	return u.MerchantRepo.CreateMerchant(merchant)
}
func (u *merchantUsecase) UpdateMerchant(merchantData restdto.MerchantCreateRequestDTO, userId string) (model.Merchant, error) {
	existingMerchant, err := u.MerchantRepo.GetMerchantByUserId(userId)
	existingMerchant.Name = merchantData.Name
	if err != nil {
		return existingMerchant, errors.New("user with this email already exists")
	}

	return u.MerchantRepo.UpdateMerchant(existingMerchant)
}
func (u *merchantUsecase) GetMerchantByUserId(userId string) (model.Merchant, error) {
	existingMerchant, err := u.MerchantRepo.GetMerchantByUserId(userId)
	if err != nil {
		return model.Merchant{}, errors.New("merchant not found")
	}
	return existingMerchant, nil
}

func (u *merchantUsecase) CreateMerchantToken(merchantTokenData restdto.MerchantTokenCreateRequestDTO, merchantId string) (model.MerchantToken, error) {
	merchantUUID, err := uuid.Parse(merchantId)
	if err != nil {
		return model.MerchantToken{}, fmt.Errorf("invalid user id format: %w", err)
	}
	existingToken, err := u.MerchantRepo.ListMerchantToken(merchantId, repository.MerchantTokenWithTokenId(merchantTokenData.TokenID.String()))
	if err != nil {
		return model.MerchantToken{}, fmt.Errorf("error checking existing merchant blockchain: %w", err)
	}
	if len(existingToken) > 0 {
		return model.MerchantToken{}, fmt.Errorf("merchant blockchain with BlockchainID %s already exists", merchantTokenData.TokenID)
	}
	merchantToken := model.MerchantToken{
		MerchantID: merchantUUID,
		TokenID:    merchantTokenData.TokenID,
		IsActive:   merchantTokenData.Active,
	}
	merchantToken, err = u.MerchantRepo.CreateMerchantToken(merchantToken)
	if err != nil {
		return model.MerchantToken{}, fmt.Errorf("error while merchant blockchain create: %w", err)
	}
	return merchantToken, err
}
func (u *merchantUsecase) ListMerchantToken(merchantId string) ([]model.MerchantToken, error) {
	objectList, err := u.MerchantRepo.ListMerchantToken(merchantId)
	if err != nil {
		return []model.MerchantToken{}, err
	}
	return objectList, err
}
