package usecase_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/1stpay/1stpay/internal/domain/enum"
	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/model"
	"github.com/1stpay/1stpay/internal/transport/rest/common/restdto"
	"github.com/1stpay/1stpay/test"
	"github.com/1stpay/1stpay/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPaymentUsecase(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../")
	t.Run("Positive balance scenario", func(t *testing.T) {
		user, _ := testConfig.TestFactory.CreateUser()
		merchant := testConfig.TestFactory.CreateMerchant(user.ID.String())
		blockchainList := testConfig.TestFactory.CreateBlockchainList()
		tokenList := testConfig.TestFactory.CreateTokenList(blockchainList)
		testConfig.TestFactory.CreateMerchantTokenList(merchant, tokenList)
		mockPriceService := new(mock.MockPriceService)
		mockPriceService.On("GetPrice", "USDT").Return(100.0, nil)
		paymentUsecase := usecase.NewPaymentUsecase(
			testConfig.Database.GormDB,
			testConfig.Repos.PaymentRepo,
			testConfig.Repos.PaymentAddressRepo,
			testConfig.Repos.MerchantRepo,
			mockPriceService,
		)
		payment, err := paymentUsecase.CreatePaymentWithWallets(restdto.PaymentCreateRestDTO{RequestedAmount: 100}, merchant.ID)
		assert.NoError(t, err)
		assert.Equal(t, enum.PaymentStatusPending, payment.Status)
		assert.WithinDuration(t, time.Now().Add(time.Hour), *payment.ExpiresAt, time.Minute)
		paymentAddressList, err := testConfig.Repos.PaymentAddressRepo.ListByPaymentId(payment.ID)
		assert.NoError(t, err)
		assert.Greater(t, len(paymentAddressList), 0)
		var merchantToken model.MerchantToken
		err = testConfig.Database.GormDB.
			Where("merchant_id = ? AND token_id = ?", merchant.ID, paymentAddressList[0].Token.ID).
			First(&merchantToken).Error
		assert.NoError(t, err)
		initialBalance := merchantToken.Balance

		balance := new(big.Int).SetInt64(int64(paymentAddressList[0].RequestedAmountWei))
		err = paymentUsecase.ConfirmInvoice(payment, paymentAddressList[0], balance)
		assert.NoError(t, err)

		var updatedPayment model.Payment
		err = testConfig.Database.GormDB.Where("id = ?", payment.ID).First(&updatedPayment).Error
		assert.NoError(t, err)
		assert.Equal(t, enum.PaymentStatusCompleted, updatedPayment.Status)

		var updatedPaymentAddress model.PaymentAddress
		err = testConfig.Database.GormDB.Where("id = ?", paymentAddressList[0].ID).First(&updatedPaymentAddress).Error
		assert.NoError(t, err)
		assert.Equal(t, int(balance.Int64()), updatedPaymentAddress.PaidAmountWei)

		err = testConfig.Database.GormDB.
			Where("merchant_id = ? AND token_id = ?", merchant.ID, paymentAddressList[0].Token.ID).
			First(&merchantToken).Error
		assert.NoError(t, err)
		assert.Equal(t, initialBalance+1.0, merchantToken.Balance)

	})
}
