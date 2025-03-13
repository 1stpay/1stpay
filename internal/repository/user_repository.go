package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	Create(user model.User) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetById(id string) (model.User, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user model.User) (model.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetById(id string) (model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
