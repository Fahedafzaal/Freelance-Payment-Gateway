# Migration Guide: From Manual Payments to Payment Gateway

This guide shows how to migrate your existing manual payment system to use the automated payment gateway.

## 🔄 Changes Overview

**Before (Manual):** Users trigger blockchain transactions through UI → Your app records transaction hashes
**After (Automated):** Your app triggers blockchain transactions → Payment gateway handles everything

## 📝 Required Changes

### 1. Update PaymentController

Your existing `PaymentController` needs to be simplified since the payment gateway will handle blockchain interactions:

```go
// controllers/payment-controller.go - UPDATED VERSION
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	appContext "fyp-module-one.com/one/context"
	"fyp-module-one.com/one/models"
	"fyp-module-one.com/one/pkg/payment" // NEW: Import payment gateway client
	"fyp-module-one.com/one/views/auth_views"
	"fyp-module-one.com/one/views/types"
	"github.com/go-chi/chi/v5"
)

type PaymentController struct {
	AppService     *models.ApplicationService
	PaymentGateway *payment.PaymentGatewayService // NEW: Payment gateway client
	EthConfig      types.EthConfig
}

func NewPaymentController(as *models.ApplicationService, ethCfg types.EthConfig) *PaymentController {
	return &PaymentController{
		AppService:     as,
		PaymentGateway: payment.NewPaymentGatewayService("http://localhost:8081"), // NEW
		EthConfig:      ethCfg,
	}
}

// GET /app/applications/{applicationID}/payment/status
// NEW: Get payment status from payment gateway
func (pc *PaymentController) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	currentUser := appContext.User(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appIDStr := chi.URLParam(r, "applicationID")
	appIDInt, err := strconv.Atoi(appIDStr)
	if err != nil {
		http.Error(w, "Invalid Application ID", http.StatusBadRequest)
		return
	}

	// Check authorization
	appPaymentCtx, err := pc.AppService.GetApplicationPaymentContext(r.Context(), int32(appIDInt))
	if err != nil {
		http.Error(w, "Application context not found", http.StatusNotFound)
		return
	}
	if currentUser.ID != appPaymentCtx.PosterUserID && currentUser.ID != appPaymentCtx.ApplicantUserID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Get status from payment gateway
	status, err := pc.PaymentGateway.GetJobStatus(r.Context(), uint64(appIDInt))
	if err != nil {
		http.Error(w, "Failed to get payment status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// POST /app/applications/{applicationID}/payment/cancel
// NEW: Cancel job and initiate refund through payment gateway
func (pc *PaymentController) HandleCancelWithRefund(w http.ResponseWriter, r *http.Request) {
	currentUser := appContext.User(r.Context())
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appIDStr := chi.URLParam(r, "applicationID")
	appIDInt, err := strconv.Atoi(appIDStr)
	if err != nil {
		http.Error(w, "Invalid Application ID", http.StatusBadRequest)
		return
	}

	// Check authorization - only poster can cancel
	appPaymentCtx, err := pc.AppService.GetApplicationPaymentContext(r.Context(), int32(appIDInt))
	if err != nil {
		http.Error(w, "Application context not found", http.StatusNotFound)
		return
	}
	if appPaymentCtx.PosterUserID != currentUser.ID {
		http.Error(w, "Unauthorized to cancel this application", http.StatusForbidden)
		return
	}

	// Call payment gateway to cancel and refund
	result, err := pc.PaymentGateway.CancelJob(r.Context(), uint64(appIDInt))
	if err != nil {
		http.Error(w, "Failed to cancel job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update database
	err = pc.AppService.InitiateEscrowRefund(r.Context(), int32(appIDInt), result.TxHash)
	if err != nil {
		// Log error but don't fail since blockchain transaction succeeded
		fmt.Printf("Warning: blockchain refund successful but failed to update DB: %v\n", err)
	}

	// Return updated view
	updatedApp, err := pc.AppService.GetApplicationDetailByID(r.Context(), int32(appIDInt))
	if err != nil {
		fmt.Printf("Error fetching updated app details: %v\n", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	viewData := auth_views.ApplicationDetailData{
		User:        currentUser,
		Application: updatedApp,
		IsJobOwner:  appPaymentCtx.PosterUserID == currentUser.ID,
		EthConfig:   pc.EthConfig,
	}
	component := auth_views.ApplicationPaymentStatusSection(viewData)
	component.Render(r.Context(), w)
}

// REMOVE: These methods are no longer needed since payment gateway handles transactions automatically
// - HandleInitiateDeposit (replaced by automatic funding in RespondToOffer)
// - HandleInitiateRelease (replaced by automatic release in PosterReviewWork)  
// - HandleInitiateRefund (replaced by HandleCancelWithRefund above)
```

### 2. Update ApplicationService

Your `ApplicationService` needs the payment gateway integration. Add these changes:

```go
// models/application_service.go - ADD THESE CHANGES

import (
	"fyp-module-one.com/one/pkg/payment" // NEW: Import payment gateway
)

type ApplicationService struct {
	// ... existing fields ...
	PaymentGateway *payment.PaymentGatewayService // NEW: Add payment gateway
}

// Update constructor
func NewApplicationService(pool *pgxpool.Pool, queries *db.Queries, fileStore FileStore, jobService *JobService, userService *UserService) *ApplicationService {
	return &ApplicationService{
		// ... existing fields ...
		PaymentGateway: payment.NewPaymentGatewayService("http://localhost:8081"), // NEW
	}
}

// MODIFY: RespondToOffer - Add automatic escrow funding
func (as *ApplicationService) RespondToOffer(ctx context.Context, params RespondToOfferParams) error {
	// ... existing transaction and validation code ...

	// NEW: Fund escrow when candidate accepts offer
	if params.Accept && newAppStatus == StatusHired {
		// Get required data for payment gateway
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

		// Validate wallet addresses
		if !applicant.WalletAddress.Valid || !poster.WalletAddress.Valid {
			return fmt.Errorf("both users must have wallet addresses set")
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

		// Update payment status
		err = as.InitiateEscrowDeposit(ctx, params.ApplicationID, result.TxHash)
		if err != nil {
			// Log warning but don't fail - blockchain transaction succeeded
			fmt.Printf("Warning: escrow funded but failed to update DB: %v\n", err)
		}

		fmt.Printf("Escrow funded automatically! Application: %d, TxHash: %s\n", 
			params.ApplicationID, result.TxHash)
	}

	return tx.Commit(ctx)
}

// MODIFY: PosterReviewWork - Add automatic payment release
func (as *ApplicationService) PosterReviewWork(ctx context.Context, params PosterReviewWorkParams) error {
	// ... existing transaction and validation code ...

	// NEW: Release payment when work is approved
	if params.NewStatus == StatusWorkApproved {
		// Check if escrow is funded
		appPaymentDetails, err := as.Queries.GetApplicationPaymentDetails(ctx, params.ApplicationID)
		if err != nil {
			return fmt.Errorf("failed to get payment details: %w", err)
		}

		if appPaymentDetails.PaymentStatus.String == "deposited" {
			// Call payment gateway to release payment
			result, err := as.PaymentGateway.CompleteJob(ctx, uint64(params.ApplicationID))
			if err != nil {
				return fmt.Errorf("failed to release payment: %w", err)
			}

			// Update payment status
			err = as.InitiatePaymentRelease(ctx, params.ApplicationID, result.TxHash)
			if err != nil {
				fmt.Printf("Warning: payment released but failed to update DB: %v\n", err)
			}

			fmt.Printf("Payment released automatically! Application: %d, TxHash: %s\n", 
				params.ApplicationID, result.TxHash)
		}
	}

	return tx.Commit(ctx)
}
```

### 3. Update main.go

Remove EthService dependency and update PaymentController initialization:

```go
// main.go - CHANGES

func main() {
	// ... existing code ...

	// REMOVE: EthService initialization - no longer needed
	// ethService, err := models.NewEthService(...)

	// ... existing service initialization ...

	// UPDATE: PaymentController - remove ethService parameter
	paymentC := controllers.NewPaymentController(applicationService, ethConfigForTemplates)

	// ... rest of main.go stays the same ...
}
```

### 4. Update Routes

Modify your payment routes to match the new automated approach:

```go
// main.go - UPDATE PAYMENT ROUTES

r.Route("/applications/{applicationID}/payment", func(r chi.Router) {
	r.Get("/status", paymentC.GetPaymentStatus)           // NEW: Get status from gateway
	r.Post("/cancel", paymentC.HandleCancelWithRefund)    // NEW: Cancel with refund
	r.Get("/context", paymentC.GetPaymentContext)         // KEEP: Still useful for UI
	// REMOVE: Manual transaction routes no longer needed
	// - initiate-deposit (now automatic)
	// - initiate-release (now automatic) 
	// - initiate-refund (replaced by cancel)
})
```

### 5. Frontend Changes

Update your frontend JavaScript/HTMX calls:

**Before (Manual):**
```javascript
// User clicks "Fund Escrow" button → calls /initiate-deposit with txHash
fetch(`/app/jobs/${jobId}/applications/${appId}/payment/initiate-deposit`, {
    method: 'POST',
    body: JSON.stringify({ txHash: userProvidedTxHash })
})
```

**After (Automated):**
```javascript
// No manual transaction buttons needed!
// Escrow funding happens automatically when candidate accepts offer
// Payment release happens automatically when poster approves work

// Only need status checking:
fetch(`/app/applications/${appId}/payment/status`)
  .then(response => response.json())
  .then(status => {
    // Update UI based on payment status
    updatePaymentStatusUI(status);
  });
```

### 6. Database Migration

Your existing database schema should work, but you might want to add a migration to ensure compatibility:

```sql
-- migrations/add_payment_gateway_support.sql

-- Ensure payment_status field supports new statuses
ALTER TABLE applications 
ALTER COLUMN payment_status TYPE VARCHAR(50);

-- Add index for faster lookups of pending transactions
CREATE INDEX IF NOT EXISTS idx_applications_payment_status 
ON applications(payment_status) 
WHERE payment_status IN ('deposit_initiated', 'release_initiated', 'refund_initiated');

-- Update any existing 'pending' statuses to the new format
UPDATE applications 
SET payment_status = 'deposit_initiated' 
WHERE payment_status = 'pending_deposit';

UPDATE applications 
SET payment_status = 'release_initiated' 
WHERE payment_status = 'pending_release';
```

## 🗂️ Files to Remove/Archive

Since the payment gateway handles blockchain interactions, you can remove:

1. **`models/ethservice.go`** - No longer needed
2. **Manual transaction UI components** - Replace with status displays
3. **Frontend wallet connection code** - Payment gateway handles this

## 📋 Testing Migration

1. **Start payment gateway:** `cd go-integration && make run`
2. **Test acceptance flow:** Candidate accepts offer → Check escrow funded automatically
3. **Test completion flow:** Poster approves work → Check payment released automatically  
4. **Test cancellation flow:** Poster cancels → Check refund initiated

## 🎯 Benefits After Migration

- **Automatic transactions** - No manual user interaction needed
- **Consistent gas prices** - Payment gateway optimizes gas
- **Better error handling** - Centralized blockchain error management
- **Simplified UI** - Remove complex wallet connection flows
- **Faster transactions** - Payment gateway can batch and optimize

Your users will have a much smoother experience with automatic blockchain payments! 🚀 