# 1stpay

## 1. Project Overview

1stpay is a crypto payments service designed to streamline the acceptance and processing of cryptocurrency transactions. The platform provides a comprehensive solution for businesses and merchants, enabling seamless payment processing, conversion to stablecoins, mass payouts, and in-depth analytics. Its primary goal is to offer a secure, scalable, and versatile payment system that caters to various blockchain networks.

## 2. Key Features

- **Payment Acceptance:**  
  - **Multi-Blockchain Support:**  
    - **EVM-compatible chains:** Ethereum, Binance Smart Chain (BSC), Arbitrum, etc.
    - **Tron**
    - **Ton**
    - **Solana**
  - **Token Support:** Accept payments in various tokens across these blockchains.
  - **Payment Creation:**  
    - Ability to create payments via API.
    - Option to create payments manually.
  - **Webhooks:** Real-time notifications for payment status updates.
  - **AML Integration:** Optional Anti-Money Laundering checks to ensure compliance.

- **Payment Conversion:**  
  - Convert received payments into stablecoins to maintain financial stability.

- **Mass Payouts:**  
  - **Multi-Blockchain Support:**  
    - Support for different blockchains and tokens.
  - **Payout Options:**  
    - Initiate payouts via API.
    - Option to create payouts manually.

- **Analytics:**  
  - Monitor payment volume trends.
  - Track revenue dynamics.
  - Analyze various financial performance indicators.

- **User and Merchant Management:**  
  - User registration and authentication.
  - Merchant account creation and management.

## 3. Architecture and Technologies

- **Backend:**  
  - Developed using **Golang** with the **Gin** framework for robust and efficient API development.

- **Database:**  
  - Uses **PostgreSQL** for reliable and scalable data management.

- **Blockchain Integration:**  
  - Direct support for multiple blockchains including EVM-compatible chains, Tron, Ton, and Solana.

*Note: This list of technologies and integrations will be updated as the project evolves.*

## 📚 Documentation

Comprehensive API documentation is available in the `docs/` directory:

- **[API Documentation](./docs/API_DOCUMENTATION.md)** - Complete API reference with endpoints and examples
- **[Developer Guide](./docs/DEVELOPER_GUIDE.md)** - Getting started guide and integration patterns  
- **[SDK Examples](./docs/SDK_EXAMPLES.md)** - Code examples for multiple programming languages
- **[Components Documentation](./docs/COMPONENTS.md)** - Internal architecture and component details
- **[OpenAPI Specification](./docs/openapi.yaml)** - Machine-readable API specification

### Quick API Reference

#### Base URLs
- Merchant API: `/merchant/api/v1`
- Frontend API: `/frontend/api/v1`  
- Integration API: `/integration/api/v1`

#### Key Endpoints
- `POST /auth/register/` - Register new user
- `POST /auth/login/` - User authentication
- `POST /merchant/` - Create merchant account
- `POST /payment/` - Create payment invoice
- `GET /payment/:id/` - Get payment details

#### Authentication
```http
Authorization: Bearer <access_token>
```

For detailed API documentation, examples, and integration guides, see the [docs directory](./docs/).

## 🚀 Quick Start

1. **Start the server**:
   ```bash
   go run cmd/app/main.go
   ```

2. **Register a merchant**:
   ```bash
   curl -X POST http://localhost:8080/merchant/api/v1/auth/register/ \
     -H "Content-Type: application/json" \
     -d '{"email":"merchant@example.com","password":"password123"}'
   ```

3. **Create a payment**:
   ```bash
   curl -X POST http://localhost:8080/merchant/api/v1/payment/ \
     -H "Authorization: Bearer <token>" \
     -d '{"requested_amount":100.50}'
   ```

For complete setup and usage instructions, see the [Developer Guide](./docs/DEVELOPER_GUIDE.md).