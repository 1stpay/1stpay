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
	if err := r.db.Create(&payment).Error; err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}
