package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type blockchainUsecase struct {
	BlockchainRepo repository.BlockchainRepository
}

type BlockchainUsecase interface {
	ListActive() ([]model.Blockchain, error)
	Create(blockchain model.Blockchain) (model.Blockchain, error)
}

func NewBlockchainUsecase(repo repository.BlockchainRepository) BlockchainUsecase {
	return &blockchainUsecase{
		BlockchainRepo: repo,
	}
}

func (u *blockchainUsecase) ListActive() ([]model.Blockchain, error) {
	blockachainList, err := u.BlockchainRepo.ListActive()
	if err != nil {
		return []model.Blockchain{}, err
	}
	return blockachainList, err
}

func (u *blockchainUsecase) Create(blockchain model.Blockchain) (model.Blockchain, error) {
	blockchain, err := u.BlockchainRepo.Create(blockchain)
	if err != nil {
		return model.Blockchain{}, err
	}
	return blockchain, err
}
