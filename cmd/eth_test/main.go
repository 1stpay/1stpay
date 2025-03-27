package main

import (
	"context"
	"fmt"

	"github.com/1stpay/1stpay/internal/infrastructure/blockchain_service"
)

func queryBalance(ctx context.Context, service blockchain_service.EthereumService, sender, receiver string) (float64, float64) {
	senderWalletBalance, err := service.GetNativeBalance(ctx, sender)
	if err != nil {
		panic("error while balance query")
	}
	receiverWalletBalance, err := service.GetNativeBalance(ctx, receiver)
	if err != nil {
		panic("error while balance query")
	}
	senderWalletBalanceFloat := blockchain_service.ConvertBigIntToFloat(senderWalletBalance, 18)
	receiverWalletBalanceFloat := blockchain_service.ConvertBigIntToFloat(receiverWalletBalance, 18)
	fmt.Printf(
		"Sender wallet balance: %f ETH\nReceiver wallet balance: %f ETH",
		senderWalletBalanceFloat,
		receiverWalletBalanceFloat,
	)
	return senderWalletBalanceFloat, receiverWalletBalanceFloat
}

func transfer(ctx context.Context, service *blockchain_service.EthereumService, senderWallet, receiverWallet string) {
	sendAmount := blockchain_service.ConvertFloatToBigInt(1.5, 18)
	result, err := service.TransferNative(
		ctx,
		"0x1369173afdaeecf3accd001efe8e7cfc59abae26c904cb1463f2823ca6da62fb",
		receiverWallet,
		sendAmount,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	queryBalance(ctx, *service, senderWallet, receiverWallet)
}

func transferAll(ctx context.Context, service *blockchain_service.EthereumService, receiverWallet string) {
	result, err := service.TransferNativeRemaining(
		ctx,
		"0x1369173afdaeecf3accd001efe8e7cfc59abae26c904cb1463f2823ca6da62fb",
		receiverWallet,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func transferAllToken(ctx context.Context, service *blockchain_service.EthereumService, receiverWallet string) {
	tokenAddress := "0x85F2b9552b097E7CEc7BE9791BB6E437BE72e9a0"
	res, err := service.TransferTokenRemaining(
		ctx,
		"0xd2caa4fa1ce834b07c1ed8b55a61707dc209840093f24a8e40b3bba1cc09ac99",
		tokenAddress,
		receiverWallet,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func main() {
	service, err := blockchain_service.NewEthereumService("http://127.0.0.1:8545", 1337)
	if err != nil {
		panic("Error while connecting to ETH RPC")
	}
	senderWallet := "0x8ac5148Fd5669BA6bcc6A2F061AE655F86EA79F4"
	receiverWallet := "0x7c5B9e1b438b4c8383c23eacaB4F3Ea47953A8bE"
	ctx := context.Background()
	queryBalance(ctx, *service, senderWallet, receiverWallet)
	transferAllToken(ctx, service, receiverWallet)
	queryBalance(ctx, *service, senderWallet, receiverWallet)
}
