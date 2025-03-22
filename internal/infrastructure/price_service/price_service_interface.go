package price_service

import "fmt"

type priceService struct {
	apiKey string
}

type PriceService interface {
	GetPrice(asset string) (float64, error)
}

func NewPriceService(apiKey string) PriceService {
	return &priceService{
		apiKey: apiKey,
	}
}

func (s *priceService) GetPrice(asset string) (float64, error) {
	prices := map[string]float64{
		"TRX":  0.065,
		"USDT": 1.0,
		"ETH":  1800,
		"ZRO":  3.23,
	}
	price, ok := prices[asset]
	if !ok {
		return 0, fmt.Errorf("unknown asset: %s", asset)
	}
	return price, nil
}
