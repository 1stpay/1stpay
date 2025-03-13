package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

type TokenRepositoryInterface interface {
	ListActive() ([]model.Token, error)
	Create(Token model.Token) (model.Token, error)
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) ListActive() ([]model.Token, error) {
	var TokenList []model.Token
	if err := r.db.Where("is_active = ?", true).Find(&TokenList).Error; err != nil {
		return []model.Token{}, err
	}
	return TokenList, nil
}

func (r *TokenRepository) Create(Token model.Token) (model.Token, error) {
	if err := r.db.Create(&Token).Error; err != nil {
		return model.Token{}, err
	}
	return Token, nil
}
