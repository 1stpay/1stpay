package route_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/1stpay/1stpay/test"
	"github.com/stretchr/testify/assert"
)

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

func TestRegisterRoute(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine

	t.Run("Positive registration", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "secret",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Error marshaling request body")
		response := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusCreated, response.Code, "Expected status 201 on registration")

		var registerResp RegisterResponse
		err = json.Unmarshal(response.Body.Bytes(), &registerResp)
		assert.NoError(t, err, "Error parsing JSON response")
		assert.NotEmpty(t, registerResp.AccessToken, "Access token should not be empty")
	})

	t.Run("Missing email", func(t *testing.T) {
		reqBody := map[string]string{
			"password": "secret",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Error marshaling request body for missing email")
		response := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusBadRequest, response.Code, "Expected status 400 for missing email")
	})

	t.Run("Missing password", func(t *testing.T) {
		reqBody := map[string]string{
			"email": "test2@example.com",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Error marshaling request body for missing password")
		response := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusBadRequest, response.Code, "Expected status 400 for missing password")
	})

	t.Run("Invalid email format", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "not-an-email",
			"password": "secret",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Error marshaling request body for invalid email format")
		response := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusBadRequest, response.Code, "Expected status 400 for invalid email format")
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "duplicate@example.com",
			"password": "secret",
		}
		reqBodyBytes, err := json.Marshal(reqBody)
		assert.NoError(t, err, "Error marshaling request body for duplicate registration")

		// First registration attempt
		response1 := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusCreated, response1.Code, "Expected status 201 on first registration")

		// Second registration attempt with the same email should fail
		response2 := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(reqBodyBytes))
		assert.Equal(t, http.StatusBadRequest, response2.Code, "Expected status 400 on duplicate registration")
	})
}

func TestLoginRoute(t *testing.T) {
	testConfig := test.NewIntegrationTest(t, "../../../../../")
	router := testConfig.GinEngine

	// Pre-register a user for login tests
	regReqBody := map[string]string{
		"email":    "login-test@example.com",
		"password": "correctpassword",
	}
	regBodyBytes, err := json.Marshal(regReqBody)
	assert.NoError(t, err, "Error marshaling registration request body")
	regResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/register/", bytes.NewReader(regBodyBytes))
	assert.Equal(t, http.StatusCreated, regResponse.Code, "Expected status 201 on registration for login tests")

	t.Run("Positive login", func(t *testing.T) {
		loginReqBody := map[string]string{
			"email":    "login-test@example.com",
			"password": "correctpassword",
		}
		loginBodyBytes, err := json.Marshal(loginReqBody)
		assert.NoError(t, err, "Error marshaling login request body")
		loginResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/login/", bytes.NewReader(loginBodyBytes))
		assert.Equal(t, http.StatusOK, loginResponse.Code, "Expected status 200 on login")

		var loginResp RegisterResponse
		err = json.Unmarshal(loginResponse.Body.Bytes(), &loginResp)
		assert.NoError(t, err, "Error parsing JSON login response")
		assert.NotEmpty(t, loginResp.AccessToken, "Access token should not be empty")
	})

	t.Run("Wrong password", func(t *testing.T) {
		loginReqBody := map[string]string{
			"email":    "login-test@example.com",
			"password": "wrongpassword",
		}
		loginBodyBytes, err := json.Marshal(loginReqBody)
		assert.NoError(t, err, "Error marshaling login request body for wrong password")
		loginResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/login/", bytes.NewReader(loginBodyBytes))
		assert.Equal(t, http.StatusForbidden, loginResponse.Code, "Expected status 401 for wrong password")
	})

	t.Run("Non-existing email", func(t *testing.T) {
		loginReqBody := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "anyPassword",
		}
		loginBodyBytes, err := json.Marshal(loginReqBody)
		assert.NoError(t, err, "Error marshaling login request body for non-existing email")
		loginResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/login/", bytes.NewReader(loginBodyBytes))
		assert.Equal(t, http.StatusForbidden, loginResponse.Code, "Expected status 401 for non-existing email")
	})

	t.Run("Missing email", func(t *testing.T) {
		loginReqBody := map[string]string{
			"password": "correctpassword",
		}
		loginBodyBytes, err := json.Marshal(loginReqBody)
		assert.NoError(t, err, "Error marshaling login request body with missing email")
		loginResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/login/", bytes.NewReader(loginBodyBytes))
		assert.Equal(t, http.StatusBadRequest, loginResponse.Code, "Expected status 400 for missing email")
	})

	t.Run("Missing password", func(t *testing.T) {
		loginReqBody := map[string]string{
			"email": "login-test@example.com",
		}
		loginBodyBytes, err := json.Marshal(loginReqBody)
		assert.NoError(t, err, "Error marshaling login request body with missing password")
		loginResponse := test.PerformRequest(router, http.MethodPost, "/frontend/api/v1/auth/login/", bytes.NewReader(loginBodyBytes))
		assert.Equal(t, http.StatusBadRequest, loginResponse.Code, "Expected status 400 for missing password")
	})
}
