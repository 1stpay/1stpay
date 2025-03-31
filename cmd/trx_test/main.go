package main

import (
	"context"
	"fmt"
	"log"

	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
)

const (
	senderPk        string = "0x"
	senderAddress   string = "0x"
	receiverAddress string = "0x"
)

func queryBalance(
	ctx context.Context,
	service blockchain_service.TronService,
	sender,
	receiver string,
) (float64, float64) {
	senderWalletBalance, err := service.GetNativeBalance(ctx, sender)
	if err != nil {
		panic(err)
	}
	receiverWalletBalance, err := service.GetNativeBalance(ctx, receiver)
	if err != nil {
		panic(err)
	}
	senderWalletBalanceFloat := blockchain_service.ConvertBigIntToFloat(
		senderWalletBalance, 6,
	)
	receiverWalletBalanceFloat := blockchain_service.ConvertBigIntToFloat(
		receiverWalletBalance, 6,
	)
	fmt.Printf(
		"Sender wallet balance: %f TRX\nReceiver wallet balance: %f TRX",
		senderWalletBalanceFloat,
		receiverWalletBalanceFloat,
	)
	return senderWalletBalanceFloat, receiverWalletBalanceFloat
}

func transfer(ctx context.Context, service *blockchain_service.TronService, senderWallet, receiverWallet string) {
	sendAmount := blockchain_service.ConvertFloatToBigInt(0.5, 6)
	result, err := service.TransferNative(
		ctx,
		senderPk,
		receiverWallet,
		sendAmount,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	queryBalance(ctx, *service, senderWallet, receiverWallet)
}

func transferAll(ctx context.Context, service *blockchain_service.TronService, receiverWallet string) {
	result, err := service.TransferNativeRemaining(
		ctx,
		senderPk,
		receiverWallet,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func transferAllToken(ctx context.Context, service *blockchain_service.TronService, receiverWallet string) {
	tokenAddress := "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"
	res, err := service.TransferTokenRemaining(
		ctx,
		senderPk,
		tokenAddress,
		receiverWallet,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	service, err := blockchain_service.NewTronService("https://nile.trongrid.io")
	// if err != nil {
	// 	panic("Error while connecting to TRON RPC")
	// }
	// senderWallet := senderAddress
	receiverWallet := receiverAddress
	ctx := context.Background()
	// addr, err := tron.GetTronAddressFromPrivateKey(senderPk)
	// res, err := service.CreateTransaction(ctx, senderWallet, receiverWallet, *big.NewInt(1 * 1_000_000))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(addr)
	// log.Println(res)

	// queryBalance(ctx, *service, senderWallet, receiverWallet)
	transferAllToken(ctx, service, receiverWallet)
	// transfer(ctx, service, senderAddress, receiverWallet)
	// queryBalance(ctx, *service, senderWallet, receiverWallet)
}
