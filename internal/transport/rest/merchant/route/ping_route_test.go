package route_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/1stpay/1stpay/test"
	"github.com/stretchr/testify/assert"
)

type PingResponse struct {
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

func TestPingRoute(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")

	router := testConfig.GinEngine
	response := test.PerformRequest(router, http.MethodGet, "/frontend/api/v1/ping", nil)
	assert.Equal(t, http.StatusOK, response.Code, "Expected status 200")

	var pingResp PingResponse
	err := json.Unmarshal(response.Body.Bytes(), &pingResp)
	assert.NoError(t, err, "Error while response JSON parsing")
	assert.Equal(t, 200, pingResp.Status, "Status field not equal to expected value")
	assert.True(t, pingResp.Success, "success field should be true")
}
