package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type UserUsecase struct {
	UserRepo repository.UserRepositoryInterface
}

type UserUsecaseInterface interface {
	GetById(id string) (model.User, error)
	GetByEmail(email string) (model.User, error)
}

func NewUserUsecase(userRepo repository.UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{
		UserRepo: userRepo,
	}
}

func (u *UserUsecase) GetByEmail(email string) (model.User, error) {
	user, err := u.UserRepo.GetByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *UserUsecase) GetById(id string) (model.User, error) {
	user, err := u.UserRepo.GetById(id)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
