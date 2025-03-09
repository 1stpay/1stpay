package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type MerchantRepository struct {
	db *gorm.DB
}

type MerchantRepositoryInterface interface {
	Create(merchant model.Merchant) (model.Merchant, error)
	GetById(id string) (model.Merchant, error)
	GetByUserId(userId string) (model.Merchant, error)
}

func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(merchant model.Merchant) (model.Merchant, error) {
	if err := r.db.Create(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) GetById(id string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("id = ?", id).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) GetByUserId(userId string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("user_id = ?", userId).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}
