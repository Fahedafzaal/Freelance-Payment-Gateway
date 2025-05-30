package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

// ApplicationPaymentDetails represents payment-related data from your existing schema
type ApplicationPaymentDetails struct {
	ApplicationID          int32
	JobID                  int32
	ApplicantUserID        int32
	PosterUserID           int32
	AgreedUSDAmount        *int32
	PaymentStatus          string
	EscrowJobID            *int32
	EscrowTxHashDeposit    *string
	EscrowTxHashRelease    *string
	EscrowTxHashRefund     *string
	ApplicantWalletAddress *string
	PosterWalletAddress    *string
	ApplicationStatus      string
}

// NewDB creates a new database connection using pgx
func NewDB(connStr string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &DB{Pool: pool}, nil
}

// GetApplicationPaymentDetails retrieves application and payment details for blockchain operations
func (db *DB) GetApplicationPaymentDetails(ctx context.Context, applicationID int32) (*ApplicationPaymentDetails, error) {
	query := `
		SELECT 
			a.id as application_id,
			a.job_id,
			a.user_id as applicant_user_id,
			j.user_id as poster_user_id,
			a.agreed_usd_amount,
			COALESCE(a.payment_status, 'pending_deposit') as payment_status,
			a.escrow_job_id,
			a.escrow_tx_hash_deposit,
			a.escrow_tx_hash_release,
			a.escrow_tx_hash_refund,
			applicant.wallet_address as applicant_wallet_address,
			poster.wallet_address as poster_wallet_address,
			a.status as application_status
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN users applicant ON a.user_id = applicant.id
		JOIN users poster ON j.user_id = poster.id
		WHERE a.id = $1
	`

	details := &ApplicationPaymentDetails{}
	err := db.Pool.QueryRow(ctx, query, applicationID).Scan(
		&details.ApplicationID,
		&details.JobID,
		&details.ApplicantUserID,
		&details.PosterUserID,
		&details.AgreedUSDAmount,
		&details.PaymentStatus,
		&details.EscrowJobID,
		&details.EscrowTxHashDeposit,
		&details.EscrowTxHashRelease,
		&details.EscrowTxHashRefund,
		&details.ApplicantWalletAddress,
		&details.PosterWalletAddress,
		&details.ApplicationStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying application payment details: %v", err)
	}

	return details, nil
}

// UpdatePaymentStatus updates the payment status and transaction hash
func (db *DB) UpdatePaymentStatus(ctx context.Context, applicationID int32, status string, txHash *string, txType string) error {
	var query string
	var args []interface{}

	switch txType {
	case "deposit":
		query = `
			UPDATE applications 
			SET payment_status = $1, escrow_tx_hash_deposit = $2
			WHERE id = $3
		`
		args = []interface{}{status, txHash, applicationID}
	case "release":
		query = `
			UPDATE applications 
			SET payment_status = $1, escrow_tx_hash_release = $2
			WHERE id = $3
		`
		args = []interface{}{status, txHash, applicationID}
	case "refund":
		query = `
			UPDATE applications 
			SET payment_status = $1, escrow_tx_hash_refund = $2
			WHERE id = $3
		`
		args = []interface{}{status, txHash, applicationID}
	default:
		query = `
			UPDATE applications 
			SET payment_status = $1
			WHERE id = $2
		`
		args = []interface{}{status, applicationID}
	}

	_, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating payment status: %v", err)
	}

	return nil
}

// ValidateApplicationForBlockchain checks if application is ready for blockchain operations
func (db *DB) ValidateApplicationForBlockchain(ctx context.Context, applicationID int32) error {
	var status string
	var applicantWallet, posterWallet *string
	var agreedAmount *int32

	query := `
		SELECT 
			a.status,
			a.agreed_usd_amount,
			applicant.wallet_address as applicant_wallet,
			poster.wallet_address as poster_wallet
		FROM applications a
		JOIN jobs j ON a.job_id = j.id
		JOIN users applicant ON a.user_id = applicant.id
		JOIN users poster ON j.user_id = poster.id
		WHERE a.id = $1
	`

	err := db.Pool.QueryRow(ctx, query, applicationID).Scan(
		&status,
		&agreedAmount,
		&applicantWallet,
		&posterWallet,
	)
	if err != nil {
		return fmt.Errorf("application not found: %v", err)
	}

	if applicantWallet == nil || *applicantWallet == "" {
		return fmt.Errorf("applicant wallet address not set")
	}

	if posterWallet == nil || *posterWallet == "" {
		return fmt.Errorf("poster wallet address not set")
	}

	if agreedAmount == nil || *agreedAmount <= 0 {
		return fmt.Errorf("agreed USD amount not set or invalid")
	}

	return nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	db.Pool.Close()
}
