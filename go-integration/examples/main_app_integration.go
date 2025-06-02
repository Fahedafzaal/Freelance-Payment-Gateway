package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Fahedafzaal/Freelance-Payment-Gateway/go-integration/pkg/payment"
)

// Example integration code for your main application
// This shows how to use the payment gateway service in your existing handlers

func main() {
	// Initialize the payment gateway service client
	paymentGateway := payment.NewPaymentGatewayService("http://localhost:8081")

	// Example: Call from your RespondToOffer method
	// This is what you would add to your ApplicationService.RespondToOffer method
	exampleRespondToOfferIntegration(paymentGateway)

	// Example: Call from your PosterReviewWork method
	// This is what you would add to your ApplicationService.PosterReviewWork method
	examplePosterReviewWorkIntegration(paymentGateway)
}

// Example integration for RespondToOffer (when candidate accepts offer)
func exampleRespondToOfferIntegration(paymentGateway *payment.PaymentGatewayService) {
	ctx := context.Background()

	// This would be your application data from the database
	applicationID := int32(123)
	freelancerWallet := "0x742C4356e2B18C51EB9D0CbaF6A1A6c0C8c7DBCE"
	posterWallet := "0x8ba1f109551bD432803012645Hac136c4Ce7"
	agreedUSDAmount := int32(100)

	// Create the request
	req := payment.PostJobRequest{
		JobID:             uint64(applicationID), // Using application.id as escrow job_id
		FreelancerAddress: freelancerWallet,
		USDAmount:         fmt.Sprintf("%d", agreedUSDAmount),
		ClientAddress:     posterWallet,
	}

	// Call the payment gateway to fund escrow
	result, err := paymentGateway.PostJob(ctx, req)
	if err != nil {
		log.Printf("Failed to fund escrow: %v", err)
		return
	}

	log.Printf("Escrow funded successfully! TxHash: %s", result.TxHash)

	// In your actual code, you would also:
	// 1. Update applications.payment_status = 'deposit_initiated'
	// 2. Update applications.escrow_tx_hash_deposit = result.TxHash
	// 3. You might want to poll or use webhooks to confirm when it's mined, then update to 'deposited'
}

// Example integration for PosterReviewWork (when poster approves work)
func examplePosterReviewWorkIntegration(paymentGateway *payment.PaymentGatewayService) {
	ctx := context.Background()

	applicationID := int32(123)

	// This would be called when poster approves work AND payment_status == "deposited"
	result, err := paymentGateway.CompleteJob(ctx, uint64(applicationID))
	if err != nil {
		log.Printf("Failed to release payment: %v", err)
		return
	}

	log.Printf("Payment released successfully! TxHash: %s", result.TxHash)

	// In your actual code, you would also:
	// 1. Update applications.payment_status = 'release_initiated'
	// 2. Update applications.escrow_tx_hash_release = result.TxHash
	// 3. Once confirmed, update to 'released' and potentially close the job
}

/*
Here's how you would modify your existing handlers:

=== In your ApplicationService.RespondToOffer method ===

func (as *ApplicationService) RespondToOffer(ctx context.Context, params RespondToOfferParams) error {
	// ... existing code for database updates ...

	if params.Accept && newAppStatus == StatusHired {
		// NEW: Call payment gateway to fund escrow
		paymentGateway := payment.NewPaymentGatewayService("http://localhost:8081")

		// Get wallet addresses from database
		app, _ := as.Queries.GetApplicationByID(ctx, params.ApplicationID)
		job, _ := as.Queries.GetJobByID(ctx, app.JobID)
		applicant, _ := as.Queries.GetUserByID(ctx, app.UserID)
		poster, _ := as.Queries.GetUserByID(ctx, job.UserID)

		req := payment.PostJobRequest{
			JobID:             uint64(params.ApplicationID),
			FreelancerAddress: applicant.WalletAddress.String,
			USDAmount:         fmt.Sprintf("%d", app.AgreedUsdAmount.Int32),
			ClientAddress:     poster.WalletAddress.String,
		}

		result, err := paymentGateway.PostJob(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to fund escrow: %w", err)
		}

		// Update payment status
		err = as.InitiateEscrowDeposit(ctx, params.ApplicationID, result.TxHash)
		if err != nil {
			log.Printf("Warning: blockchain transaction successful but failed to update DB: %v", err)
		}
	}

	return tx.Commit(ctx)
}

=== In your ApplicationService.PosterReviewWork method ===

func (as *ApplicationService) PosterReviewWork(ctx context.Context, params PosterReviewWorkParams) error {
	// ... existing code for database updates ...

	if params.NewStatus == StatusWorkApproved && currentAppPaymentDetails.PaymentStatus.String == PaymentStatusDeposited {
		// NEW: Call payment gateway to release payment
		paymentGateway := payment.NewPaymentGatewayService("http://localhost:8081")

		result, err := paymentGateway.CompleteJob(ctx, uint64(params.ApplicationID))
		if err != nil {
			return fmt.Errorf("failed to release payment: %w", err)
		}

		// Update payment status
		err = as.InitiatePaymentRelease(ctx, params.ApplicationID, result.TxHash)
		if err != nil {
			log.Printf("Warning: blockchain transaction successful but failed to update DB: %v", err)
		}
	}

	return tx.Commit(ctx)
}

=== You'll also want to add these helper methods ===

// CheckTransactionStatus can be called periodically to confirm transactions
func (as *ApplicationService) CheckTransactionStatus(ctx context.Context, applicationID int32) error {
	paymentGateway := payment.NewPaymentGatewayService("http://localhost:8081")

	status, err := paymentGateway.GetJobStatus(ctx, uint64(applicationID))
	if err != nil {
		return err
	}

	// Update your database based on the blockchain status
	// This could be called from a background job or webhook

	return nil
}
*/
