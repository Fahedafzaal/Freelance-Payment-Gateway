# Freelance Payment Gateway Integration

A secure, trustless payment gateway for freelance platforms using blockchain escrow.

## 🚀 Quick Start

1. **Configure Environment**
   ```bash
   cp env.example .env
   # Edit .env with your credentials
   ```

2. **Deploy Smart Contract**
   ```bash
   make deploy-contract
   ```

3. **Start Payment Gateway**
   ```bash
   make run
   ```

## 📋 Integration Guide

### 1. Database Schema

The payment gateway works with your existing PostgreSQL schema:

```sql
applications (
    id SERIAL PRIMARY KEY,
    payment_status VARCHAR(50),
    escrow_tx_hash_deposit VARCHAR(66),
    escrow_tx_hash_release VARCHAR(66),
    escrow_tx_hash_refund VARCHAR(66),
    agreed_usd_amount DECIMAL(10,2)
)

users (
    id SERIAL PRIMARY KEY,
    wallet_address VARCHAR(42)
)

jobs (
    id SERIAL PRIMARY KEY,
    -- other fields
)
```

### 2. API Endpoints

#### POST /post-job
Called when candidate accepts offer → funds escrow
```json
{
    "job_id": "123",              // applications.id
    "freelancer_address": "0x...", // applicant wallet
    "usd_amount": "100.00",       // agreed_usd_amount
    "client_address": "0x..."     // poster wallet
}
```

#### POST /complete-job
Called when poster approves work → releases payment
```json
{
    "job_id": "123"  // applications.id
}
```

#### POST /cancel-job
Called for refunds
```json
{
    "job_id": "123"  // applications.id
}
```

#### GET /job-status
Returns payment status
```json
{
    "job_id": "123"  // applications.id
}
```

### 3. Integration Example

```go
// In your RespondToOffer handler
func (h *Handler) RespondToOffer(w http.ResponseWriter, r *http.Request) {
    // ... your existing code ...

    // Call payment gateway to fund escrow
    resp, err := h.paymentGateway.PostJob(ctx, &blockchain.PostJobRequest{
        JobID:             application.ID,
        FreelancerAddress: freelancer.WalletAddress,
        USDAmount:         application.AgreedUSDAmount,
        ClientAddress:     client.WalletAddress,
    })
    if err != nil {
        // Handle error
        return
    }

    // Update application status
    application.PaymentStatus = "deposit_initiated"
    application.EscrowTxHashDeposit = resp.TxHash
    // Save to database
}
```

## 🔧 Configuration

### Environment Variables
```env
# Server
PORT=8081
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=freelance_platform

# Blockchain
PRIVATE_KEY=your_private_key
RPC_URL=your_rpc_url
CONTRACT_ADDRESS=your_contract_address
```

### Docker Support
```bash
# Build and run with Docker
docker-compose up --build
```

## 🛠 Development

### Prerequisites
- Go 1.21+
- PostgreSQL 14+
- Foundry (for smart contracts)

### Commands
```bash
# Install dependencies
go mod tidy

# Run tests
make test

# Build
make build

# Run
make run
```

## 📝 Notes

- Uses `applications.id` as the escrow `jobId` on blockchain
- All payment tracking happens at the application level
- Multiple freelancers can work on different applications for the same job
- Failed blockchain calls don't corrupt your database
- All transaction hashes are recorded for transparency

## 🔍 Troubleshooting

1. **Database Connection Issues**
   - Check PostgreSQL is running
   - Verify credentials in .env
   - Ensure database exists

2. **Blockchain Transaction Failures**
   - Check RPC URL is correct
   - Verify private key has enough gas
   - Check contract address is correct

3. **API Errors**
   - Check application exists
   - Verify wallet addresses are valid
   - Ensure amounts are properly formatted

## 📚 Additional Resources

- [Smart Contract Documentation](./contracts/README.md)
- [API Documentation](./pkg/blockchain/README.md)
- [Database Schema](./pkg/database/README.md) 