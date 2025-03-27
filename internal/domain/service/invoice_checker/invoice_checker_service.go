package invoicechecker

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/repository"
	"gorm.io/gorm"
)

type InvoiceChecker interface {
	Start(ctx context.Context) error
	CheckInvoices(ctx context.Context) error
	CheckInvoiceReceipt() error
	CheckExpiredPayments(ctx context.Context) error
}

type invoiceChecker struct {
	paymentRepo        repository.PaymentRepository
	paymentAddressRepo repository.PaymentAddressRepository
	paymentUsecase     usecase.PaymentUsecase
	db                 *gorm.DB
	blockchainServices map[string]blockchain_service.BlockchainService
	pollInterval       time.Duration
}

func NewInvoiceChecker(
	db *gorm.DB,
	paymentRepo repository.PaymentRepository,
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
	expiredTicker := time.NewTicker(ic.pollInterval)
	defer expiredTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := ic.CheckInvoices(ctx); err != nil {
				log.Printf("Invoice check error: %v", err)
			}
		case <-expiredTicker.C:
			if err := ic.CheckExpiredPayments(ctx); err != nil {
				log.Printf("Expired payment check error: %v", err)
			}
		}
	}
}

func (ic *invoiceChecker) CheckInvoiceReceipt() error {
	return nil
}

func (ic *invoiceChecker) CheckInvoices(ctx context.Context) error {
	fmt.Println("Starting invoice check...")
	activeInvoices, err := ic.paymentRepo.GetPaymentListWithStatus(enum.PaymentStatusPending)
	if err != nil {
		return err
	}

	for _, inv := range activeInvoices {
		addresses, err := ic.paymentAddressRepo.ListByPaymentId(inv.ID)
		if err != nil {
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
				balance, err = service.GetNativeBalance(ctx, addr.PublicKey)
			} else {
				if addr.Token.ContractAddress == "" {
					log.Printf("Token %s is non-native but ContractAddress is empty", addr.Token.Symbol)
					continue
				}
				balance, err = service.GetTokenBalance(ctx, addr.PublicKey, addr.Token.ContractAddress)
			}
			if err != nil {
				log.Printf("Error getting balance for address %s on blockchain %s: %v", addr.PublicKey, bcID, err)
				continue
			}

			reqAmount := big.NewInt(int64(addr.RequestedAmountWei))
			if balance.Cmp(reqAmount) >= 0 {
				if err := ic.paymentUsecase.ConfirmInvoice(inv, addr, balance); err != nil {
					log.Printf("Error confirming invoice %s: %v", inv.ID, err)
					continue
				}
				log.Printf("Invoice %s confirmed. Address %s balance (%s minimal units) >= requested (%s minimal units).",
					inv.ID, addr.PublicKey, balance.String(), reqAmount.String())

				break
			} else {
				log.Printf("Invoice %s, address %s: balance (%s minimal units) is less than requested (%s minimal units)",
					inv.ID, addr.PublicKey, balance.String(), reqAmount.String())
			}
		}
	}

	return nil
}

func (ic *invoiceChecker) CheckExpiredPayments(ctx context.Context) error {
	fmt.Println("payment cancel")
	now := time.Now()
	result := ic.db.Model(&model.Payment{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND status = ?", now, "pending").
		Update("status", enum.PaymentStatusFailed)
	if result.Error != nil {
		return result.Error
	}
	log.Printf("%d expired payments cancelled", result.RowsAffected)
	return nil
}
