package route_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/1stpay/1stpay/internal/transport/rest/merchant/restdto"
	"github.com/1stpay/1stpay/test"
	"github.com/stretchr/testify/assert"
)

func TestMerchantCreate(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine

	t.Run("Positive merchant creation", func(t *testing.T) {
		user, accessToken := testConfig.TestFactory.CreateUser()

		createPayload := restdto.MerchantCreateRequestDTO{
			Name: "Test merchant",
		}
		createPayloadBytes, err := json.Marshal(createPayload)
		assert.NoError(t, err, "Error marshaling request body")

		merchantResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(createPayloadBytes), accessToken)
		assert.Equal(t, http.StatusCreated, merchantResponse.Code, "Expected status 201 for merchant creation")

		var merchantInfo restdto.MerchantCreateResponseDTO
		err = json.Unmarshal(merchantResponse.Body.Bytes(), &merchantInfo)
		assert.NoError(t, err, "Error parsing JSON response for merchant creation")
		assert.Equal(t, user.ID, merchantInfo.UserID, "Returned user ID should match the registered user's ID")
	})

	// Negative scenario: missing token
	t.Run("Missing token", func(t *testing.T) {
		createPayload := restdto.MerchantCreateRequestDTO{
			Name: "Test merchant",
		}
		createPayloadBytes, err := json.Marshal(createPayload)
		assert.NoError(t, err, "Error marshaling request body for missing token")

		// No token provided
		resp := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(createPayloadBytes))
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "Expected status 401 when token is missing")
	})

	// Negative scenario: invalid token
	t.Run("Invalid token", func(t *testing.T) {
		createPayload := restdto.MerchantCreateRequestDTO{
			Name: "Test merchant",
		}
		createPayloadBytes, err := json.Marshal(createPayload)
		assert.NoError(t, err, "Error marshaling request body for invalid token")

		invalidToken := "this_is_an_invalid_token"
		resp := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(createPayloadBytes), invalidToken)
		assert.Equal(t, http.StatusUnauthorized, resp.Code, "Expected status 401 when token is invalid")
	})

	// Negative scenario: invalid JSON body
	t.Run("Invalid JSON body", func(t *testing.T) {
		invalidBody := []byte("not a valid json")
		_, accessToken := testConfig.TestFactory.CreateUser()

		resp := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(invalidBody), accessToken)
		assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected status 400 for invalid JSON body")
	})

	// Negative scenario: missing required field (e.g. missing name)
	t.Run("Missing required field", func(t *testing.T) {
		payload := map[string]string{}
		payloadBytes, err := json.Marshal(payload)
		assert.NoError(t, err, "Error marshaling request body with missing required field")
		_, accessToken := testConfig.TestFactory.CreateUser()

		resp := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(payloadBytes), accessToken)
		assert.Equal(t, http.StatusBadRequest, resp.Code, "Expected status 400 when required field is missing")
	})
	t.Run("User attempts to create multiple merchants", func(t *testing.T) {
		_, accessToken := testConfig.TestFactory.CreateUser()

		createPayload1 := restdto.MerchantCreateRequestDTO{
			Name: "First Merchant",
		}
		createPayloadBytes1, err := json.Marshal(createPayload1)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")
		firstResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(createPayloadBytes1), accessToken)
		assert.Equal(t, http.StatusCreated, firstResponse.Code, "Expected status 201 for first merchant creation")

		var firstMerchant restdto.MerchantCreateResponseDTO
		err = json.Unmarshal(firstResponse.Body.Bytes(), &firstMerchant)
		assert.NoError(t, err, "Error parsing JSON response for first merchant creation")

		createPayload2 := restdto.MerchantCreateRequestDTO{
			Name: "Second Merchant Attempt",
		}
		createPayloadBytes2, err := json.Marshal(createPayload2)
		assert.NoError(t, err, "Error marshaling request body for second merchant creation attempt")
		secondResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/merchant/", bytes.NewReader(createPayloadBytes2), accessToken)
		assert.Equal(t, http.StatusBadRequest, secondResponse.Code, "Expected status 201 for second merchant creation attempt")
	})
}

func TestMerchantDetail(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine
	t.Run("Positive merchant creation", func(t *testing.T) {
		user, accessToken := testConfig.TestFactory.CreateUser()
		merchant := testConfig.TestFactory.CreateMerchant(user.ID.String())
		merchantResponse := test.PerformRequest(router, http.MethodGet, "/frontend/api/v1/merchant/me/", nil, accessToken)
		assert.Equal(t, http.StatusOK, merchantResponse.Code, "Expected status 200 for merchant creation")

		var merchantInfo restdto.MerchantDetailResponseDTO
		err := json.Unmarshal(merchantResponse.Body.Bytes(), &merchantInfo)
		assert.NoError(t, err, "Error parsing JSON response for merchant creation")
		assert.Equal(t, merchant.ID, merchantInfo.ID, "Returned user ID should match the registered user's ID")
	})
	t.Run("Merchant doesn't exists", func(t *testing.T) {
		_, accessToken := testConfig.TestFactory.CreateUser()
		merchantResponse := test.PerformRequest(router, http.MethodGet, "/frontend/api/v1/merchant/me/", nil, accessToken)
		assert.Equal(t, http.StatusNotFound, merchantResponse.Code, "Expected status 200 for merchant creation")
	})
}

func TestMerchantTokenCreate(t *testing.T) {
	const testUrl = "/frontend/api/v1/merchant/me/tokens/"
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine
	t.Run("Positive merchant token creation", func(t *testing.T) {
		user, accessToken := testConfig.TestFactory.CreateUser()
		merchant := testConfig.TestFactory.CreateMerchant(user.ID.String())
		chainList := testConfig.TestFactory.CreateBlockchainList()
		tokenList := testConfig.TestFactory.CreateTokenList(chainList)
		createPayload := restdto.MerchantTokenCreateRequestDTO{
			TokenID: tokenList[0].ID,
			Active:  true,
		}
		createPayloadBytes, err := json.Marshal(createPayload)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")

		createResponse := test.PerformRequest(router, http.MethodPost, testUrl, bytes.NewReader(createPayloadBytes), accessToken)
		var createdObj restdto.MerchantTokenCreateResponseDTO
		err = json.Unmarshal(createResponse.Body.Bytes(), &createdObj)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")
		assert.Equal(t, createPayload.TokenID, createdObj.TokenID, "Returned user ID should match the registered user's ID")
		assert.Equal(t, merchant.ID, createdObj.MerchantID, "Returned user ID should match the registered user's ID")
	})
	t.Run("Negative merchant token creation", func(t *testing.T) {
		user, accessToken := testConfig.TestFactory.CreateUser()
		merchant := testConfig.TestFactory.CreateMerchant(user.ID.String())
		chainList := testConfig.TestFactory.CreateBlockchainList()
		tokenList := testConfig.TestFactory.CreateTokenList(chainList)
		createPayload := restdto.MerchantTokenCreateRequestDTO{
			TokenID: tokenList[0].ID,
			Active:  true,
		}
		createPayloadBytes, err := json.Marshal(createPayload)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")

		createResponse := test.PerformRequest(router, http.MethodPost, testUrl, bytes.NewReader(createPayloadBytes), accessToken)
		var createdObj restdto.MerchantTokenCreateResponseDTO
		err = json.Unmarshal(createResponse.Body.Bytes(), &createdObj)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")
		assert.Equal(t, createPayload.TokenID, createdObj.TokenID, "Returned user ID should match the registered user's ID")
		assert.Equal(t, merchant.ID, createdObj.MerchantID, "Returned user ID should match the registered user's ID")

		createResponse = test.PerformRequest(router, http.MethodPost, testUrl, bytes.NewReader(createPayloadBytes), accessToken)
		assert.NoError(t, err, "Error marshaling request body for first merchant creation")
		assert.Equal(t, http.StatusBadRequest, createResponse.Code, "Returned user ID should match the registered user's ID")
	})
}
