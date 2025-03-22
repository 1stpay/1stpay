package controller

import (
	"fmt"
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/frontend/restdto"
	"github.com/gin-gonic/gin"
)

type FrontendPaymentController struct {
	PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentController(paymentUsecase usecase.PaymentUsecase) *FrontendPaymentController {
	return &FrontendPaymentController{
		PaymentUsecase: paymentUsecase,
	}
}

func (con *FrontendPaymentController) Get(c *gin.Context) {
	paymentId := c.Param("id")
	payment, paymentAddresses, err := con.PaymentUsecase.GetPaymentWithAddresses(paymentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error"})
		return
	}
	merchantDTO := restdto.MerchatDetailResponseDTO{
		Name: payment.Merchant.Name,
	}

	var paymentAddressDTOs []restdto.PaymentAddressDetailResponseDTO
	for _, addr := range paymentAddresses {
		paDTO := restdto.PaymentAddressDetailResponseDTO{
			PublicKey:       addr.PublicKey,
			RequestedAmount: fmt.Sprintf("%.8f", addr.RequestedAmount),
			Token: restdto.TokenDetailRestDTO{
				Name:   addr.Token.Name,
				Symbol: addr.Token.Symbol,
				Logo:   addr.Token.Logo,
				Blockchain: restdto.BlockchainDetailRestDTO{
					Name: addr.Token.Blockchain.Name,
					Logo: addr.Token.Blockchain.Logo,
				},
			},
		}
		paymentAddressDTOs = append(paymentAddressDTOs, paDTO)
	}

	response := restdto.PaymentDetailResponseRestDTO{
		RequestedAmount:    payment.RequestedAmount,
		Email:              payment.InvoiceEmail,
		ExpiresAt:          payment.ExpiresAt,
		Merchant:           merchantDTO,
		PaymentAddressList: paymentAddressDTOs,
	}

	c.JSON(http.StatusOK, response)

}
