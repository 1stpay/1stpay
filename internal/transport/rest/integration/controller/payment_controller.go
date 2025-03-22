package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/common/restdto"
	"github.com/1stpay/1stpay/internal/transport/rest/integration/helpers"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentUsecase usecase.PaymentUsecase
}

func NewPaymentController(paymentUsecase usecase.PaymentUsecase) *PaymentController {
	return &PaymentController{
		PaymentUsecase: paymentUsecase,
	}
}

func (con *PaymentController) Create(c *gin.Context) {
	var req restdto.PaymentCreateRestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	merchant, ok := helpers.GetMerchantOrAbort(c)
	if !ok {
		return
	}
	payment, err := con.PaymentUsecase.CreatePaymentWithWallets(req, merchant.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Err while payment create"})
		return
	}
	c.JSON(http.StatusOK, restdto.PaymentCreateResponseRestDTO{
		Id:        payment.ID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	})
}
