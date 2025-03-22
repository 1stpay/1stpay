package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/integration/restdto"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/helpers"
	frontend_dto "github.com/1stpay/1stpay/internal/transport/rest/merchant/rest_dto"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentUsecase  usecase.PaymentUsecaseInterface
	MerchantUsecase usecase.MerchantUsecase
	UserUsecase     usecase.UserUsecase
}

func NewPaymentController(paymentUsecase usecase.PaymentUsecaseInterface, merchantUsecase usecase.MerchantUsecase, userUsecase usecase.UserUsecase) *PaymentController {
	return &PaymentController{
		PaymentUsecase:  paymentUsecase,
		MerchantUsecase: merchantUsecase,
		UserUsecase:     userUsecase,
	}
}

func (con *PaymentController) Create(c *gin.Context) {
	var req restdto.InvoiceCreateRestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	user, ok := helpers.GetUserOrAbort(c, con.UserUsecase)
	if !ok {
		return
	}
	merchant, err := con.MerchantUsecase.GetMerchantByUserId(user.ID.String())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	payment, err := con.PaymentUsecase.CreatePaymentWithWallets(req, merchant.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Err while payment create"})
		return
	}
	c.JSON(http.StatusOK, frontend_dto.PaymentCreateResponseRestDTO{
		Id:        payment.ID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	})
}
