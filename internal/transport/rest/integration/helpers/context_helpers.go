package helpers

import (
	"net/http"

	"github.com/1stpay/1stpay/internal/model"
	"github.com/gin-gonic/gin"
)

func GetMerchantOrAbort(c *gin.Context) (model.Merchant, bool) {
	merchantData, exists := c.Get("merchant")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant"})
		return model.Merchant{}, false
	}
	merchant, ok := merchantData.(model.Merchant)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant type"})
		return model.Merchant{}, false
	}
	return merchant, true
}
