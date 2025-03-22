package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1stpay/1stpay/internal/config"
	invoicechecker "github.com/1stpay/1stpay/internal/domain/service/invoice_checker"
	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
)

func runInvoiceChecker() {
	app := config.App()

	blockService, err := blockchain_service.InitBlockchainServices(app.Deps.Repos.BlockchainRepo)
	if err != nil {
		panic(err)
	}

	invChecker := invoicechecker.NewInvoiceChecker(
		app.Postgres,
		app.Deps.Repos.PaymentRepo,
		app.Deps.Repos.PaymentAddressRepo,
		blockService,
		1*time.Second,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, cancelling context...")
		cancel()
	}()

	log.Println("Starting Invoice Checker service...")
	if err := invChecker.Start(ctx); err != nil {
		log.Printf("Invoice Checker terminated: %v", err)
	}
}

func main() {
	runInvoiceChecker()
}
