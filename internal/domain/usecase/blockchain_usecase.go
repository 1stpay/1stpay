package usecase

import (
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
)

type BlockchainUsecase struct {
	BlockchainRepo repository.BlockchainRepositoryInterface
}

type BlockchainUsecaseInterface interface {
	ListActive() ([]model.Blockchain, error)
	Create(blockchain model.Blockchain) (model.Blockchain, error)
}

func NewBlockchainUsecase(repo repository.BlockchainRepositoryInterface) *BlockchainUsecase {
	return &BlockchainUsecase{
		BlockchainRepo: repo,
	}
}

func (u *BlockchainUsecase) ListActive() ([]model.Blockchain, error) {
	blockachainList, err := u.BlockchainRepo.ListActive()
	if err != nil {
		return []model.Blockchain{}, err
	}
	return blockachainList, err
}

func (u *BlockchainUsecase) Create(blockchain model.Blockchain) (model.Blockchain, error) {
	blockchain, err := u.BlockchainRepo.Create(blockchain)
	if err != nil {
		return model.Blockchain{}, err
	}
	return blockchain, err
}
