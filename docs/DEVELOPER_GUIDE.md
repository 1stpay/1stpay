# 1stPay Developer Guide

## Getting Started

This guide provides comprehensive information for developers integrating with the 1stPay cryptocurrency payments API.

## Quick Start

### 1. Authentication Setup

Before using any protected endpoints, you need to authenticate:

```bash
# Register a new account
curl -X POST http://localhost:8080/merchant/api/v1/auth/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-secure-password"
  }'
```

Save the returned `access_token` for subsequent requests.

### 2. Create Merchant Profile

```bash
curl -X POST http://localhost:8080/merchant/api/v1/merchant/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Your Store Name"
  }'
```

### 3. Configure Accepted Tokens

First, get available tokens:

```bash
curl -X GET http://localhost:8080/merchant/api/v1/token/list/
```

Then configure which tokens your merchant accepts:

```bash
curl -X POST http://localhost:8080/merchant/api/v1/merchant/me/tokens/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "token_id": "TOKEN_UUID_FROM_LIST",
    "active": true
  }'
```

### 4. Create a Payment

```bash
curl -X POST http://localhost:8080/merchant/api/v1/payment/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "requested_amount": 100.50,
    "email": "customer@example.com"
  }'
```

### 5. Direct Customer to Payment Page

Use the returned payment ID to direct customers to:
```
http://localhost:8080/frontend/api/v1/payment/{PAYMENT_ID}/
```

## API Integration Patterns

### 1. E-commerce Integration

For e-commerce platforms:

```javascript
// Example Node.js integration
const axios = require('axios');

class FirstPayAPI {
  constructor(baseURL, accessToken) {
    this.client = axios.create({
      baseURL: baseURL,
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    });
  }

  async createPayment(amount, customerEmail) {
    try {
      const response = await this.client.post('/merchant/api/v1/payment/', {
        requested_amount: amount,
        email: customerEmail
      });
      return response.data;
    } catch (error) {
      throw new Error(`Payment creation failed: ${error.response.data.error}`);
    }
  }

  async getPaymentDetails(paymentId) {
    try {
      const response = await this.client.get(`/frontend/api/v1/payment/${paymentId}/`);
      return response.data;
    } catch (error) {
      throw new Error(`Payment retrieval failed: ${error.response.data.error}`);
    }
  }
}

// Usage
const api = new FirstPayAPI('http://localhost:8080', 'YOUR_ACCESS_TOKEN');
const payment = await api.createPayment(100.50, 'customer@example.com');
console.log('Payment created:', payment.id);
```

### 2. Webhook Integration

Configure webhooks to receive real-time payment status updates:

```javascript
// Express.js webhook handler example
app.post('/webhooks/1stpay', (req, res) => {
  const { payment_id, status, amount, transaction_hash } = req.body;
  
  switch(status) {
    case 'completed':
      // Payment completed successfully
      console.log(`Payment ${payment_id} completed: ${amount}`);
      // Update your database, fulfill order, etc.
      break;
    case 'failed':
      // Payment failed
      console.log(`Payment ${payment_id} failed`);
      // Handle failed payment
      break;
    case 'pending':
      // Payment is still pending
      console.log(`Payment ${payment_id} is pending`);
      break;
  }
  
  res.status(200).send('OK');
});
```

### 3. Mobile App Integration

For mobile applications:

```swift
// Swift iOS example
import Foundation

class FirstPayAPI {
    private let baseURL: String
    private let accessToken: String
    
    init(baseURL: String, accessToken: String) {
        self.baseURL = baseURL
        self.accessToken = accessToken
    }
    
    func createPayment(amount: Double, email: String, completion: @escaping (Result<PaymentResponse, Error>) -> Void) {
        guard let url = URL(string: "\(baseURL)/merchant/api/v1/payment/") else { return }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(accessToken)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        
        let body = [
            "requested_amount": amount,
            "email": email
        ]
        
        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: body)
        } catch {
            completion(.failure(error))
            return
        }
        
        URLSession.shared.dataTask(with: request) { data, response, error in
            // Handle response
        }.resume()
    }
}
```

## Supported Blockchains and Tokens

### Blockchain Networks

1. **EVM Compatible**
   - Ethereum Mainnet
   - Binance Smart Chain (BSC)
   - Arbitrum
   - Polygon
   - And other EVM-compatible networks

2. **Tron Network**
   - TRX and TRC-20 tokens

3. **TON Network**
   - TON and jetton tokens

4. **Solana Network**
   - SOL and SPL tokens

### Token Configuration

Each token has the following properties:
- **Name**: Human-readable token name
- **Symbol**: Token ticker symbol
- **Contract Address**: Smart contract address (for non-native tokens)
- **Decimals**: Number of decimal places
- **Is Native**: Whether the token is native to the blockchain
- **Is Active**: Whether the token is currently supported

## Payment Flow Architecture

### 1. Payment Creation Flow

```mermaid
sequenceDiagram
    participant M as Merchant
    participant API as 1stPay API
    participant DB as Database
    participant BC as Blockchain

    M->>API: POST /payment/ (amount, email)
    API->>DB: Create payment record
    API->>BC: Generate payment addresses
    API->>DB: Store payment addresses
    API->>M: Return payment ID
```

### 2. Customer Payment Flow

```mermaid
sequenceDiagram
    participant C as Customer
    participant Frontend as Frontend API
    participant BC as Blockchain
    participant API as 1stPay API

    C->>Frontend: GET /payment/{id}/
    Frontend->>API: Fetch payment details
    API->>C: Return payment addresses
    C->>BC: Send crypto to address
    BC->>API: Payment detected (webhook/polling)
    API->>API: Update payment status
```

## Error Handling Best Practices

### 1. Retry Logic

Implement exponential backoff for transient errors:

```javascript
async function apiCallWithRetry(apiCall, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      return await apiCall();
    } catch (error) {
      if (error.response?.status >= 500 && i < maxRetries - 1) {
        await new Promise(resolve => setTimeout(resolve, Math.pow(2, i) * 1000));
        continue;
      }
      throw error;
    }
  }
}
```

### 2. Error Response Handling

```javascript
function handleAPIError(error) {
  if (error.response) {
    const { status, data } = error.response;
    
    switch (status) {
      case 400:
        console.error('Bad Request:', data.error);
        // Handle validation errors
        break;
      case 401:
        console.error('Unauthorized:', data.error);
        // Redirect to login or refresh token
        break;
      case 403:
        console.error('Forbidden:', data.error);
        // Handle permission errors
        break;
      case 404:
        console.error('Not Found:', data.error);
        // Handle missing resources
        break;
      case 500:
        console.error('Server Error:', data.error);
        // Handle server errors, maybe retry
        break;
      default:
        console.error('Unexpected Error:', data.error);
    }
  }
}
```

## Security Considerations

### 1. Token Security

- Store JWT tokens securely (encrypted storage)
- Implement token refresh mechanisms
- Never expose tokens in client-side code
- Use environment variables for API credentials

### 2. API Security

- Always use HTTPS in production
- Validate all input data
- Implement rate limiting on your side
- Log API calls for monitoring

### 3. Payment Security

- Verify payment amounts and addresses
- Implement transaction monitoring
- Use webhooks for real-time updates
- Validate blockchain transactions independently

## Testing

### 1. Unit Testing

Test your API integration:

```javascript
// Jest example
describe('1stPay API Integration', () => {
  test('should create payment successfully', async () => {
    const api = new FirstPayAPI(testBaseURL, testToken);
    const payment = await api.createPayment(100.50, 'test@example.com');
    
    expect(payment).toHaveProperty('id');
    expect(payment).toHaveProperty('created_at');
    expect(payment).toHaveProperty('updated_at');
  });
  
  test('should handle invalid amount', async () => {
    const api = new FirstPayAPI(testBaseURL, testToken);
    
    await expect(api.createPayment(-100, 'test@example.com'))
      .rejects.toThrow('Payment creation failed');
  });
});
```

### 2. Integration Testing

Use the ping endpoints to verify connectivity:

```bash
curl -X GET http://localhost:8080/merchant/api/v1/ping
```

Expected response:
```json
{
  "message": "pong"
}
```

## Rate Limiting

While not explicitly implemented in the current API, follow these best practices:

- Implement client-side rate limiting
- Use exponential backoff for retries
- Monitor API usage patterns
- Cache reference data (tokens, blockchains) when possible

## Monitoring and Logging

### 1. API Call Logging

Log all API interactions:

```javascript
function logAPICall(method, url, requestData, responseData, duration) {
  console.log({
    timestamp: new Date().toISOString(),
    method,
    url,
    requestData,
    responseData,
    duration,
    status: responseData.status
  });
}
```

### 2. Payment Monitoring

Track payment statuses:

```javascript
function monitorPayment(paymentId) {
  const checkStatus = async () => {
    try {
      const payment = await api.getPaymentDetails(paymentId);
      
      switch (payment.status) {
        case 'completed':
          console.log(`Payment ${paymentId} completed`);
          return true; // Stop monitoring
        case 'failed':
          console.log(`Payment ${paymentId} failed`);
          return true; // Stop monitoring
        case 'pending':
          console.log(`Payment ${paymentId} still pending`);
          return false; // Continue monitoring
      }
    } catch (error) {
      console.error('Error checking payment status:', error);
      return false;
    }
  };
  
  const interval = setInterval(async () => {
    const shouldStop = await checkStatus();
    if (shouldStop) {
      clearInterval(interval);
    }
  }, 30000); // Check every 30 seconds
}
```

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check if access token is valid
   - Verify token is included in Authorization header
   - Ensure token hasn't expired

2. **404 Merchant not found**
   - Verify merchant profile has been created
   - Check if user is associated with a merchant

3. **400 Invalid request**
   - Validate request body format
   - Check required fields are present
   - Verify data types match expected format

4. **500 Internal Server Error**
   - Check server logs
   - Verify database connectivity
   - Contact API support if persistent

### Debug Mode

Enable debug logging to troubleshoot issues:

```javascript
const api = new FirstPayAPI(baseURL, accessToken, { debug: true });
```

This will log all requests and responses for debugging purposes.

## Production Deployment

### Environment Configuration

Set up environment variables:

```bash
# API Configuration
FIRSTPAY_API_URL=https://api.1stpay.com
FIRSTPAY_ACCESS_TOKEN=your_production_token

# Security
HTTPS_ENABLED=true
JWT_SECRET=your_jwt_secret

# Database
DATABASE_URL=postgresql://user:pass@localhost/1stpay

# Blockchain Configuration
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID
BSC_RPC_URL=https://bsc-dataseed.binance.org/
```

### Production Checklist

- [ ] Use HTTPS for all API calls
- [ ] Implement proper error handling
- [ ] Set up monitoring and alerting
- [ ] Configure webhook endpoints
- [ ] Test payment flows end-to-end
- [ ] Implement transaction verification
- [ ] Set up backup and recovery procedures
- [ ] Configure rate limiting
- [ ] Implement logging and auditing

## Support and Resources

### Documentation
- [API Reference](./API_DOCUMENTATION.md)
- [OpenAPI Specification](./openapi.yaml)
- [Component Documentation](./COMPONENTS.md)

### Community
- GitHub Issues: Report bugs and feature requests
- Developer Forum: Community discussions
- API Status: Monitor API uptime and incidents

### Contact
- Technical Support: support@1stpay.com
- Sales: sales@1stpay.com
- Security Issues: security@1stpay.com