package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentAddressRepository struct {
	db *gorm.DB
}

type PaymentAddressRepository interface {
	Create(paymentAddress model.PaymentAddress) (model.PaymentAddress, error)
	CreateTx(tx *gorm.DB, paymentAddress model.PaymentAddress) (model.PaymentAddress, error)
	BulkCreateTx(tx *gorm.DB, paymentAddressList []model.PaymentAddress) ([]model.PaymentAddress, error)
	ListByPaymentId(paymentId uuid.UUID) ([]model.PaymentAddress, error)
}

func NewPaymentAddressRepository(db *gorm.DB) PaymentAddressRepository {
	return &paymentAddressRepository{
		db: db,
	}
}

func (r *paymentAddressRepository) Create(paymentAddress model.PaymentAddress) (model.PaymentAddress, error) {
	if err := r.db.Create(&paymentAddress).Error; err != nil {
		return model.PaymentAddress{}, err
	}
	return paymentAddress, nil
}

func (r *paymentAddressRepository) CreateTx(tx *gorm.DB, paymentAddress model.PaymentAddress) (model.PaymentAddress, error) {
	if err := tx.Create(&paymentAddress).Error; err != nil {
		return model.PaymentAddress{}, err
	}
	return paymentAddress, nil
}

func (r *paymentAddressRepository) BulkCreateTx(tx *gorm.DB, paymentAddressList []model.PaymentAddress) ([]model.PaymentAddress, error) {
	if err := tx.Create(&paymentAddressList).Error; err != nil {
		return []model.PaymentAddress{}, err
	}
	return paymentAddressList, nil
}

func (r *paymentAddressRepository) ListByPaymentId(paymentId uuid.UUID) ([]model.PaymentAddress, error) {
	var addresses []model.PaymentAddress
	if err := r.db.
		Where("payment_id = ?", paymentId).
		Preload("Token").
		Preload("Token.Blockchain").
		Find(&addresses).Error; err != nil {
		return []model.PaymentAddress{}, err
	}
	return addresses, nil
}
