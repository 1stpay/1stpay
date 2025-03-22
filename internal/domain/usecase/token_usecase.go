package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type TokenUsecase struct {
	TokenRepo repository.TokenRepository
}

type TokenUsecaseInterface interface {
	ListActive() ([]model.Token, error)
	Create(Token model.Token) (model.Token, error)
}

func NewTokenUsecase(repo repository.TokenRepository) *TokenUsecase {
	return &TokenUsecase{
		TokenRepo: repo,
	}
}

func (u *TokenUsecase) ListActive() ([]model.Token, error) {
	tokenList, err := u.TokenRepo.ListActive()
	if err != nil {
		return []model.Token{}, err
	}
	return tokenList, err
}

func (u *TokenUsecase) Create(Token model.Token) (model.Token, error) {
	Token, err := u.TokenRepo.Create(Token)
	if err != nil {
		return model.Token{}, err
	}
	return Token, err
}
