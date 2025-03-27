package blockchain_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/repository"
)

func InitBlockchainServices(blockchainRepo repository.BlockchainRepository) (map[string]BlockchainService, error) {
	blockchains, err := blockchainRepo.ListActive()
	if err != nil {
		return nil, fmt.Errorf("failed to query blockchain list")
	}
	services := make(map[string]BlockchainService)

	for _, bc := range blockchains {
		var cfg map[string]string
		if err := json.Unmarshal(bc.Config, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config for blockchain %s: %w", bc.ID, err)
		}

		rpcURL, ok := cfg["rpc_url"]
		if !ok || rpcURL == "" {
			return nil, fmt.Errorf("rpc_url not found or empty for blockchain %s", bc.ID)
		}

		switch bc.ChainType {
		case enum.EVM:
			chainId, ok := cfg["chain_id"]
			if !ok || chainId == "" {
				return nil, fmt.Errorf("chain id not found or empty for blockchain %s", bc.ID)
			}
			chainIdInt, err := strconv.Atoi(chainId)
			if err != nil {
				return nil, fmt.Errorf("incorrect chain id for blockchain %s", bc.ID)
			}
			service, err := NewEthereumService(rpcURL, chainIdInt)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize Ethereum service for blockchain %s: %w", bc.ID, err)
			}
			services[bc.ID.String()] = service

		case enum.TRON:
			service, err := NewTronService(rpcURL)
			if err != nil {
				return nil, fmt.Errorf("failed to initialize Ethereum service for blockchain %s: %w", bc.ID, err)
			}
			services[bc.ID.String()] = service
		case enum.TON:
			return nil, errors.New("TON service not implemented yet")
		case enum.SOLANA:
			return nil, errors.New("solana service not implemented yet")
		default:
			return nil, fmt.Errorf("unsupported chain type: %s", bc.ChainType)
		}
	}

	return services, nil
}
