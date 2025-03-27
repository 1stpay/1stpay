package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"github.com/google/uuid"
)

type merchantAPIKeyUsecase struct {
	repo repository.MerchantRepository
}

type MerchantAPIKeyUsecase interface {
	CreateAPIKey(merchantID uuid.UUID, expiresAt *time.Time) (model.MerchantAPIKey, string, error)
	ListAPIKeys(merchantID uuid.UUID) ([]model.MerchantAPIKey, error)
	ValidateAPIKey(apiKey string) (model.Merchant, error)
	DeactivateAPIKey(apiKeyID uuid.UUID) error
}

func NewMerchantAPIKeyUsecase(repo repository.MerchantRepository) MerchantAPIKeyUsecase {
	return &merchantAPIKeyUsecase{
		repo: repo,
	}
}

// CreateAPIKey creates a new API key for a merchant
func (u *merchantAPIKeyUsecase) CreateAPIKey(merchantID uuid.UUID, expiresAt *time.Time) (model.MerchantAPIKey, string, error) {
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return model.MerchantAPIKey{}, "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	rawKey := hex.EncodeToString(rawBytes)

	hashBytes := sha256.Sum256(rawBytes)
	hashedKey := hex.EncodeToString(hashBytes[:])

	name := ""
	if len(rawKey) >= 5 {
		name = rawKey[len(rawKey)-5:]
	} else {
		name = rawKey
	}

	apiKeyModel := model.MerchantAPIKey{
		MerchantID: merchantID,
		Name:       name,
		APIKey:     hashedKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	createdKey, err := u.repo.CreateMerchantAPIKey(apiKeyModel)
	if err != nil {
		return model.MerchantAPIKey{}, "", err
	}
	return createdKey, rawKey, nil
}

func (u *merchantAPIKeyUsecase) ListAPIKeys(merchantID uuid.UUID) ([]model.MerchantAPIKey, error) {
	keys, err := u.repo.ListMerchantAPIKey(merchantID.String())
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (u *merchantAPIKeyUsecase) DeactivateAPIKey(tokenID uuid.UUID) error {
	if err := u.repo.DeactivateMerchantAPIKey(tokenID.String()); err != nil {
		return fmt.Errorf("failed to deactivate merchant token %s: %w", tokenID.String(), err)
	}
	return nil
}

func (u *merchantAPIKeyUsecase) ValidateAPIKey(rawKey string) (model.Merchant, error) {
	rawBytes, err := hex.DecodeString(rawKey)
	if err != nil {
		return model.Merchant{}, fmt.Errorf("failed to decode API key: %w", err)
	}
	hashBytes := sha256.Sum256(rawBytes)
	hashedKey := hex.EncodeToString(hashBytes[:])

	apiKeyRecord, err := u.repo.GetMerchantAPIKeyByHash(hashedKey)
	if err != nil {
		return model.Merchant{}, err
	}
	if !apiKeyRecord.IsActive {
		return model.Merchant{}, errors.New("API key is deactivated")
	}

	merchant, err := u.repo.GetMerchantById(apiKeyRecord.MerchantID.String())
	if err != nil {
		return model.Merchant{}, err
	}
	return merchant, nil
}
