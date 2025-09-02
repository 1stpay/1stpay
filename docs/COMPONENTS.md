# 1stPay Components Documentation

## Overview

This document provides detailed information about the internal components and architecture of the 1stPay system.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend UI   │    │   Merchant UI   │    │  Integration    │
│                 │    │                 │    │    Systems      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │                      │                      │
    ┌─────▼──────────────────────▼──────────────────────▼─────┐
    │                REST API Layer                           │
    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐   │
    │  │  Frontend   │ │  Merchant   │ │  Integration    │   │
    │  │     API     │ │     API     │ │      API        │   │
    │  └─────────────┘ └─────────────┘ └─────────────────┘   │
    └─────────────────────────┬───────────────────────────────┘
                              │
    ┌─────────────────────────▼───────────────────────────────┐
    │                Application Layer                        │
    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐   │
    │  │  Use Cases  │ │  Services   │ │   Controllers   │   │
    │  └─────────────┘ └─────────────┘ └─────────────────┘   │
    └─────────────────────────┬───────────────────────────────┘
                              │
    ┌─────────────────────────▼───────────────────────────────┐
    │                Domain Layer                             │
    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐   │
    │  │   Models    │ │    Enums    │ │  Repositories   │   │
    │  └─────────────┘ └─────────────┘ └─────────────────┘   │
    └─────────────────────────┬───────────────────────────────┘
                              │
    ┌─────────────────────────▼───────────────────────────────┐
    │              Infrastructure Layer                       │
    │  ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐   │
    │  │  Database   │ │ Blockchain  │ │   External      │   │
    │  │             │ │ Integrations│ │   Services      │   │
    │  └─────────────┘ └─────────────┘ └─────────────────┘   │
    └─────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Transport Layer

#### REST API Routes (`internal/transport/rest/`)

The REST layer is organized into three main API groups:

##### Merchant API (`/merchant/api/v1`)
- **Purpose**: Merchant operations and management
- **Authentication**: Required for most endpoints
- **Features**: User management, merchant profiles, payment creation, token configuration

##### Frontend API (`/frontend/api/v1`)
- **Purpose**: Customer-facing payment interfaces
- **Authentication**: Optional (public payment viewing)
- **Features**: Payment details, customer payment interface

##### Integration API (`/integration/api/v1`)
- **Purpose**: External system integrations
- **Authentication**: Varies by endpoint
- **Features**: Webhook handling, third-party integrations

#### Controllers

Controllers handle HTTP requests and coordinate between the transport layer and business logic.

**Key Controllers:**
- `AuthController`: User authentication and authorization
- `MerchantController`: Merchant account management
- `PaymentController`: Payment creation and management
- `UserController`: User profile management
- `TokenController`: Token listing and management
- `BlockchainController`: Blockchain network information

#### Middleware

**JWT Authentication Middleware:**
```go
// Protects endpoints requiring authentication
protectedRouter.Use(deps.Middleware.JWTAuth)
```

**CORS Middleware:**
```go
// Enables cross-origin requests
gin.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

### 2. Application Layer

#### Use Cases (`internal/domain/usecase/`)

Use cases contain the business logic and coordinate between different domain services.

##### AuthUsecase
- **Purpose**: User authentication and authorization
- **Methods**:
  - `Register(req RegisterRequest) (User, string, error)`: Register new user
  - `Login(req LoginRequest) (User, string, error)`: Authenticate user

##### MerchantUsecase
- **Purpose**: Merchant management operations
- **Methods**:
  - `CreateMerchant(req MerchantCreateRequest, userID string) (Merchant, error)`
  - `GetMerchantByUserId(userID string) (Merchant, error)`
  - `UpdateMerchant(req MerchantCreateRequest, userID string) (Merchant, error)`
  - `ListMerchantToken(merchantID string) ([]MerchantToken, error)`
  - `CreateMerchantToken(req MerchantTokenCreateRequest, merchantID string) (MerchantToken, error)`

##### PaymentUsecase
- **Purpose**: Payment processing and management
- **Methods**:
  - `CreatePaymentWithWallets(paymentData InvoiceCreateRestDTO, merchantId uuid.UUID) (Payment, error)`
  - `GetPaymentWithAddresses(paymentID string) (Payment, []PaymentAddress, error)`

##### UserUsecase
- **Purpose**: User profile management
- **Methods**:
  - `GetUserById(userID string) (User, error)`

##### TokenUsecase
- **Purpose**: Token management
- **Methods**:
  - `ListActive() ([]Token, error)`: Get all active tokens

##### BlockchainUsecase
- **Purpose**: Blockchain network management
- **Methods**:
  - `ListActive() ([]Blockchain, error)`: Get all active blockchains

### 3. Domain Layer

#### Models (`internal/model/`)

Domain models represent the core business entities:

##### User Model
- Manages user accounts and authentication
- Links to merchant accounts
- Stores encrypted passwords and profile information

##### Merchant Model
- Represents business accounts
- Links to user accounts
- Manages commission rates and business settings

##### Payment Model
- Tracks payment transactions
- Manages payment status and amounts
- Links to merchants and tokens
- Supports AML checking

##### Token Model
- Represents cryptocurrency tokens
- Links to blockchain networks
- Manages token metadata and configuration

##### Blockchain Model
- Represents supported blockchain networks
- Manages network configuration
- Supports multiple network types (EVM, Tron, TON, Solana)

#### Enums (`internal/domain/enum/`)

##### PaymentStatus
- `pending`: Payment awaiting completion
- `completed`: Payment successfully processed
- `failed`: Payment processing failed
- `not_filled`: Payment expired without completion

##### PaymentAMLStatus
- `passed`: AML verification successful
- `failed`: AML verification failed
- `pending`: AML verification in progress

##### NetworkType
- `evm`: Ethereum Virtual Machine compatible networks
- `tron`: Tron blockchain network
- `ton`: TON blockchain network
- `solana`: Solana blockchain network

### 4. Infrastructure Layer

#### Database (`internal/repository/`)

The system uses PostgreSQL as the primary database with GORM as the ORM.

**Key Features:**
- UUID primary keys for all entities
- Automatic timestamp management
- Foreign key relationships
- JSON configuration fields
- Database migrations support

#### Configuration (`internal/config/`)

##### Application Configuration
```go
type Application struct {
    Env      *Env
    Postgres *gorm.DB
    Deps     *Dependencies
}
```

##### Environment Configuration
```go
type Env struct {
    HttpPort string
    // Database configuration
    // JWT configuration
    // Blockchain RPC URLs
    // External service configurations
}
```

##### Dependencies
```go
type Dependencies struct {
    Controllers *Controllers
    Middleware  *Middleware
}
```

## Component Interactions

### 1. Request Flow

```
HTTP Request → Router → Middleware → Controller → UseCase → Repository → Database
                ↓
HTTP Response ← DTO ← Controller ← UseCase ← Model ← Repository ← Database
```

### 2. Authentication Flow

```
Client → AuthController → AuthUsecase → UserRepository → Database
                                    ↓
JWT Token ← AuthController ← JWT Service ← AuthUsecase
```

### 3. Payment Creation Flow

```
Merchant → PaymentController → PaymentUsecase → PaymentRepository → Database
                                           ↓
Blockchain Service → Wallet Generation → PaymentAddress Creation
                                           ↓
Payment Response ← PaymentController ← Payment + Addresses
```

## Extension Points

### 1. Adding New Blockchain Support

To add support for a new blockchain:

1. **Add Network Type**: Update `NetworkType` enum
2. **Implement Blockchain Service**: Create blockchain-specific integration
3. **Update Token Model**: Add blockchain-specific configuration
4. **Add Migration**: Create database migration for new network type
5. **Update Controllers**: Ensure API endpoints support new blockchain

### 2. Adding New Payment Methods

To add new payment methods:

1. **Extend Payment Model**: Add new payment method fields
2. **Update Payment Status**: Add new status types if needed
3. **Implement Payment Service**: Create payment method-specific logic
4. **Update DTOs**: Add request/response structures
5. **Update Controllers**: Add new endpoints if needed

### 3. Adding Webhook Support

To implement webhooks:

1. **Create Webhook Model**: Store webhook configurations
2. **Implement Webhook Service**: Handle webhook delivery
3. **Add Webhook Endpoints**: Create webhook management API
4. **Update Payment Flow**: Trigger webhooks on status changes

## Performance Considerations

### 1. Database Optimization

- Use database indexes for frequently queried fields
- Implement connection pooling
- Use prepared statements for repeated queries
- Consider read replicas for heavy read workloads

### 2. API Optimization

- Implement response caching for reference data
- Use pagination for large result sets
- Implement request/response compression
- Consider API rate limiting

### 3. Blockchain Integration

- Use connection pooling for blockchain RPC calls
- Implement retry logic for failed blockchain calls
- Cache blockchain data when possible
- Use webhooks instead of polling when available

## Security Components

### 1. JWT Authentication

- Uses industry-standard JWT tokens
- Tokens include user and merchant information
- Configurable token expiration
- Secure token signing and verification

### 2. Input Validation

- Request body validation using struct tags
- Email format validation
- UUID format validation
- Required field validation

### 3. Database Security

- Encrypted password storage
- SQL injection prevention via ORM
- Database connection encryption
- Audit logging for sensitive operations

## Testing Components

### 1. Unit Tests

Located in `test/` directory:
- Model validation tests
- Use case logic tests
- Controller endpoint tests
- Service integration tests

### 2. Integration Tests

- End-to-end API flow tests
- Database integration tests
- Blockchain integration tests
- Authentication flow tests

## Deployment Components

### 1. Application Startup

The application starts with:
```go
func main() {
    app.Run()
}
```

Which initializes:
- Configuration loading
- Database connections
- Route setup
- Middleware configuration
- Server startup

### 2. Configuration Management

- Environment-based configuration
- Viper for configuration management
- Support for multiple configuration sources
- Runtime configuration validation

### 3. Database Migrations

Located in `migrations/` directory:
- Schema versioning
- Automatic migration execution
- Rollback support
- Data seeding scripts

---

This component documentation provides a comprehensive overview of the 1stPay system architecture and can be used by developers to understand the codebase structure and extend functionality.