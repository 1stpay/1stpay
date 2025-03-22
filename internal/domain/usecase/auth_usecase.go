package usecase

import (
	"errors"
	"time"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	UserRepo  repository.UserRepository
	jwtSecret string
}

type AuthUsecase interface {
	Register(registerData restdto.RegisterRequest) (model.User, string, error)
	Login(loginData restdto.LoginRequest) (model.User, string, error)
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{
		UserRepo:  userRepo,
		jwtSecret: "hehe",
	}
}

func (u *authUsecase) Register(registerData restdto.RegisterRequest) (model.User, string, error) {
	_, err := u.UserRepo.GetByEmail(registerData.Email)
	if err == nil {
		return model.User{}, "", errors.New("user with this email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, "", err
	}
	user := model.User{
		Email:    registerData.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	createdUser, err := u.UserRepo.Create(user)
	if err != nil {
		return model.User{}, "", err
	}
	token, err := u.generateToken(createdUser)
	if err != nil {
		return model.User{}, "", err
	}
	return createdUser, token, nil
}

func (u *authUsecase) Login(loginData restdto.LoginRequest) (model.User, string, error) {
	user, err := u.UserRepo.GetByEmail(loginData.Email)
	if err != nil {
		return model.User{}, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return model.User{}, "", errors.New("invalid credentials")
	}

	token, err := u.generateToken(user)
	if err != nil {
		return model.User{}, "", err
	}

	return user, token, nil
}

func (u *authUsecase) generateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
