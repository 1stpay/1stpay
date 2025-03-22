package invoicechecker

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"gorm.io/gorm"
)

type InvoiceChecker interface {
	Start(ctx context.Context) error
	CheckInvoices(ctx context.Context) error
	CheckInvoiceReceipt() error
}

type invoiceChecker struct {
	paymentRepo        repository.PaymentRepositoryInterface
	paymentAddressRepo repository.PaymentAddressRepository
	db                 *gorm.DB
	blockchainServices map[string]blockchain_service.BlockchainService
	pollInterval       time.Duration
}

func NewInvoiceChecker(
	db *gorm.DB,
	paymentRepo repository.PaymentRepositoryInterface,
	paymentAddressRepo repository.PaymentAddressRepository,
	blockchainServices map[string]blockchain_service.BlockchainService,
	pollInterval time.Duration,
) InvoiceChecker {
	return &invoiceChecker{
		db:                 db,
		paymentRepo:        paymentRepo,
		paymentAddressRepo: paymentAddressRepo,
		blockchainServices: blockchainServices,
		pollInterval:       pollInterval,
	}
}

func (ic *invoiceChecker) Start(ctx context.Context) error {
	ticker := time.NewTicker(ic.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := ic.CheckInvoices(ctx); err != nil {
				log.Printf("Invoice check error: %v", err)
			}
		}
	}
}

func (ic *invoiceChecker) CheckInvoiceReceipt() error {
	return nil
}

func (ic *invoiceChecker) CheckInvoices(ctx context.Context) error {
	fmt.Println("Starting invoice check...")

	var activeInvoices []model.Payment
	if err := ic.db.Where("status = ?", "pending").Find(&activeInvoices).Error; err != nil {
		return err
	}

	for _, inv := range activeInvoices {
		var addresses []model.PaymentAddress
		if err := ic.db.
			Where("payment_id = ?", inv.ID).
			Preload("Token").
			Preload("Token.Blockchain").
			Find(&addresses).Error; err != nil {
			log.Printf("Error retrieving addresses for invoice %s: %v", inv.ID, err)
			continue
		}

		for _, addr := range addresses {
			bcID := addr.Token.Blockchain.ID.String()
			service, ok := ic.blockchainServices[bcID]
			if !ok {
				log.Printf("No blockchain service found for blockchain ID %s", bcID)
				continue
			}

			var balance *big.Int
			var err error

			if addr.Token.IsNative {
				balance, err = service.GetNativeBalance(addr.PublicKey)
			} else {
				if addr.Token.ContractAddress == "" {
					log.Printf("Token %s is non-native but ContractAddress is empty", addr.Token.Symbol)
					continue
				}
				balance, err = service.GetTokenBalance(addr.PublicKey, addr.Token.ContractAddress)
			}
			if err != nil {
				log.Printf("Error getting balance for address %s on blockchain %s: %v", addr.PublicKey, bcID, err)
				continue
			}

			reqAmount := big.NewInt(int64(addr.RequestedAmountWei))
			if balance.Cmp(reqAmount) >= 0 {
				decimals := addr.Token.Decimals
				factor := math.Pow10(decimals)

				fBalance := new(big.Float).SetInt(balance)
				fmt.Println(fBalance)
				paidAmountFloat, _ := new(big.Float).Quo(fBalance, big.NewFloat(factor)).Float64()

				// if err := ic.db.Model(&model.PaymentAddress{}).
				// 	Where("id = ?", addr.ID).
				// 	Updates(map[string]interface{}{
				// 		"paid_amount":     paidAmountFloat,
				// 		"paid_amount_wei": balance.Int64(),
				// 	}).Error; err != nil {
				// 	log.Printf("Error updating PaymentAddress %s: %v", addr.ID, err)
				// 	continue
				// }

				// if err := ic.db.Model(&model.Payment{}).
				// 	Where("id = ?", inv.ID).
				// 	Updates(map[string]interface{}{
				// 		"status":        enum.PaymentStatusCompleted,
				// 		"used_token_id": addr.Token.ID,
				// 	}).Error; err != nil {
				// 	log.Printf("Error updating Payment %s: %v", inv.ID, err)
				// } else {
				log.Printf("Invoice %s confirmed. Address %s balance (%s minimal units) >= requested (%s minimal units). Converted value: %f",
					inv.ID, addr.PublicKey, balance.String(), reqAmount.String(), paidAmountFloat)
				// }
				break
			} else {
				log.Printf("Invoice %s, address %s: balance (%s minimal units) is less than requested (%s minimal units)",
					inv.ID, addr.PublicKey, balance.String(), reqAmount.String())
			}
		}
	}

	return nil
}
