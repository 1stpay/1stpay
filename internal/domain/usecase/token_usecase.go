package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type tokenUsecase struct {
	TokenRepo repository.TokenRepository
}

type TokenUsecase interface {
	ListActive() ([]model.Token, error)
	Create(Token model.Token) (model.Token, error)
}

func NewTokenUsecase(repo repository.TokenRepository) TokenUsecase {
	return &tokenUsecase{
		TokenRepo: repo,
	}
}

func (u *tokenUsecase) ListActive() ([]model.Token, error) {
	tokenList, err := u.TokenRepo.ListActive()
	if err != nil {
		return []model.Token{}, err
	}
	return tokenList, err
}

func (u *tokenUsecase) Create(Token model.Token) (model.Token, error) {
	Token, err := u.TokenRepo.Create(Token)
	if err != nil {
		return model.Token{}, err
	}
	return Token, err
}
