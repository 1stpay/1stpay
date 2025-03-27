package controller

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/domain/usecase"
	"github.com/1stpay/1stpay/internal/transport/rest/common/restdto"
	"github.com/1stpay/1stpay/internal/transport/rest/merchant/helpers"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentUsecase  usecase.PaymentUsecase
	MerchantUsecase usecase.MerchantUsecase
	UserUsecase     usecase.UserUsecase
}

func NewPaymentController(paymentUsecase usecase.PaymentUsecase, merchantUsecase usecase.MerchantUsecase, userUsecase usecase.UserUsecase) *PaymentController {
	return &PaymentController{
		PaymentUsecase:  paymentUsecase,
		MerchantUsecase: merchantUsecase,
		UserUsecase:     userUsecase,
	}
}

// Create godoc
// @Summary Create a new payment
// @Description Creates a new payment with associated wallets for the current merchant.
// @Tags Payment
// @Accept json
// @Produce json
// @Param payload body restdto.PaymentCreateRestDTO true "Payment creation payload"
// @Success 200 {object} restdto.PaymentCreateResponseRestDTO "Payment created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Merchant not found or error during payment creation"
// @Router /merchant/api/v1/payments/ [post]
func (con *PaymentController) Create(c *gin.Context) {
	var req restdto.PaymentCreateRestDTO
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
	c.JSON(http.StatusOK, restdto.PaymentCreateResponseRestDTO{
		Id:        payment.ID,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	})
}
