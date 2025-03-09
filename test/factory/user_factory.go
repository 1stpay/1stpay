package factory

import (
	"fmt"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	restdto "github.com/1stpay/1stpay/internal/transport/rest/frontend/rest_dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestUserFactory struct {
	db          *gorm.DB
	userRepo    *repository.UserRepository
	userUsecase *usecase.UserUsecase
	authUsecase *usecase.AuthUsecase
}

func NewTestUserFactory(db *gorm.DB) *TestUserFactory {
	repo := repository.NewUserRepository(db)
	return &TestUserFactory{
		db:          db,
		userRepo:    repo,
		userUsecase: usecase.NewUserUsecase(repo),
		authUsecase: usecase.NewAuthUsecase(repo),
	}
}

func (f *TestUserFactory) CreateUser() (model.User, string) {
	uniqueEmail := fmt.Sprintf("testuser_%s@example.com", uuid.New().String())
	registerData := restdto.RegisterRequest{
		Email:    uniqueEmail,
		Password: "Secret",
	}
	user, accessToken, err := f.authUsecase.Register(registerData)
	if err != nil {
		panic("Error while test user creation")
	}
	return user, accessToken
}
