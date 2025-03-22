package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

type TokenRepository interface {
	ListActive() ([]model.Token, error)
	Create(Token model.Token) (model.Token, error)
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (r *tokenRepository) ListActive() ([]model.Token, error) {
	var TokenList []model.Token
	if err := r.db.Where("is_active = ?", true).Find(&TokenList).Error; err != nil {
		return []model.Token{}, err
	}
	return TokenList, nil
}

func (r *tokenRepository) Create(Token model.Token) (model.Token, error) {
	if err := r.db.Create(&Token).Error; err != nil {
		return model.Token{}, err
	}
	return Token, nil
}
