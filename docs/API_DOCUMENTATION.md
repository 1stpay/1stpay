# 1stPay API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Data Models](#data-models)
4. [API Endpoints](#api-endpoints)
   - [Merchant API](#merchant-api)
   - [Frontend API](#frontend-api)
   - [Integration API](#integration-api)
5. [Error Handling](#error-handling)
6. [Examples](#examples)

## Overview

1stPay is a comprehensive cryptocurrency payments service that provides APIs for payment processing, merchant management, and multi-blockchain support. The API is organized into three main groups:

- **Merchant API** (`/merchant/api/v1`): For merchant operations and management
- **Frontend API** (`/frontend/api/v1`): For frontend application interactions
- **Integration API** (`/integration/api/v1`): For external system integrations

## Authentication

The API uses JWT (JSON Web Token) based authentication. Most endpoints require authentication via the `Authorization` header.

### Authentication Flow

1. **Register**: Create a new user account
2. **Login**: Authenticate and receive an access token
3. **Use Token**: Include the token in subsequent requests

```http
Authorization: Bearer <access_token>
```

## Data Models

### User

Represents a user in the system.

```go
type User struct {
    ID        uuid.UUID `json:"id"`         // Unique user identifier
    CreatedAt time.Time `json:"created_at"` // Account creation timestamp
    UpdatedAt time.Time `json:"updated_at"` // Last update timestamp
    Email     string    `json:"email"`      // User email (unique)
    Password  string    `json:"-"`          // Encrypted password (not exposed in API)
    Role      string    `json:"role"`       // User role
}
```

### Merchant

Represents a merchant account linked to a user.

```go
type Merchant struct {
    ID             uuid.UUID `json:"id"`              // Unique merchant identifier
    CreatedAt      time.Time `json:"created_at"`      // Creation timestamp
    UpdatedAt      time.Time `json:"updated_at"`      // Last update timestamp
    UserID         uuid.UUID `json:"user_id"`         // Associated user ID
    Name           string    `json:"name"`            // Merchant name
    CommissionRate float64   `json:"commission_rate"` // Commission rate (0-100)
}
```

### Payment

Represents a payment transaction.

```go
type Payment struct {
    ID               uuid.UUID            `json:"id"`                // Unique payment identifier
    CreatedAt        time.Time            `json:"created_at"`        // Creation timestamp
    UpdatedAt        time.Time            `json:"updated_at"`        // Last update timestamp
    MerchantID       uuid.UUID            `json:"merchant_id"`       // Associated merchant ID
    RequestedAmount  float64              `json:"requested_amount"`  // Amount requested
    PaidAmount       float64              `json:"paid_amount"`       // Amount actually paid
    CommissionAmount float64              `json:"commission_amount"` // Commission charged
    ExpiresAt        *time.Time           `json:"expires_at"`        // Payment expiration (optional)
    AMLStatus        *PaymentAMLStatus    `json:"aml_status"`        // AML check status (optional)
    Status           PaymentStatus        `json:"status"`            // Payment status
    InvoiceEmail     *string              `json:"invoice_email"`     // Invoice email (optional)
    UsedTokenID      *uuid.UUID           `json:"used_token_id"`     // Token used for payment (optional)
}
```

### Token

Represents a cryptocurrency token.

```go
type Token struct {
    ID              uuid.UUID  `json:"id"`               // Unique token identifier
    Name            string     `json:"name"`             // Token name
    Symbol          string     `json:"symbol"`           // Token symbol (e.g., "ETH", "BTC")
    ContractAddress string     `json:"contract_address"` // Smart contract address
    Decimals        int        `json:"decimals"`         // Token decimal places
    Logo            *string    `json:"logo"`             // Logo URL (optional)
    BlockchainID    uuid.UUID  `json:"blockchain_id"`    // Associated blockchain ID
    IsNative        bool       `json:"is_native"`        // Whether token is native to blockchain
    IsActive        bool       `json:"is_active"`        // Whether token is active
    Config          JSON       `json:"config"`           // Additional configuration
}
```

### Blockchain

Represents a supported blockchain network.

```go
type Blockchain struct {
    ID        uuid.UUID   `json:"id"`         // Unique blockchain identifier
    Name      string      `json:"name"`       // Blockchain name
    Logo      *string     `json:"logo"`       // Logo URL (optional)
    IsActive  bool        `json:"is_active"`  // Whether blockchain is active
    ChainType NetworkType `json:"chain_type"` // Type of blockchain network
    Config    JSON        `json:"config"`     // Additional configuration
}
```

### MerchantToken

Represents the relationship between a merchant and supported tokens.

```go
type MerchantToken struct {
    ID         uuid.UUID `json:"id"`          // Unique identifier
    MerchantID uuid.UUID `json:"merchant_id"` // Associated merchant ID
    TokenID    uuid.UUID `json:"token_id"`    // Associated token ID
    Balance    float64   `json:"balance"`     // Current balance
    IsActive   bool      `json:"is_active"`   // Whether token is active for merchant
    CreatedAt  time.Time `json:"created_at"`  // Creation timestamp
}
```

## Enums

### PaymentStatus

```go
type PaymentStatus string

const (
    PaymentStatusPending   PaymentStatus = "pending"     // Payment is pending
    PaymentStatusCompleted PaymentStatus = "completed"   // Payment completed successfully
    PaymentStatusFailed    PaymentStatus = "failed"      // Payment failed
    PaymentStatusNotFilled PaymentStatus = "not_filled"  // Payment not filled within time limit
)
```

### PaymentAMLStatus

```go
type PaymentAMLStatus string

const (
    PaymentAMLStatusPassed  PaymentAMLStatus = "passed"  // AML check passed
    PaymentAMLStatusFailed  PaymentAMLStatus = "failed"  // AML check failed
    PaymentAMLStatusPending PaymentAMLStatus = "pending" // AML check in progress
)
```

### NetworkType

```go
type NetworkType string

const (
    EVM    NetworkType = "evm"    // Ethereum Virtual Machine compatible
    TRON   NetworkType = "tron"   // Tron network
    TON    NetworkType = "ton"    // TON network
    SOLANA NetworkType = "solana" // Solana network
)
```

## API Endpoints

### Merchant API

Base URL: `/merchant/api/v1`

#### Authentication Endpoints

##### POST /auth/register/

Register a new user account.

**Request Body:**
```json
{
    "email": "merchant@example.com",
    "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (400 Bad Request):**
```json
{
    "error": "Invalid input data"
}
```

##### POST /auth/login/

Authenticate an existing user.

**Request Body:**
```json
{
    "email": "merchant@example.com", 
    "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (403 Forbidden):**
```json
{
    "error": "Invalid credentials"
}
```

#### User Profile Endpoints

##### GET /user/me/ 🔒

Get current user profile information.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "merchant@example.com"
}
```

#### Merchant Management Endpoints

##### POST /merchant/ 🔒

Create a new merchant account for the authenticated user.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
    "name": "My Crypto Store"
}
```

**Response (201 Created):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "name": "My Crypto Store",
    "commission_rate": 0.0
}
```

##### GET /merchant/me/ 🔒

Get current merchant details.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "name": "My Crypto Store",
    "commission_rate": 2.5
}
```

**Response (404 Not Found):**
```json
{
    "error": "Merchant not found"
}
```

##### PUT /merchant/me/ 🔒

Update current merchant information.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
    "name": "Updated Store Name"
}
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T15:45:00Z",
    "user_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "name": "Updated Store Name",
    "commission_rate": 2.5
}
```

#### Merchant Token Management

##### GET /merchant/me/tokens/ 🔒

Get list of tokens configured for the merchant.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "merchant_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
        "token_id": "456e7890-e89b-12d3-a456-426614174111",
        "active": true,
        "created_at": "2024-01-15T10:30:00Z"
    }
]
```

##### POST /merchant/me/tokens/ 🔒

Add a new token to the merchant's accepted tokens.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
    "token_id": "456e7890-e89b-12d3-a456-426614174111",
    "active": true
}
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "merchant_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "token_id": "456e7890-e89b-12d3-a456-426614174111",
    "active": true,
    "created_at": "2024-01-15T10:30:00Z"
}
```

#### Payment Management

##### POST /payment/ 🔒

Create a new payment invoice.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
    "requested_amount": 100.50,
    "email": "customer@example.com"
}
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Reference Data Endpoints

##### GET /blockchain/list/

Get list of active blockchain networks.

**Response (200 OK):**
```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Ethereum",
        "logo": "https://example.com/eth-logo.png",
        "is_active": true,
        "chain_type": "evm"
    },
    {
        "id": "456e7890-e89b-12d3-a456-426614174111",
        "name": "Binance Smart Chain",
        "logo": "https://example.com/bsc-logo.png", 
        "is_active": true,
        "chain_type": "evm"
    }
]
```

##### GET /token/list/

Get list of active tokens.

**Response (200 OK):**
```json
[
    {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Ethereum",
        "symbol": "ETH",
        "blockchain_id": "456e7890-e89b-12d3-a456-426614174111",
        "logo": "https://example.com/eth-logo.png",
        "is_native": true,
        "is_active": true
    },
    {
        "id": "789abcde-e89b-12d3-a456-426614174222",
        "name": "USD Coin",
        "symbol": "USDC",
        "blockchain_id": "456e7890-e89b-12d3-a456-426614174111",
        "logo": "https://example.com/usdc-logo.png",
        "is_native": false,
        "is_active": true
    }
]
```

#### Health Check

##### GET /ping

Health check endpoint.

**Response (200 OK):**
```json
{
    "message": "pong"
}
```

### Frontend API

Base URL: `/frontend/api/v1`

The Frontend API provides endpoints for customer-facing payment interfaces and public payment information.

#### Authentication Endpoints

##### POST /auth/register/

Register a new user account (same as merchant API).

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```

**Response (201 Created):**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

##### POST /auth/login/

Authenticate an existing user (same as merchant API).

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "securepassword123"
}
```

**Response (200 OK):**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### User Profile Endpoints

##### GET /user/me/ 🔒

Get current user profile information.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com"
}
```

#### Payment Information

##### GET /payment/:id/

Get detailed payment information for a specific payment ID. This endpoint is typically used by customers to view payment details and make payments.

**Parameters:**
- `id` (path): Payment UUID

**Response (200 OK):**
```json
{
    "requested_amount": 100.50,
    "email": "customer@example.com",
    "expires_at": "2024-01-15T11:30:00Z",
    "merchant": {
        "name": "My Crypto Store"
    },
    "payment_address_list": [
        {
            "public_key": "0x742d35Cc6634C0532925a3b8D4e7B4B9C664E4b1",
            "requested_amount": "100.50000000",
            "token": {
                "name": "Ethereum",
                "symbol": "ETH",
                "logo": "https://example.com/eth-logo.png",
                "blockchain": {
                    "name": "Ethereum",
                    "logo": "https://example.com/eth-logo.png"
                }
            }
        },
        {
            "public_key": "0x8ba1f109551bD432803012645Hac136c",
            "requested_amount": "100.50000000",
            "token": {
                "name": "USD Coin",
                "symbol": "USDC",
                "logo": "https://example.com/usdc-logo.png",
                "blockchain": {
                    "name": "Ethereum",
                    "logo": "https://example.com/eth-logo.png"
                }
            }
        }
    ]
}
```

**Response (404 Not Found):**
```json
{
    "error": "Error"
}
```

#### Reference Data Endpoints

##### GET /blockchain/list/

Get list of active blockchain networks (same as merchant API).

##### GET /token/list/

Get list of active tokens (same as merchant API).

#### Merchant Information

##### POST /merchant/ 🔒

Create merchant account (same as merchant API).

##### GET /merchant/me/ 🔒

Get merchant details (same as merchant API).

##### PUT /merchant/me/ 🔒

Update merchant information (same as merchant API).

##### GET /merchant/me/tokens/ 🔒

Get merchant tokens (same as merchant API).

##### POST /merchant/me/tokens/ 🔒

Add merchant token (same as merchant API).

#### Health Check

##### GET /ping

Health check endpoint.

**Response (200 OK):**
```json
{
    "message": "pong"
}
```

### Integration API

Base URL: `/integration/api/v1`

The Integration API provides endpoints for external system integrations and webhook handling.

#### Payment Integration

##### POST /payment/

Create a payment invoice via integration API.

**Request Body:**
```json
{
    "requested_amount": 100.50,
    "email": "customer@example.com"
}
```

**Response (200 OK):**
```json
{
    "message": "pong"
}
```

##### GET /payment/

Get payment information via integration API.

**Response (200 OK):**
```json
{
    "message": "pong"
}
```

##### GET /payment/:id/

Get specific payment details via integration API.

**Parameters:**
- `id` (path): Payment UUID

**Response (200 OK):**
```json
{
    "message": "pong"
}
```

## Error Handling

### Standard Error Response Format

All API endpoints return errors in a consistent format:

```json
{
    "error": "Error description"
}
```

### HTTP Status Codes

- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data or parameters
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Authentication failed or insufficient permissions
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

### Common Error Messages

- `"Invalid input data"`: Request body validation failed
- `"Wrong input data"`: Request data format is incorrect
- `"Invalid request"`: General request validation error
- `"Merchant not found"`: Merchant resource not found
- `"Invalid user"`: User authentication/context error
- `"Invalid user type"`: User type assertion failed
- `"Invalid credentials"`: Login authentication failed

## Examples

### Complete Payment Flow Example

#### 1. Register a Merchant

```bash
curl -X POST http://localhost:8080/merchant/api/v1/auth/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "email": "merchant@example.com",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. Create Merchant Profile

```bash
curl -X POST http://localhost:8080/merchant/api/v1/merchant/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "My Crypto Store"
  }'
```

**Response:**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "name": "My Crypto Store",
    "commission_rate": 0.0
}
```

#### 3. Get Available Tokens

```bash
curl -X GET http://localhost:8080/merchant/api/v1/token/list/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### 4. Configure Accepted Tokens

```bash
curl -X POST http://localhost:8080/merchant/api/v1/merchant/me/tokens/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "token_id": "456e7890-e89b-12d3-a456-426614174111",
    "active": true
  }'
```

#### 5. Create Payment Invoice

```bash
curl -X POST http://localhost:8080/merchant/api/v1/payment/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "requested_amount": 100.50,
    "email": "customer@example.com"
  }'
```

**Response:**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
}
```

#### 6. Customer Views Payment (Frontend API)

```bash
curl -X GET http://localhost:8080/frontend/api/v1/payment/123e4567-e89b-12d3-a456-426614174000/
```

**Response:**
```json
{
    "requested_amount": 100.50,
    "email": "customer@example.com",
    "expires_at": "2024-01-15T11:30:00Z",
    "merchant": {
        "name": "My Crypto Store"
    },
    "payment_address_list": [
        {
            "public_key": "0x742d35Cc6634C0532925a3b8D4e7B4B9C664E4b1",
            "requested_amount": "100.50000000",
            "token": {
                "name": "Ethereum",
                "symbol": "ETH",
                "logo": "https://example.com/eth-logo.png",
                "blockchain": {
                    "name": "Ethereum",
                    "logo": "https://example.com/eth-logo.png"
                }
            }
        }
    ]
}
```

### Authentication Example

#### Getting User Profile

```bash
curl -X GET http://localhost:8080/merchant/api/v1/user/me/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
    "id": "987fcdeb-51a2-43d1-9c45-123456789abc",
    "email": "merchant@example.com"
}
```

### Error Handling Example

#### Invalid Authentication

```bash
curl -X GET http://localhost:8080/merchant/api/v1/user/me/ \
  -H "Authorization: Bearer invalid_token"
```

**Response (401 Unauthorized):**
```json
{
    "error": "Invalid user"
}
```

#### Missing Required Fields

```bash
curl -X POST http://localhost:8080/merchant/api/v1/auth/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email"
  }'
```

**Response (400 Bad Request):**
```json
{
    "error": "Invalid input data"
}
```

## Usage Instructions

### Setting Up Authentication

1. **Register**: Use the `/auth/register/` endpoint to create a new account
2. **Login**: Use the `/auth/login/` endpoint to get an access token
3. **Include Token**: Add the token to the `Authorization` header for protected endpoints

### Creating a Payment

1. **Setup Merchant**: Register, create merchant profile, configure accepted tokens
2. **Create Invoice**: Use the `/payment/` endpoint to create a payment
3. **Share Payment Link**: Provide the payment ID to customers via the frontend API
4. **Monitor Status**: Track payment status through webhooks or polling

### Integration Guidelines

- **Base URLs**: Use appropriate base URL for your use case (merchant/frontend/integration)
- **Authentication**: Most endpoints require JWT authentication
- **Rate Limiting**: Follow rate limiting guidelines (if implemented)
- **Webhooks**: Configure webhooks for real-time payment status updates
- **Error Handling**: Implement proper error handling for all API responses

### Best Practices

1. **Security**: Always use HTTPS in production
2. **Token Management**: Securely store and refresh access tokens
3. **Error Handling**: Handle all possible error responses
4. **Validation**: Validate all input data before sending requests
5. **Monitoring**: Monitor API usage and payment statuses
6. **Testing**: Use the ping endpoints for health checks

---

**Note**: This documentation is based on the current API implementation. Some integration endpoints may return placeholder responses and require further implementation. The 🔒 symbol indicates endpoints that require authentication.
