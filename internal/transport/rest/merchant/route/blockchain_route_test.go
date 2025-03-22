package route_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/1stpay/1stpay/test"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainList(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine
	t.Run("Positive blockchain list", func(t *testing.T) {
		blockachainList := testConfig.TestFactory.CreateBlockchainList()
		blockachainListResponse := test.PerformRequest(router, http.MethodGet, "/frontend/api/v1/blockchain/list/", nil)
		assert.Equal(t, http.StatusOK, blockachainListResponse.Code, "Expected status 200 for merchant creation")

		var blockachainListJson []restdto.BlockchainListResponseDTO
		err := json.Unmarshal(blockachainListResponse.Body.Bytes(), &blockachainListJson)
		assert.NoError(t, err, "Error parsing JSON response for merchant creation")
		assert.Equal(t, len(blockachainList), len(blockachainListJson), "Returned user ID should match the registered user's ID")
	})
}
