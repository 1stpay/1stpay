# 1stPay SDK Examples

## Overview

This document provides SDK examples and code snippets for integrating with the 1stPay API in various programming languages.

## JavaScript/Node.js SDK

### Installation

```bash
npm install axios uuid
```

### Basic SDK Implementation

```javascript
const axios = require('axios');

class FirstPaySDK {
  constructor(baseURL, accessToken = null) {
    this.baseURL = baseURL;
    this.accessToken = accessToken;
    this.client = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // Add auth interceptor
    this.client.interceptors.request.use((config) => {
      if (this.accessToken) {
        config.headers.Authorization = `Bearer ${this.accessToken}`;
      }
      return config;
    });

    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          throw new Error(error.response.data.error || 'API Error');
        }
        throw error;
      }
    );
  }

  // Authentication methods
  async register(email, password) {
    const response = await this.client.post('/merchant/api/v1/auth/register/', {
      email,
      password
    });
    this.accessToken = response.data.access_token;
    return response.data;
  }

  async login(email, password) {
    const response = await this.client.post('/merchant/api/v1/auth/login/', {
      email,
      password
    });
    this.accessToken = response.data.access_token;
    return response.data;
  }

  // User methods
  async getUserProfile() {
    const response = await this.client.get('/merchant/api/v1/user/me/');
    return response.data;
  }

  // Merchant methods
  async createMerchant(name) {
    const response = await this.client.post('/merchant/api/v1/merchant/', {
      name
    });
    return response.data;
  }

  async getMerchant() {
    const response = await this.client.get('/merchant/api/v1/merchant/me/');
    return response.data;
  }

  async updateMerchant(name) {
    const response = await this.client.put('/merchant/api/v1/merchant/me/', {
      name
    });
    return response.data;
  }

  // Token methods
  async listTokens() {
    const response = await this.client.get('/merchant/api/v1/token/list/');
    return response.data;
  }

  async getMerchantTokens() {
    const response = await this.client.get('/merchant/api/v1/merchant/me/tokens/');
    return response.data;
  }

  async addMerchantToken(tokenId, active = true) {
    const response = await this.client.post('/merchant/api/v1/merchant/me/tokens/', {
      token_id: tokenId,
      active
    });
    return response.data;
  }

  // Blockchain methods
  async listBlockchains() {
    const response = await this.client.get('/merchant/api/v1/blockchain/list/');
    return response.data;
  }

  // Payment methods
  async createPayment(amount, email = null) {
    const response = await this.client.post('/merchant/api/v1/payment/', {
      requested_amount: amount,
      email
    });
    return response.data;
  }

  async getPaymentDetails(paymentId) {
    const response = await this.client.get(`/frontend/api/v1/payment/${paymentId}/`);
    return response.data;
  }

  // Health check
  async ping() {
    const response = await this.client.get('/merchant/api/v1/ping');
    return response.data;
  }
}

module.exports = FirstPaySDK;
```

### Usage Examples

```javascript
const FirstPaySDK = require('./firstpay-sdk');

async function example() {
  // Initialize SDK
  const sdk = new FirstPaySDK('http://localhost:8080');

  try {
    // Register and authenticate
    await sdk.register('merchant@example.com', 'securepassword123');
    console.log('Registered successfully');

    // Create merchant profile
    const merchant = await sdk.createMerchant('My Crypto Store');
    console.log('Merchant created:', merchant);

    // Get available tokens
    const tokens = await sdk.listTokens();
    console.log('Available tokens:', tokens);

    // Configure accepted tokens
    if (tokens.length > 0) {
      const tokenConfig = await sdk.addMerchantToken(tokens[0].id);
      console.log('Token configured:', tokenConfig);
    }

    // Create a payment
    const payment = await sdk.createPayment(100.50, 'customer@example.com');
    console.log('Payment created:', payment);

    // Get payment details (customer view)
    const paymentDetails = await sdk.getPaymentDetails(payment.id);
    console.log('Payment details:', paymentDetails);

  } catch (error) {
    console.error('Error:', error.message);
  }
}

example();
```

## Python SDK

### Installation

```bash
pip install requests
```

### Basic SDK Implementation

```python
import requests
import json
from typing import Optional, Dict, Any, List

class FirstPaySDK:
    def __init__(self, base_url: str, access_token: Optional[str] = None):
        self.base_url = base_url
        self.access_token = access_token
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def _make_request(self, method: str, endpoint: str, data: Optional[Dict] = None, auth_required: bool = True) -> Dict[str, Any]:
        url = f"{self.base_url}{endpoint}"
        headers = {}
        
        if auth_required and self.access_token:
            headers['Authorization'] = f'Bearer {self.access_token}'
        
        try:
            response = self.session.request(
                method=method,
                url=url,
                json=data,
                headers=headers
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.HTTPError as e:
            error_msg = e.response.json().get('error', 'API Error') if e.response.content else str(e)
            raise Exception(f"API Error: {error_msg}")
    
    # Authentication methods
    def register(self, email: str, password: str) -> Dict[str, Any]:
        response = self._make_request('POST', '/merchant/api/v1/auth/register/', {
            'email': email,
            'password': password
        }, auth_required=False)
        self.access_token = response['access_token']
        return response
    
    def login(self, email: str, password: str) -> Dict[str, Any]:
        response = self._make_request('POST', '/merchant/api/v1/auth/login/', {
            'email': email,
            'password': password
        }, auth_required=False)
        self.access_token = response['access_token']
        return response
    
    # User methods
    def get_user_profile(self) -> Dict[str, Any]:
        return self._make_request('GET', '/merchant/api/v1/user/me/')
    
    # Merchant methods
    def create_merchant(self, name: str) -> Dict[str, Any]:
        return self._make_request('POST', '/merchant/api/v1/merchant/', {
            'name': name
        })
    
    def get_merchant(self) -> Dict[str, Any]:
        return self._make_request('GET', '/merchant/api/v1/merchant/me/')
    
    def update_merchant(self, name: str) -> Dict[str, Any]:
        return self._make_request('PUT', '/merchant/api/v1/merchant/me/', {
            'name': name
        })
    
    # Token methods
    def list_tokens(self) -> List[Dict[str, Any]]:
        return self._make_request('GET', '/merchant/api/v1/token/list/', auth_required=False)
    
    def get_merchant_tokens(self) -> List[Dict[str, Any]]:
        return self._make_request('GET', '/merchant/api/v1/merchant/me/tokens/')
    
    def add_merchant_token(self, token_id: str, active: bool = True) -> Dict[str, Any]:
        return self._make_request('POST', '/merchant/api/v1/merchant/me/tokens/', {
            'token_id': token_id,
            'active': active
        })
    
    # Blockchain methods
    def list_blockchains(self) -> List[Dict[str, Any]]:
        return self._make_request('GET', '/merchant/api/v1/blockchain/list/', auth_required=False)
    
    # Payment methods
    def create_payment(self, amount: float, email: Optional[str] = None) -> Dict[str, Any]:
        data = {'requested_amount': amount}
        if email:
            data['email'] = email
        return self._make_request('POST', '/merchant/api/v1/payment/', data)
    
    def get_payment_details(self, payment_id: str) -> Dict[str, Any]:
        return self._make_request('GET', f'/frontend/api/v1/payment/{payment_id}/', auth_required=False)
    
    # Health check
    def ping(self) -> Dict[str, Any]:
        return self._make_request('GET', '/merchant/api/v1/ping', auth_required=False)
```

### Usage Examples

```python
from firstpay_sdk import FirstPaySDK

def main():
    # Initialize SDK
    sdk = FirstPaySDK('http://localhost:8080')
    
    try:
        # Register and authenticate
        auth_response = sdk.register('merchant@example.com', 'securepassword123')
        print('Registered successfully')
        
        # Create merchant profile
        merchant = sdk.create_merchant('My Crypto Store')
        print(f'Merchant created: {merchant}')
        
        # Get available tokens
        tokens = sdk.list_tokens()
        print(f'Available tokens: {len(tokens)}')
        
        # Configure accepted tokens
        if tokens:
            token_config = sdk.add_merchant_token(tokens[0]['id'])
            print(f'Token configured: {token_config}')
        
        # Create a payment
        payment = sdk.create_payment(100.50, 'customer@example.com')
        print(f'Payment created: {payment}')
        
        # Get payment details (customer view)
        payment_details = sdk.get_payment_details(payment['id'])
        print(f'Payment details: {payment_details}')
        
    except Exception as error:
        print(f'Error: {error}')

if __name__ == '__main__':
    main()
```

## PHP SDK

### Basic SDK Implementation

```php
<?php

class FirstPaySDK {
    private $baseURL;
    private $accessToken;
    private $httpClient;

    public function __construct($baseURL, $accessToken = null) {
        $this->baseURL = $baseURL;
        $this->accessToken = $accessToken;
        $this->httpClient = new \GuzzleHttp\Client([
            'base_uri' => $this->baseURL,
            'headers' => [
                'Content-Type' => 'application/json'
            ]
        ]);
    }

    private function makeRequest($method, $endpoint, $data = null, $authRequired = true) {
        $options = [];
        
        if ($authRequired && $this->accessToken) {
            $options['headers']['Authorization'] = 'Bearer ' . $this->accessToken;
        }
        
        if ($data) {
            $options['json'] = $data;
        }

        try {
            $response = $this->httpClient->request($method, $endpoint, $options);
            return json_decode($response->getBody(), true);
        } catch (\GuzzleHttp\Exception\ClientException $e) {
            $errorResponse = json_decode($e->getResponse()->getBody(), true);
            throw new Exception($errorResponse['error'] ?? 'API Error');
        }
    }

    // Authentication methods
    public function register($email, $password) {
        $response = $this->makeRequest('POST', '/merchant/api/v1/auth/register/', [
            'email' => $email,
            'password' => $password
        ], false);
        $this->accessToken = $response['access_token'];
        return $response;
    }

    public function login($email, $password) {
        $response = $this->makeRequest('POST', '/merchant/api/v1/auth/login/', [
            'email' => $email,
            'password' => $password
        ], false);
        $this->accessToken = $response['access_token'];
        return $response;
    }

    // Merchant methods
    public function createMerchant($name) {
        return $this->makeRequest('POST', '/merchant/api/v1/merchant/', [
            'name' => $name
        ]);
    }

    public function getMerchant() {
        return $this->makeRequest('GET', '/merchant/api/v1/merchant/me/');
    }

    // Payment methods
    public function createPayment($amount, $email = null) {
        $data = ['requested_amount' => $amount];
        if ($email) {
            $data['email'] = $email;
        }
        return $this->makeRequest('POST', '/merchant/api/v1/payment/', $data);
    }

    public function getPaymentDetails($paymentId) {
        return $this->makeRequest('GET', "/frontend/api/v1/payment/{$paymentId}/", null, false);
    }

    // Token methods
    public function listTokens() {
        return $this->makeRequest('GET', '/merchant/api/v1/token/list/', null, false);
    }

    public function addMerchantToken($tokenId, $active = true) {
        return $this->makeRequest('POST', '/merchant/api/v1/merchant/me/tokens/', [
            'token_id' => $tokenId,
            'active' => $active
        ]);
    }
}
```

### Usage Example

```php
<?php
require_once 'vendor/autoload.php';

$sdk = new FirstPaySDK('http://localhost:8080');

try {
    // Register and authenticate
    $authResponse = $sdk->register('merchant@example.com', 'securepassword123');
    echo "Registered successfully\n";
    
    // Create merchant
    $merchant = $sdk->createMerchant('My PHP Store');
    echo "Merchant created: " . $merchant['name'] . "\n";
    
    // Create payment
    $payment = $sdk->createPayment(100.50, 'customer@example.com');
    echo "Payment created: " . $payment['id'] . "\n";
    
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
?>
```

## Go SDK

### Basic SDK Implementation

```go
package firstpay

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type SDK struct {
    BaseURL     string
    AccessToken string
    HTTPClient  *http.Client
}

type APIError struct {
    Error string `json:"error"`
}

func NewSDK(baseURL string, accessToken ...string) *SDK {
    sdk := &SDK{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
    if len(accessToken) > 0 {
        sdk.AccessToken = accessToken[0]
    }
    return sdk
}

func (s *SDK) makeRequest(method, endpoint string, body interface{}, authRequired bool) ([]byte, error) {
    var reqBody io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, s.BaseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if authRequired && s.AccessToken != "" {
        req.Header.Set("Authorization", "Bearer "+s.AccessToken)
    }

    resp, err := s.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode >= 400 {
        var apiErr APIError
        if err := json.Unmarshal(respBody, &apiErr); err == nil {
            return nil, fmt.Errorf("API error: %s", apiErr.Error)
        }
        return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
    }

    return respBody, nil
}

// Authentication methods
type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AccessTokenResponse struct {
    AccessToken string `json:"access_token"`
}

func (s *SDK) Register(email, password string) (*AccessTokenResponse, error) {
    req := RegisterRequest{Email: email, Password: password}
    respBody, err := s.makeRequest("POST", "/merchant/api/v1/auth/register/", req, false)
    if err != nil {
        return nil, err
    }

    var response AccessTokenResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    s.AccessToken = response.AccessToken
    return &response, nil
}

func (s *SDK) Login(email, password string) (*AccessTokenResponse, error) {
    req := RegisterRequest{Email: email, Password: password}
    respBody, err := s.makeRequest("POST", "/merchant/api/v1/auth/login/", req, false)
    if err != nil {
        return nil, err
    }

    var response AccessTokenResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    s.AccessToken = response.AccessToken
    return &response, nil
}

// Merchant methods
type MerchantCreateRequest struct {
    Name string `json:"name"`
}

type MerchantResponse struct {
    ID             string    `json:"id"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    UserID         string    `json:"user_id"`
    Name           string    `json:"name"`
    CommissionRate float64   `json:"commission_rate"`
}

func (s *SDK) CreateMerchant(name string) (*MerchantResponse, error) {
    req := MerchantCreateRequest{Name: name}
    respBody, err := s.makeRequest("POST", "/merchant/api/v1/merchant/", req, true)
    if err != nil {
        return nil, err
    }

    var response MerchantResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    return &response, nil
}

// Payment methods
type PaymentCreateRequest struct {
    RequestedAmount float64 `json:"requested_amount"`
    Email          *string `json:"email,omitempty"`
}

type PaymentCreateResponse struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (s *SDK) CreatePayment(amount float64, email *string) (*PaymentCreateResponse, error) {
    req := PaymentCreateRequest{RequestedAmount: amount, Email: email}
    respBody, err := s.makeRequest("POST", "/merchant/api/v1/payment/", req, true)
    if err != nil {
        return nil, err
    }

    var response PaymentCreateResponse
    if err := json.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }

    return &response, nil
}
```

### Usage Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/yourorg/firstpay-go-sdk"
)

func main() {
    sdk := firstpay.NewSDK("http://localhost:8080")

    // Register
    authResp, err := sdk.Register("merchant@example.com", "securepassword123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Registered successfully")

    // Create merchant
    merchant, err := sdk.CreateMerchant("My Go Store")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Merchant created: %s\n", merchant.Name)

    // Create payment
    email := "customer@example.com"
    payment, err := sdk.CreatePayment(100.50, &email)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Payment created: %s\n", payment.ID)
}
```

## cURL Examples

### Complete Workflow

```bash
#!/bin/bash

# Configuration
API_BASE="http://localhost:8080"
EMAIL="merchant@example.com"
PASSWORD="securepassword123"

# 1. Register
echo "Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "${API_BASE}/merchant/api/v1/auth/register/" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}")

ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.access_token')
echo "Access token: ${ACCESS_TOKEN}"

# 2. Create merchant
echo "Creating merchant..."
MERCHANT_RESPONSE=$(curl -s -X POST "${API_BASE}/merchant/api/v1/merchant/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"name":"My Bash Store"}')

echo "Merchant created: $(echo $MERCHANT_RESPONSE | jq -r '.name')"

# 3. Get available tokens
echo "Getting available tokens..."
TOKENS_RESPONSE=$(curl -s -X GET "${API_BASE}/merchant/api/v1/token/list/")
FIRST_TOKEN_ID=$(echo $TOKENS_RESPONSE | jq -r '.[0].id')
echo "First token ID: ${FIRST_TOKEN_ID}"

# 4. Configure token
echo "Configuring token..."
TOKEN_CONFIG_RESPONSE=$(curl -s -X POST "${API_BASE}/merchant/api/v1/merchant/me/tokens/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d "{\"token_id\":\"${FIRST_TOKEN_ID}\",\"active\":true}")

echo "Token configured"

# 5. Create payment
echo "Creating payment..."
PAYMENT_RESPONSE=$(curl -s -X POST "${API_BASE}/merchant/api/v1/payment/" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"requested_amount":100.50,"email":"customer@example.com"}')

PAYMENT_ID=$(echo $PAYMENT_RESPONSE | jq -r '.id')
echo "Payment created: ${PAYMENT_ID}"

# 6. Get payment details (customer view)
echo "Getting payment details..."
PAYMENT_DETAILS=$(curl -s -X GET "${API_BASE}/frontend/api/v1/payment/${PAYMENT_ID}/")
echo "Payment details: $(echo $PAYMENT_DETAILS | jq '.')"
```

## React.js Integration

### API Service

```javascript
// services/firstPayAPI.js
import axios from 'axios';

class FirstPayAPI {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.client = axios.create({
      baseURL: this.baseURL,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    // Add auth interceptor
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('firstpay_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  async register(email, password) {
    const response = await this.client.post('/merchant/api/v1/auth/register/', {
      email,
      password
    });
    localStorage.setItem('firstpay_token', response.data.access_token);
    return response.data;
  }

  async createPayment(amount, email) {
    const response = await this.client.post('/merchant/api/v1/payment/', {
      requested_amount: amount,
      email
    });
    return response.data;
  }

  async getPaymentDetails(paymentId) {
    const response = await this.client.get(`/frontend/api/v1/payment/${paymentId}/`);
    return response.data;
  }
}

export default new FirstPayAPI('http://localhost:8080');
```

### React Component

```jsx
// components/PaymentForm.jsx
import React, { useState } from 'react';
import firstPayAPI from '../services/firstPayAPI';

const PaymentForm = () => {
  const [amount, setAmount] = useState('');
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [paymentId, setPaymentId] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      const payment = await firstPayAPI.createPayment(
        parseFloat(amount),
        email || null
      );
      setPaymentId(payment.id);
      alert('Payment created successfully!');
    } catch (error) {
      alert('Error creating payment: ' + error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="payment-form">
      <h2>Create Payment</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Amount:</label>
          <input
            type="number"
            step="0.01"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Customer Email (optional):</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Create Payment'}
        </button>
      </form>

      {paymentId && (
        <div className="payment-result">
          <h3>Payment Created!</h3>
          <p>Payment ID: {paymentId}</p>
          <p>Share this link with your customer:</p>
          <a 
            href={`http://localhost:8080/frontend/api/v1/payment/${paymentId}/`}
            target="_blank"
            rel="noopener noreferrer"
          >
            Payment Link
          </a>
        </div>
      )}
    </div>
  );
};

export default PaymentForm;
```

## Payment Status Monitoring

### JavaScript Polling Example

```javascript
class PaymentMonitor {
  constructor(api, paymentId, onStatusChange) {
    this.api = api;
    this.paymentId = paymentId;
    this.onStatusChange = onStatusChange;
    this.polling = false;
    this.interval = null;
  }

  start(intervalMs = 30000) {
    if (this.polling) return;
    
    this.polling = true;
    this.checkStatus(); // Initial check
    
    this.interval = setInterval(() => {
      this.checkStatus();
    }, intervalMs);
  }

  stop() {
    this.polling = false;
    if (this.interval) {
      clearInterval(this.interval);
      this.interval = null;
    }
  }

  async checkStatus() {
    try {
      const payment = await this.api.getPaymentDetails(this.paymentId);
      
      if (this.onStatusChange) {
        this.onStatusChange(payment.status, payment);
      }

      // Stop monitoring if payment is completed or failed
      if (['completed', 'failed', 'not_filled'].includes(payment.status)) {
        this.stop();
      }
    } catch (error) {
      console.error('Error checking payment status:', error);
    }
  }
}

// Usage
const monitor = new PaymentMonitor(
  firstPayAPI,
  'payment-id-here',
  (status, payment) => {
    console.log(`Payment status changed to: ${status}`);
    
    switch (status) {
      case 'completed':
        alert('Payment completed successfully!');
        break;
      case 'failed':
        alert('Payment failed');
        break;
      case 'not_filled':
        alert('Payment expired');
        break;
    }
  }
);

monitor.start(); // Start monitoring every 30 seconds
```

## Error Handling Patterns

### Comprehensive Error Handling

```javascript
class FirstPayError extends Error {
  constructor(message, statusCode, originalError) {
    super(message);
    this.name = 'FirstPayError';
    this.statusCode = statusCode;
    this.originalError = originalError;
  }
}

function handleFirstPayError(error) {
  if (error.response) {
    const { status, data } = error.response;
    
    switch (status) {
      case 400:
        throw new FirstPayError(
          `Validation Error: ${data.error}`,
          400,
          error
        );
      case 401:
        // Token expired or invalid
        localStorage.removeItem('firstpay_token');
        throw new FirstPayError(
          'Authentication required. Please login again.',
          401,
          error
        );
      case 403:
        throw new FirstPayError(
          `Access Denied: ${data.error}`,
          403,
          error
        );
      case 404:
        throw new FirstPayError(
          `Resource Not Found: ${data.error}`,
          404,
          error
        );
      case 500:
        throw new FirstPayError(
          'Server Error. Please try again later.',
          500,
          error
        );
      default:
        throw new FirstPayError(
          `API Error: ${data.error || 'Unknown error'}`,
          status,
          error
        );
    }
  }
  
  // Network or other errors
  throw new FirstPayError(
    'Network Error. Please check your connection.',
    0,
    error
  );
}
```

## Best Practices

### 1. Security

```javascript
// ✅ Good: Store tokens securely
const secureStorage = {
  setToken(token) {
    // Use secure storage in production
    localStorage.setItem('firstpay_token', token);
  },
  
  getToken() {
    return localStorage.getItem('firstpay_token');
  },
  
  removeToken() {
    localStorage.removeItem('firstpay_token');
  }
};

// ✅ Good: Validate input data
function validatePaymentAmount(amount) {
  if (typeof amount !== 'number' || amount <= 0) {
    throw new Error('Amount must be a positive number');
  }
  if (amount > 1000000) {
    throw new Error('Amount exceeds maximum limit');
  }
  return true;
}
```

### 2. Performance

```javascript
// ✅ Good: Cache reference data
class CachedFirstPayAPI extends FirstPaySDK {
  constructor(baseURL, accessToken) {
    super(baseURL, accessToken);
    this.cache = new Map();
    this.cacheTimeout = 5 * 60 * 1000; // 5 minutes
  }

  async listTokens() {
    const cacheKey = 'tokens';
    const cached = this.cache.get(cacheKey);
    
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data;
    }

    const tokens = await super.listTokens();
    this.cache.set(cacheKey, {
      data: tokens,
      timestamp: Date.now()
    });

    return tokens;
  }
}
```

### 3. Testing

```javascript
// Jest test example
describe('FirstPay SDK', () => {
  let sdk;
  
  beforeEach(() => {
    sdk = new FirstPaySDK('http://localhost:8080');
  });

  test('should register user successfully', async () => {
    const mockResponse = { access_token: 'mock-token' };
    jest.spyOn(sdk.client, 'post').mockResolvedValue({ data: mockResponse });

    const result = await sdk.register('test@example.com', 'password');
    
    expect(result.access_token).toBe('mock-token');
    expect(sdk.accessToken).toBe('mock-token');
  });

  test('should create payment successfully', async () => {
    sdk.accessToken = 'valid-token';
    const mockPayment = {
      id: 'payment-id',
      created_at: '2024-01-15T10:30:00Z',
      updated_at: '2024-01-15T10:30:00Z'
    };
    
    jest.spyOn(sdk.client, 'post').mockResolvedValue({ data: mockPayment });

    const result = await sdk.createPayment(100.50, 'customer@example.com');
    
    expect(result.id).toBe('payment-id');
  });
});
```

## Production Considerations

### 1. Environment Configuration

```javascript
const config = {
  development: {
    apiUrl: 'http://localhost:8080',
    debug: true
  },
  production: {
    apiUrl: 'https://api.1stpay.com',
    debug: false
  }
};

const env = process.env.NODE_ENV || 'development';
const apiConfig = config[env];
```

### 2. Monitoring and Logging

```javascript
class MonitoredFirstPaySDK extends FirstPaySDK {
  constructor(baseURL, accessToken, logger) {
    super(baseURL, accessToken);
    this.logger = logger;
    
    // Add request/response logging
    this.client.interceptors.request.use((config) => {
      this.logger.info('API Request', {
        method: config.method,
        url: config.url,
        timestamp: new Date().toISOString()
      });
      return config;
    });

    this.client.interceptors.response.use(
      (response) => {
        this.logger.info('API Response', {
          status: response.status,
          url: response.config.url,
          timestamp: new Date().toISOString()
        });
        return response;
      },
      (error) => {
        this.logger.error('API Error', {
          status: error.response?.status,
          url: error.config?.url,
          error: error.message,
          timestamp: new Date().toISOString()
        });
        return Promise.reject(error);
      }
    );
  }
}
```

---

These SDK examples provide a solid foundation for integrating with the 1stPay API across different programming languages and platforms. Choose the example that best fits your technology stack and customize as needed.