package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type MerchantRepository struct {
	db *gorm.DB
}

type MerchantRepositoryInterface interface {
	CreateMerchant(merchant model.Merchant) (model.Merchant, error)
	GetMerchantById(id string) (model.Merchant, error)
	UpdateMerchant(merchant model.Merchant) (model.Merchant, error)
	GetMerchantByUserId(userId string) (model.Merchant, error)
	CreateMerchantToken(merchantToken model.MerchantToken) (model.MerchantToken, error)
	ListMerchantToken(merchantId string, opts ...MerchantTokenOption) ([]model.MerchantToken, error)
}

func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) CreateMerchant(merchant model.Merchant) (model.Merchant, error) {
	if err := r.db.Create(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) GetMerchantById(id string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("id = ?", id).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) UpdateMerchant(merchant model.Merchant) (model.Merchant, error) {
	if err := r.db.Save(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) GetMerchantByUserId(userId string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("user_id = ?", userId).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *MerchantRepository) CreateMerchantToken(merchantToken model.MerchantToken) (model.MerchantToken, error) {
	if err := r.db.Create(&merchantToken).Error; err != nil {
		return model.MerchantToken{}, err
	}
	return merchantToken, nil
}

func (r *MerchantRepository) ListMerchantToken(merchantId string, opts ...MerchantTokenOption) ([]model.MerchantToken, error) {
	var tokenList []model.MerchantToken
	dbQuery := r.db.Where("merchant_id = ?", merchantId).
		Preload("Token").
		Preload("Token.Blockchain")
	for _, opt := range opts {
		dbQuery = opt(dbQuery)
	}
	if err := dbQuery.Find(&tokenList).Error; err != nil {
		return []model.MerchantToken{}, err
	}
	return tokenList, nil
}
