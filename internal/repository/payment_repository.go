package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type PaymentRepositoryInterface interface {
	Create(payment model.Payment) (model.Payment, error)
	CreateTx(tx *gorm.DB, payment model.Payment) (model.Payment, error)
	GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error)
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Create(payment model.Payment) (model.Payment, error) {
	if err := r.db.Create(&payment).Error; err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *PaymentRepository) CreateTx(tx *gorm.DB, payment model.Payment) (model.Payment, error) {
	if err := tx.Create(&payment).Error; err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *PaymentRepository) GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error) {
	var payment model.Payment
	var paymentAddressList []model.PaymentAddress

	if err := r.db.Where("id = ?", paymentID).Preload("Merchant").First(&payment).Error; err != nil {
		return model.Payment{}, []model.PaymentAddress{}, err
	}

	if err := r.db.Where("payment_id = ?", paymentID).Preload("Token").Preload("Token.Blockchain").Find(&paymentAddressList).Error; err != nil {
		return model.Payment{}, []model.PaymentAddress{}, err
	}

	return payment, paymentAddressList, nil
}
