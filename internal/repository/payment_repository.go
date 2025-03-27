package repository

import (
	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

type PaymentRepository interface {
	Create(payment model.Payment) (model.Payment, error)
	CreateTx(tx *gorm.DB, payment model.Payment) (model.Payment, error)
	GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error)
	GetPaymentListWithStatus(status enum.PaymentStatus) ([]model.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(payment model.Payment) (model.Payment, error) {
	if err := r.db.Create(&payment).Error; err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *paymentRepository) CreateTx(tx *gorm.DB, payment model.Payment) (model.Payment, error) {
	if err := tx.Create(&payment).Error; err != nil {
		return model.Payment{}, err
	}
	return payment, nil
}

func (r *paymentRepository) GetPaymentWithAddresses(paymentID string) (model.Payment, []model.PaymentAddress, error) {
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

func (r *paymentRepository) GetPaymentListWithStatus(status enum.PaymentStatus) ([]model.Payment, error) {
	var paymentList []model.Payment
	if err := r.db.Where("status = ?", status).Find(&paymentList).Error; err != nil {
		return []model.Payment{}, err
	}

	return paymentList, nil
}
