package route_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/1stpay/1stpay/test"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizedUserAccess(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine
	user, accessToken := testConfig.TestFactory.CreateUser()

	userReponse := test.PerformRequest(router, http.MethodGet, "/frontend/api/v1/user/me/", nil, accessToken)

	assert.Equal(t, http.StatusOK, userReponse.Code, "Expected status 200 for authorized GET request")

	var userInfo map[string]interface{}
	err := json.Unmarshal(userReponse.Body.Bytes(), &userInfo)
	assert.NoError(t, err, "Error parsing JSON response for authorized access")
	assert.Equal(t, user.Email, userInfo["email"], "Returned email should match the registered email")
}
