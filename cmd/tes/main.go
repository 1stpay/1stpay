package main

import (
	"fmt"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/domain/service/kms"
)

func main() {
	provider, err := kms.GetProvider(enum.TON)
	if err != nil {
		panic(err)
	}
	wallet, err := provider.Create()
	if err != nil {
		panic(err)
	}
	fmt.Println(wallet)

}
