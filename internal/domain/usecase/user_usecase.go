package usecase

import (
	"errors"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserUsecase struct {
	UserRepo repository.UserRepositoryInterface
}

type UserUsecaseInterface interface {
	GetByEmail(remail string) (model.User, error)
	GetUserFromContext(c *gin.Context) (model.User, error)
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

func (u *UserUsecase) GetUserFromContext(c *gin.Context) (model.User, error) {
	claims, exists := c.Get(middleware.ContextUserKey)
	if !exists {
		return model.User{}, errors.New("user not found in context")
	}

	userClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return model.User{}, errors.New("invalid user data in context")
	}

	userId, ok := userClaims["user_id"].(string)
	if !ok || userId == "" {
		return model.User{}, errors.New("email not found in token claims")
	}

	user, err := u.UserRepo.GetById(userId)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
