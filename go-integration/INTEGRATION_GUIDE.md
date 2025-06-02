# Complete Integration Guide for Freelance Payment Gateway

This guide provides step-by-step instructions to integrate the blockchain payment gateway with your existing freelance platform.

## 📋 Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Import Methods](#import-methods)
4. [Database Updates](#database-updates)
5. [Code Integration](#code-integration)
6. [Testing](#testing)
7. [Production Deployment](#production-deployment)

## 🎯 Overview

The payment gateway integrates with your application at three key points:
1. **Candidate accepts offer** → Fund escrow (`PostJob`)
2. **Poster approves work** → Release payment (`CompleteJob`)
3. **Job cancellation** → Refund client (`CancelJob`)

## 📋 Prerequisites

### 1. Database Schema
Ensure your database has these fields in the `applications` table:
```sql
-- Required fields (should already exist)
applications.payment_status VARCHAR(50)
applications.escrow_tx_hash_deposit VARCHAR(66)
applications.escrow_tx_hash_release VARCHAR(66)
applications.escrow_tx_hash_refund VARCHAR(66)
applications.agreed_usd_amount DECIMAL(10,2)

-- Required fields in users table
users.wallet_address VARCHAR(42)
```

### 2. Payment Gateway Running
Start the payment gateway service:
```bash
cd go-integration
make run
# Service should be running on http://localhost:8081
```

## 🔗 Import Methods

You have **two options** for integrating with the payment gateway:

### Option A: Using Go Client Package (Recommended)

#### Step 1: Copy the Client Package
Copy the client package to your main project:
```bash
# From your main project root
mkdir -p pkg/payment
cp -r /path/to/Freelance-Payment-Gateway/go-integration/pkg/blockchain/* pkg/payment/
```

#### Step 2: Update Import Paths
Edit `pkg/payment/service.go` and update the package declaration:
```go
package payment  // Change from 'package blockchain'
```

#### Step 3: Import in Your Code
```go
import "your-project/pkg/payment"
```

### Option B: Direct HTTP Calls
If you prefer not to import the package, use direct HTTP calls (examples provided below).

## 🔄 Database Updates

Add these helper methods to your `ApplicationService`:

```go
// InitiateEscrowDeposit updates payment status when escrow funding starts
func (as *ApplicationService) InitiateEscrowDeposit(ctx context.Context, applicationID int32, txHash string) error {
    return as.Queries.UpdateApplicationPaymentStatus(ctx, UpdateApplicationPaymentStatusParams{
        ID:                     applicationID,
        PaymentStatus:          sql.NullString{String: "deposit_initiated", Valid: true},
        EscrowTxHashDeposit:    sql.NullString{String: txHash, Valid: true},
    })
}

// InitiatePaymentRelease updates payment status when payment release starts
func (as *ApplicationService) InitiatePaymentRelease(ctx context.Context, applicationID int32, txHash string) error {
    return as.Queries.UpdateApplicationPaymentStatus(ctx, UpdateApplicationPaymentStatusParams{
        ID:                     applicationID,
        PaymentStatus:          sql.NullString{String: "release_initiated", Valid: true},
        EscrowTxHashRelease:    sql.NullString{String: txHash, Valid: true},
    })
}

// InitiateRefund updates payment status when refund starts
func (as *ApplicationService) InitiateRefund(ctx context.Context, applicationID int32, txHash string) error {
    return as.Queries.UpdateApplicationPaymentStatus(ctx, UpdateApplicationPaymentStatusParams{
        ID:                     applicationID,
        PaymentStatus:          sql.NullString{String: "refund_initiated", Valid: true},
        EscrowTxHashRefund:     sql.NullString{String: txHash, Valid: true},
    })
}
```

Add corresponding SQL queries to your `queries.sql`:
```sql
-- name: UpdateApplicationPaymentStatus :exec
UPDATE applications 
SET 
    payment_status = COALESCE(sqlc.narg('payment_status'), payment_status),
    escrow_tx_hash_deposit = COALESCE(sqlc.narg('escrow_tx_hash_deposit'), escrow_tx_hash_deposit),
    escrow_tx_hash_release = COALESCE(sqlc.narg('escrow_tx_hash_release'), escrow_tx_hash_release),
    escrow_tx_hash_refund = COALESCE(sqlc.narg('escrow_tx_hash_refund'), escrow_tx_hash_refund),
    updated_at = NOW()
WHERE id = $1;
```

Then regenerate your queries:
```bash
sqlc generate
```

## 💻 Code Integration

### 1. Initialize Payment Gateway Client

Add to your application service struct:
```go
type ApplicationService struct {
    Queries        *db.Queries
    PaymentGateway *payment.PaymentGatewayService  // Add this
    // ... other fields
}

// Update your constructor
func NewApplicationService(queries *db.Queries) *ApplicationService {
    return &ApplicationService{
        Queries:        queries,
        PaymentGateway: payment.NewPaymentGatewayService("http://localhost:8081"),
    }
}
```

### 2. Modify RespondToOffer Method

Update your `RespondToOffer` method to fund escrow when candidate accepts:

```go
func (as *ApplicationService) RespondToOffer(ctx context.Context, params RespondToOfferParams) error {
    // Start transaction
    tx, err := as.Queries.BeginTx(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // ... your existing code for status updates ...

    // NEW: Fund escrow when candidate accepts offer
    if params.Accept && newAppStatus == StatusHired {
        // Get required data
        app, err := as.Queries.GetApplicationByID(ctx, params.ApplicationID)
        if err != nil {
            return fmt.Errorf("failed to get application: %w", err)
        }

        job, err := as.Queries.GetJobByID(ctx, app.JobID)
        if err != nil {
            return fmt.Errorf("failed to get job: %w", err)
        }

        applicant, err := as.Queries.GetUserByID(ctx, app.UserID)
        if err != nil {
            return fmt.Errorf("failed to get applicant: %w", err)
        }

        poster, err := as.Queries.GetUserByID(ctx, job.UserID)
        if err != nil {
            return fmt.Errorf("failed to get poster: %w", err)
        }

        // Validate wallet addresses exist
        if !applicant.WalletAddress.Valid || !poster.WalletAddress.Valid {
            return fmt.Errorf("wallet addresses required for both users")
        }

        // Call payment gateway to fund escrow
        req := payment.PostJobRequest{
            JobID:             uint64(params.ApplicationID),
            FreelancerAddress: applicant.WalletAddress.String,
            USDAmount:         fmt.Sprintf("%.2f", float64(app.AgreedUsdAmount.Int32)),
            ClientAddress:     poster.WalletAddress.String,
        }

        result, err := as.PaymentGateway.PostJob(ctx, req)
        if err != nil {
            return fmt.Errorf("failed to fund escrow: %w", err)
        }

        // Update payment status in database
        err = as.InitiateEscrowDeposit(ctx, params.ApplicationID, result.TxHash)
        if err != nil {
            log.Printf("Warning: blockchain transaction successful but failed to update DB: %v", err)
            // Continue - don't fail the entire transaction for DB update issues
        }

        log.Printf("Escrow funded successfully! Application: %d, TxHash: %s", 
            params.ApplicationID, result.TxHash)
    }

    return tx.Commit(ctx)
}
```

### 3. Modify PosterReviewWork Method

Update your `PosterReviewWork` method to release payment when work is approved:

```go
func (as *ApplicationService) PosterReviewWork(ctx context.Context, params PosterReviewWorkParams) error {
    // Start transaction
    tx, err := as.Queries.BeginTx(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // ... your existing code for status updates ...

    // Get current application payment details
    currentAppPaymentDetails, err := as.Queries.GetApplicationPaymentDetails(ctx, params.ApplicationID)
    if err != nil {
        return fmt.Errorf("failed to get application payment details: %w", err)
    }

    // NEW: Release payment when work is approved and escrow is funded
    if params.NewStatus == StatusWorkApproved && 
       currentAppPaymentDetails.PaymentStatus.String == "deposited" {
        
        result, err := as.PaymentGateway.CompleteJob(ctx, uint64(params.ApplicationID))
        if err != nil {
            return fmt.Errorf("failed to release payment: %w", err)
        }

        // Update payment status in database
        err = as.InitiatePaymentRelease(ctx, params.ApplicationID, result.TxHash)
        if err != nil {
            log.Printf("Warning: blockchain transaction successful but failed to update DB: %v", err)
        }

        log.Printf("Payment released successfully! Application: %d, TxHash: %s", 
            params.ApplicationID, result.TxHash)
    }

    return tx.Commit(ctx)
}
```

### 4. Add Job Cancellation Support

Add a new method for handling job cancellations:

```go
func (as *ApplicationService) CancelJobWithRefund(ctx context.Context, applicationID int32, reason string) error {
    // Start transaction
    tx, err := as.Queries.BeginTx(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Get current application payment details
    currentAppPaymentDetails, err := as.Queries.GetApplicationPaymentDetails(ctx, applicationID)
    if err != nil {
        return fmt.Errorf("failed to get application payment details: %w", err)
    }

    // Only refund if escrow is funded but not yet released
    if currentAppPaymentDetails.PaymentStatus.String == "deposited" {
        result, err := as.PaymentGateway.CancelJob(ctx, uint64(applicationID))
        if err != nil {
            return fmt.Errorf("failed to cancel job and refund: %w", err)
        }

        // Update payment status in database
        err = as.InitiateRefund(ctx, applicationID, result.TxHash)
        if err != nil {
            log.Printf("Warning: blockchain transaction successful but failed to update DB: %v", err)
        }

        log.Printf("Job cancelled and refund initiated! Application: %d, TxHash: %s", 
            applicationID, result.TxHash)
    }

    // Update application status to cancelled
    err = as.Queries.UpdateApplicationStatus(ctx, UpdateApplicationStatusParams{
        ID:     applicationID,
        Status: StatusCancelled,
    })
    if err != nil {
        return fmt.Errorf("failed to update application status: %w", err)
    }

    return tx.Commit(ctx)
}
```

### 5. Add Transaction Status Checking

Add a background service to check transaction status:

```go
// CheckTransactionStatus periodically checks blockchain transaction status
func (as *ApplicationService) CheckTransactionStatus(ctx context.Context, applicationID int32) error {
    status, err := as.PaymentGateway.GetJobStatus(ctx, uint64(applicationID))
    if err != nil {
        return fmt.Errorf("failed to get job status: %w", err)
    }

    // Update database based on blockchain status
    updates := UpdateApplicationPaymentStatusParams{
        ID: applicationID,
    }

    // Map blockchain status to database status
    switch status.PaymentStatus {
    case "deposited":
        updates.PaymentStatus = sql.NullString{String: "deposited", Valid: true}
    case "released":
        updates.PaymentStatus = sql.NullString{String: "released", Valid: true}
    case "refunded":
        updates.PaymentStatus = sql.NullString{String: "refunded", Valid: true}
    }

    if updates.PaymentStatus.Valid {
        err = as.Queries.UpdateApplicationPaymentStatus(ctx, updates)
        if err != nil {
            return fmt.Errorf("failed to update payment status: %w", err)
        }
    }

    return nil
}

// StartTransactionMonitor runs a background service to monitor transaction status
func (as *ApplicationService) StartTransactionMonitor(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Get all applications with pending transactions
            apps, err := as.Queries.GetApplicationsWithPendingTransactions(ctx)
            if err != nil {
                log.Printf("Failed to get pending applications: %v", err)
                continue
            }

            for _, app := range apps {
                err := as.CheckTransactionStatus(ctx, app.ID)
                if err != nil {
                    log.Printf("Failed to check transaction status for app %d: %v", app.ID, err)
                }
            }
        }
    }
}
```

## 🔄 Alternative: Direct HTTP Integration

If you prefer not to import the package, use direct HTTP calls:

```go
import (
    "bytes"
    "encoding/json"
    "net/http"
)

// Fund escrow with direct HTTP call
func (as *ApplicationService) fundEscrowHTTP(ctx context.Context, applicationID int32, freelancerAddr, clientAddr string, amount string) (string, error) {
    reqBody := map[string]interface{}{
        "job_id":             applicationID,
        "freelancer_address": freelancerAddr,
        "usd_amount":         amount,
        "client_address":     clientAddr,
    }

    jsonBody, err := json.Marshal(reqBody)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8081/post-job", bytes.NewReader(jsonBody))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
    }

    var result struct {
        TxHash string `json:"tx_hash"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result.TxHash, err
}

// Release payment with direct HTTP call
func (as *ApplicationService) releasePaymentHTTP(ctx context.Context, applicationID int32) (string, error) {
    url := fmt.Sprintf("http://localhost:8081/complete-job?job_id=%d", applicationID)
    
    req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
    if err != nil {
        return "", err
    }

    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
    }

    var result struct {
        TxHash string `json:"tx_hash"`
    }
    err = json.NewDecoder(resp.Body).Decode(&result)
    return result.TxHash, err
}
```

## 🧪 Testing

### 1. Unit Tests
Create tests for your modified methods:

```go
func TestRespondToOfferWithEscrow(t *testing.T) {
    // Mock payment gateway
    mockGateway := &MockPaymentGateway{}
    service := &ApplicationService{
        Queries:        mockQueries,
        PaymentGateway: mockGateway,
    }

    // Test accepting offer triggers escrow funding
    err := service.RespondToOffer(ctx, RespondToOfferParams{
        ApplicationID: 123,
        Accept:       true,
    })

    assert.NoError(t, err)
    assert.True(t, mockGateway.PostJobCalled)
}
```

### 2. Integration Tests
Test the complete flow:

```bash
# Start payment gateway
cd go-integration
make run

# In another terminal, test your main application
cd your-main-app
go test ./... -tags=integration
```

### 3. Manual Testing
1. Create a job posting
2. Apply as a freelancer
3. Accept the application (should fund escrow)
4. Submit work
5. Approve work (should release payment)

## 🚀 Production Deployment

### 1. Environment Configuration
Set production environment variables:
```env
# In your main application
PAYMENT_GATEWAY_URL=https://payment-gateway.yourdomain.com

# In payment gateway
CONTRACT_ADDRESS=0xYourProductionContractAddress
ETHEREUM_RPC_URL=https://mainnet.infura.io/v3/YOUR_PROJECT_ID
PRIVATE_KEY=your_production_private_key
```

### 2. Monitoring
Add monitoring for:
- Payment gateway health
- Transaction success rates
- Gas costs
- Failed transactions

```go
// Add health check endpoint
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
    // Check payment gateway health
    resp, err := http.Get("http://localhost:8081/health")
    if err != nil || resp.StatusCode != 200 {
        http.Error(w, "Payment gateway unavailable", http.StatusServiceUnavailable)
        return
    }
    
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
```

### 3. Error Handling
Implement robust error handling:
- Retry logic for failed transactions
- Fallback mechanisms
- User notifications for transaction status

### 4. Security
- Use environment variables for sensitive data
- Implement proper authentication
- Monitor for suspicious transactions
- Set up alerts for failed payments

## 📚 Additional Resources

- **Smart Contract**: See `/src/EthJobEscrow.sol` for contract details
- **API Documentation**: All endpoints documented in main README
- **Error Codes**: Check payment gateway logs for detailed error messages
- **Gas Optimization**: Monitor gas usage and adjust GAS_PRICE as needed

## 🆘 Troubleshooting

**Common Issues:**

1. **"Wallet address required"**
   - Ensure all users have `wallet_address` set
   - Validate addresses are valid Ethereum addresses

2. **"Payment gateway unavailable"**
   - Check if payment gateway is running on port 8081
   - Verify network connectivity

3. **Transaction failures**
   - Check Ethereum network status
   - Ensure sufficient ETH for gas fees
   - Verify contract address is correct

4. **Database inconsistencies**
   - Use transactions for atomic operations
   - Implement retry logic for DB updates
   - Monitor for orphaned blockchain transactions

That's it! Your freelance platform now has secure, trustless blockchain payments integrated into your existing workflow. 🎉 