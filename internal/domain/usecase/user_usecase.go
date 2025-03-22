package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type userUsecase struct {
	UserRepo repository.UserRepository
}

type UserUsecase interface {
	GetById(id string) (model.User, error)
	GetByEmail(email string) (model.User, error)
}

func NewUserUsecase(userRepo repository.UserRepository) *userUsecase {
	return &userUsecase{
		UserRepo: userRepo,
	}
}

func (u *userUsecase) GetByEmail(email string) (model.User, error) {
	user, err := u.UserRepo.GetByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userUsecase) GetById(id string) (model.User, error) {
	user, err := u.UserRepo.GetById(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
