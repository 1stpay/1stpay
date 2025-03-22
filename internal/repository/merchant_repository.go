package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type merchantRepository struct {
	db *gorm.DB
}

type MerchantRepository interface {
	CreateMerchant(merchant model.Merchant) (model.Merchant, error)
	GetMerchantById(id string) (model.Merchant, error)
	UpdateMerchant(merchant model.Merchant) (model.Merchant, error)
	GetMerchantByUserId(userId string) (model.Merchant, error)
	CreateMerchantToken(merchantToken model.MerchantToken) (model.MerchantToken, error)
	ListMerchantToken(merchantId string, opts ...MerchantTokenOption) ([]model.MerchantToken, error)
	CreateMerchantAPIKey(apiKey model.MerchantAPIKey) (model.MerchantAPIKey, error)
	ListMerchantAPIKey(merchantId string) ([]model.MerchantAPIKey, error)
	DeactivateMerchantAPIKey(tokenID string) error
	GetMerchantAPIKeyByHash(apiKeyHash string) (model.MerchantAPIKey, error)
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r merchantRepository) CreateMerchant(merchant model.Merchant) (model.Merchant, error) {
	if err := r.db.Create(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *merchantRepository) GetMerchantById(id string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("id = ?", id).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *merchantRepository) UpdateMerchant(merchant model.Merchant) (model.Merchant, error) {
	if err := r.db.Save(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *merchantRepository) GetMerchantByUserId(userId string) (model.Merchant, error) {
	var merchant model.Merchant
	if err := r.db.Where("user_id = ?", userId).First(&merchant).Error; err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}

func (r *merchantRepository) CreateMerchantToken(merchantToken model.MerchantToken) (model.MerchantToken, error) {
	if err := r.db.Create(&merchantToken).Error; err != nil {
		return model.MerchantToken{}, err
	}
	return merchantToken, nil
}

func (r *merchantRepository) ListMerchantToken(merchantId string, opts ...MerchantTokenOption) ([]model.MerchantToken, error) {
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

func (r *merchantRepository) CreateMerchantAPIKey(apiKey model.MerchantAPIKey) (model.MerchantAPIKey, error) {
	if err := r.db.Create(&apiKey).Error; err != nil {
		return model.MerchantAPIKey{}, err
	}
	return apiKey, nil
}

func (r *merchantRepository) ListMerchantAPIKey(merchantId string) ([]model.MerchantAPIKey, error) {
	var keyList []model.MerchantAPIKey
	dbQuery := r.db.Where("merchant_id = ?", merchantId).Where("is_active = ?", true)
	if err := dbQuery.Find(&keyList).Error; err != nil {
		return []model.MerchantAPIKey{}, err
	}
	return keyList, nil
}

func (r *merchantRepository) DeactivateMerchantAPIKey(apiKeyID string) error {
	return r.db.Model(&model.MerchantAPIKey{}).
		Where("id = ?", apiKeyID).
		Update("is_active", false).Error
}

func (r *merchantRepository) GetMerchantAPIKeyByHash(apiKeyHash string) (model.MerchantAPIKey, error) {
	var key model.MerchantAPIKey
	if err := r.db.Where("api_key = ?", apiKeyHash).First(&key).Error; err != nil {
		return model.MerchantAPIKey{}, err
	}
	return key, nil
}
