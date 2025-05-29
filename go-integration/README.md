# Freelance Payment Gateway - Golang Integration

This package provides a Golang integration for the EthJobEscrow smart contract, enabling blockchain-based payments for your freelancing platform.

## Features

- **Escrow Payments**: Secure payment holding using smart contracts
- **USD-ETH Conversion**: Automatic conversion using Chainlink price feeds
- **Job Management**: Post, complete, and cancel jobs on the blockchain
- **Real-time Price Data**: Get current ETH/USD prices
- **REST API**: HTTP endpoints for all blockchain operations

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Web Frontend  │◄──►│  Golang Backend  │◄──►│ Smart Contract  │
│                 │    │                  │    │   (Ethereum)    │
│   - Job UI      │    │  - REST API      │    │  - Escrow Logic │
│   - Payments    │    │  - Blockchain    │    │  - Price Feeds  │
│   - Status      │    │    Client        │    │  - Events       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Prerequisites

1. **Go 1.19+** installed
2. **Ethereum Node Access** (Infura, Alchemy, or local node)
3. **Smart Contract** deployed on testnet/mainnet
4. **Private Key** with ETH for gas fees
5. **Environment Variables** configured

## Quick Start

### 1. Clone and Setup

```bash
cd go-integration
cp env.example .env
# Edit .env with your configuration
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Deploy Smart Contract (if not deployed)

```bash
# Go back to the root directory
cd ../
forge script script/DeployPaymentGateway.s.sol:DeployEthJobEscrow --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY --broadcast
```

### 4. Configure Environment

Edit your `.env` file with the deployed contract address:

```env
ETHEREUM_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID
CONTRACT_ADDRESS=0xYourDeployedContractAddress
PRIVATE_KEY=your_private_key_without_0x
```

### 5. Run the Application

```bash
# Load environment variables
source .env

# Run the server
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Post a Job

**POST** `/post-job`

Creates a new job and locks ETH in escrow based on USD amount.

```bash
curl -X POST http://localhost:8080/post-job \
  -H "Content-Type: application/json" \
  -d '{
    "job_id": 1,
    "freelancer_address": "0x742C4356e2B18C51EB9D0CbaF6A1A6c0C8c7DBCE",
    "usd_amount": "100",
    "client_address": "0x8ba1f109551bD432803012645Hac136c4Ce7"
  }'
```

**Response:**
```json
{
  "tx_hash": "0x123...",
  "block_number": 12345,
  "gas_used": 150000,
  "success": true
}
```

### 2. Get Job Status

**GET** `/job-status?job_id=1`

Retrieves current job information from the blockchain.

```bash
curl http://localhost:8080/job-status?job_id=1
```

**Response:**
```json
{
  "job_id": 1,
  "client": "0x8ba1f109551bD432803012645Hac136c4Ce7",
  "freelancer": "0x742C4356e2B18C51EB9D0CbaF6A1A6c0C8c7DBCE",
  "usd_amount": "100",
  "eth_amount": "0.03125",
  "is_completed": false,
  "is_paid": false
}
```

### 3. Complete Job

**POST** `/complete-job?job_id=1`

Marks a job as completed and releases payment to freelancer.

```bash
curl -X POST http://localhost:8080/complete-job?job_id=1
```

### 4. Cancel Job

**POST** `/cancel-job?job_id=1`

Cancels a job and refunds the client.

```bash
curl -X POST http://localhost:8080/cancel-job?job_id=1
```

### 5. Get ETH Price

**GET** `/eth-price`

Gets current ETH/USD price from Chainlink.

```bash
curl http://localhost:8080/eth-price
```

**Response:**
```json
{
  "eth_usd_price": "320000000000"
}
```

## Integration in Your Web App

### 1. Import the Package

```go
import (
    "github.com/fahedafzaal/freelance-payment-gateway/internal/config"
    "github.com/fahedafzaal/freelance-payment-gateway/pkg/blockchain"
)
```

### 2. Initialize Client

```go
cfg := config.Load()
client, err := blockchain.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 3. Use Blockchain Operations

```go
// Post a job
result, err := client.PostJob(ctx, jobID, freelancerAddr, usdAmount, clientAddr)

// Get job details
details, err := client.GetJobDetails(ctx, jobID)

// Complete job
result, err := client.MarkJobCompleted(ctx, jobID)

// Cancel job
result, err := client.CancelJob(ctx, jobID)
```

## Configuration

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ETHEREUM_RPC_URL` | Ethereum node URL | `https://sepolia.infura.io/v3/...` |
| `NETWORK_ID` | Chain ID | `11155111` (Sepolia) |
| `CONTRACT_ADDRESS` | Deployed contract address | `0x123...` |
| `PRIVATE_KEY` | Private key for transactions | `abc123...` |
| `GAS_LIMIT` | Gas limit for transactions | `300000` |
| `GAS_PRICE` | Gas price in Gwei | `20` |
| `SERVER_PORT` | HTTP server port | `8080` |

### Network Support

- **Sepolia Testnet** (Chain ID: 11155111)
- **Ethereum Mainnet** (Chain ID: 1)

Price feeds are automatically configured based on the network.

## Security Considerations

1. **Private Key Management**: Never commit private keys to version control
2. **Environment Variables**: Use secure environment variable management
3. **Gas Limits**: Set appropriate gas limits to prevent failed transactions
4. **Input Validation**: Always validate user inputs before blockchain calls
5. **Error Handling**: Implement proper error handling for failed transactions

## Testing

### Unit Tests

```bash
go test ./...
```

### Integration Tests with Testnet

1. Deploy contract to Sepolia testnet
2. Configure environment variables
3. Run integration tests

```bash
go test -tags=integration ./...
```

## Troubleshooting

### Common Issues

1. **"insufficient funds"**: Ensure your wallet has enough ETH for gas
2. **"contract not found"**: Verify the contract address is correct
3. **"nonce too low"**: Wait for previous transactions to confirm
4. **"gas too low"**: Increase the gas limit in configuration

### Debugging

Enable debug logging:
```go
log.SetLevel(log.DebugLevel)
```

### Transaction Monitoring

Check transaction status on:
- Sepolia: https://sepolia.etherscan.io
- Mainnet: https://etherscan.io

## Production Deployment

### 1. Security Checklist

- [ ] Private keys stored securely (AWS KMS, HashiCorp Vault)
- [ ] Environment variables encrypted
- [ ] Rate limiting implemented
- [ ] Input validation on all endpoints
- [ ] Monitoring and alerting set up

### 2. Infrastructure

- [ ] Load balancer for multiple instances
- [ ] Database for off-chain data
- [ ] Redis for caching
- [ ] Backup Ethereum node endpoints

### 3. Monitoring

- [ ] Transaction success/failure rates
- [ ] Gas usage and costs
- [ ] API response times
- [ ] Blockchain node connectivity

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

MIT License - see LICENSE file for details. 