package repository

import "gorm.io/gorm"

type MerchantTokenOption func(db *gorm.DB) *gorm.DB

func MerchantTokenWithMerchantId(merchantId string) MerchantTokenOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("merchant_id = ?", merchantId)
	}
}
func MerchantTokenWithTokenId(tokenId string) MerchantTokenOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("token_id = ?", tokenId)
	}
}
