package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type PaymentAddressRepository struct {
	db *gorm.DB
}

type PaymentAddressRepositoryInterface interface {
	Create(paymentAddress model.PaymentAddress) (model.PaymentAddress, error)
	CreateTx(tx *gorm.DB, paymentAddress model.PaymentAddress) (model.PaymentAddress, error)
}

func NewPaymentAddressRepository(db *gorm.DB) *PaymentAddressRepository {
	return &PaymentAddressRepository{
		db: db,
	}
}

func (r *PaymentAddressRepository) Create(paymentAddress model.PaymentAddress) (model.PaymentAddress, error) {
	if err := r.db.Create(&paymentAddress).Error; err != nil {
		return model.PaymentAddress{}, err
	}
	return paymentAddress, nil
}

func (r *PaymentAddressRepository) CreateTx(tx *gorm.DB, paymentAddress model.PaymentAddress) (model.PaymentAddress, error) {
	if err := tx.Create(&paymentAddress).Error; err != nil {
		return model.PaymentAddress{}, err
	}
	return paymentAddress, nil
}
