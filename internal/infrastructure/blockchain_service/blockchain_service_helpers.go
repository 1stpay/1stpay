package blockchain_service

import (
	"math"
	"math/big"
)

func ConvertBigIntToFloat(amount *big.Int, decimals int) float64 {
	f := new(big.Float).SetInt(amount)
	divisor := new(big.Float).SetFloat64(math.Pow10(decimals))
	result, _ := new(big.Float).Quo(f, divisor).Float64()
	return result
}

func ConvertFloatToBigInt(amount float64, decimals int) *big.Int {
	multiplier := math.Pow10(decimals)
	bf := new(big.Float).SetFloat64(amount * multiplier)
	result := new(big.Int)
	bf.Int(result)
	return result
}
