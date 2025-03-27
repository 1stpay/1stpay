package mock

import (
	"github.com/stretchr/testify/mock"
)

type MockPriceService struct {
	mock.Mock
}

func (m *MockPriceService) GetPrice(key string) (float64, error) {
	args := m.Called(key)
	return args.Get(0).(float64), args.Error(1)
}
