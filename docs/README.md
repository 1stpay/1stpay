# 1stPay Documentation

Welcome to the comprehensive documentation for the 1stPay cryptocurrency payments API.

## Documentation Structure

### 📚 Core Documentation

- **[API Documentation](./API_DOCUMENTATION.md)** - Complete API reference with endpoints, models, and examples
- **[OpenAPI Specification](./openapi.yaml)** - Machine-readable API specification for tools like Swagger UI
- **[Developer Guide](./DEVELOPER_GUIDE.md)** - Getting started guide and integration patterns
- **[Components Documentation](./COMPONENTS.md)** - Internal architecture and component details
- **[SDK Examples](./SDK_EXAMPLES.md)** - Code examples for multiple programming languages

### 🚀 Quick Start

1. **For API Integration**: Start with [API Documentation](./API_DOCUMENTATION.md)
2. **For Development**: Read the [Developer Guide](./DEVELOPER_GUIDE.md)
3. **For Code Examples**: Check [SDK Examples](./SDK_EXAMPLES.md)
4. **For Architecture**: Review [Components Documentation](./COMPONENTS.md)

### 🔧 API Reference

The API is organized into three main groups:

#### Merchant API (`/merchant/api/v1`)
For merchant operations and management:
- User authentication and registration
- Merchant profile management
- Payment creation and tracking
- Token configuration
- Account management

#### Frontend API (`/frontend/api/v1`)
For customer-facing payment interfaces:
- Public payment information
- Customer payment details
- Payment status tracking

#### Integration API (`/integration/api/v1`)
For external system integrations:
- Webhook endpoints
- Third-party system integration
- Automated payment processing

### 📖 Key Features Documented

- **Multi-Blockchain Support**: EVM, Tron, TON, Solana networks
- **Payment Processing**: Invoice creation, payment tracking, status management
- **Merchant Management**: Account setup, token configuration, commission management
- **Authentication**: JWT-based secure authentication
- **Error Handling**: Comprehensive error codes and handling patterns
- **Security**: Best practices for secure API integration

### 🛠 Tools and Utilities

#### Swagger UI
View the interactive API documentation:
```bash
# Install swagger-ui-serve
npm install -g swagger-ui-serve

# Serve the OpenAPI spec
swagger-ui-serve docs/openapi.yaml
```

#### Postman Collection
Import the OpenAPI specification into Postman for API testing:
1. Open Postman
2. Click "Import"
3. Select "Upload Files"
4. Choose `docs/openapi.yaml`

### 📋 API Endpoints Summary

#### Authentication
- `POST /auth/register/` - Register new user
- `POST /auth/login/` - User login

#### Merchant Management
- `POST /merchant/` - Create merchant
- `GET /merchant/me/` - Get merchant details
- `PUT /merchant/me/` - Update merchant
- `GET /merchant/me/tokens/` - List merchant tokens
- `POST /merchant/me/tokens/` - Add merchant token

#### Payment Processing
- `POST /payment/` - Create payment
- `GET /payment/:id/` - Get payment details

#### Reference Data
- `GET /blockchain/list/` - List blockchains
- `GET /token/list/` - List tokens

#### Utilities
- `GET /ping` - Health check
- `GET /user/me/` - User profile

### 🔐 Authentication

All protected endpoints require JWT authentication:

```http
Authorization: Bearer <access_token>
```

Get access tokens via:
- `POST /merchant/api/v1/auth/register/`
- `POST /merchant/api/v1/auth/login/`

### 💡 Integration Examples

#### Quick Payment Creation

```bash
# 1. Register/Login
curl -X POST http://localhost:8080/merchant/api/v1/auth/register/ \
  -H "Content-Type: application/json" \
  -d '{"email":"merchant@example.com","password":"password123"}'

# 2. Create merchant
curl -X POST http://localhost:8080/merchant/api/v1/merchant/ \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"My Store"}'

# 3. Create payment
curl -X POST http://localhost:8080/merchant/api/v1/payment/ \
  -H "Authorization: Bearer <token>" \
  -d '{"requested_amount":100.50,"email":"customer@example.com"}'
```

#### Customer Payment View

```bash
# Customer views payment details
curl -X GET http://localhost:8080/frontend/api/v1/payment/<payment-id>/
```

### 🐛 Troubleshooting

#### Common Issues

1. **401 Unauthorized**
   - Verify access token is valid
   - Check Authorization header format

2. **404 Merchant not found**
   - Ensure merchant profile is created
   - Verify user-merchant association

3. **400 Invalid request**
   - Validate request body format
   - Check required fields

#### Debug Mode

Enable detailed logging in your SDK implementation to troubleshoot issues.

### 📞 Support

- **GitHub Issues**: Report bugs and request features
- **API Status**: Monitor service uptime
- **Technical Support**: Contact for integration help

### 🔄 Updates

This documentation is automatically updated with each API version. Check the changelog for recent updates and breaking changes.

---

**Last Updated**: January 2024  
**API Version**: 1.0.0  
**Documentation Version**: 1.0.0