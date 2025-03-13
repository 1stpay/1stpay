package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type BlockchainRepository struct {
	db *gorm.DB
}

type BlockchainRepositoryInterface interface {
	ListActive() ([]model.Blockchain, error)
	Create(blockchain model.Blockchain) (model.Blockchain, error)
}

func NewBlockchainRepository(db *gorm.DB) *BlockchainRepository {
	return &BlockchainRepository{
		db: db,
	}
}

func (r *BlockchainRepository) ListActive() ([]model.Blockchain, error) {
	var blockchainList []model.Blockchain
	if err := r.db.Where("is_active = ?", true).Find(&blockchainList).Error; err != nil {
		return []model.Blockchain{}, err
	}
	return blockchainList, nil
}

func (r *BlockchainRepository) Create(blockchain model.Blockchain) (model.Blockchain, error) {
	if err := r.db.Create(&blockchain).Error; err != nil {
		return model.Blockchain{}, err
	}
	return blockchain, nil
}
