package repository

import (
	"github.com/1stpay/1stpay/internal/model"
	"gorm.io/gorm"
)

type blockchainRepository struct {
	db *gorm.DB
}

type BlockchainRepository interface {
	ListActive() ([]model.Blockchain, error)
	Create(blockchain model.Blockchain) (model.Blockchain, error)
}

func NewBlockchainRepository(db *gorm.DB) BlockchainRepository {
	return &blockchainRepository{
		db: db,
	}
}

func (r *blockchainRepository) ListActive() ([]model.Blockchain, error) {
	var blockchainList []model.Blockchain
	if err := r.db.Where("is_active = ?", true).Find(&blockchainList).Error; err != nil {
		return []model.Blockchain{}, err
	}
	return blockchainList, nil
}

func (r *blockchainRepository) Create(blockchain model.Blockchain) (model.Blockchain, error) {
	if err := r.db.Create(&blockchain).Error; err != nil {
		return model.Blockchain{}, err
	}
	return blockchain, nil
}
