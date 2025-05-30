package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/fahedafzaal/freelance-payment-gateway/internal/config"
	"github.com/fahedafzaal/freelance-payment-gateway/pkg/blockchain"
	"github.com/fahedafzaal/freelance-payment-gateway/pkg/database"
)

type PaymentGateway struct {
	client *blockchain.Client
	config *config.Config
	db     *database.DB
}

// Request/Response types for your application flow
type PostJobRequest struct {
	JobID             uint64 `json:"job_id"`             // application.id (your escrow_job_id)
	FreelancerAddress string `json:"freelancer_address"` // applicant wallet
	USDAmount         string `json:"usd_amount"`         // agreed_usd_amount
	ClientAddress     string `json:"client_address"`     // poster wallet
}

type JobStatusResponse struct {
	JobID             uint64 `json:"job_id"`
	ApplicationID     int32  `json:"application_id"`
	FreelancerAddress string `json:"freelancer_address"`
	ClientAddress     string `json:"client_address"`
	USDAmount         string `json:"usd_amount"`
	PaymentStatus     string `json:"payment_status"`
	ApplicationStatus string `json:"application_status"`
	TxHashDeposit     string `json:"tx_hash_deposit,omitempty"`
	TxHashRelease     string `json:"tx_hash_release,omitempty"`
	TxHashRefund      string `json:"tx_hash_refund,omitempty"`
}

type TransactionResponse struct {
	TxHash      string `json:"tx_hash"`
	BlockNumber uint64 `json:"block_number"`
	GasUsed     uint64 `json:"gas_used"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

func NewPaymentGateway(cfg *config.Config) (*PaymentGateway, error) {
	// Initialize blockchain client
	client, err := blockchain.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize database connection
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PaymentGateway{
		client: client,
		config: cfg,
		db:     db,
	}, nil
}

// POST /post-job - Called when candidate accepts offer
func (pg *PaymentGateway) postJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PostJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Validate the application is ready for blockchain operations
	applicationID := int32(req.JobID) // Using application.id as escrow job_id
	if err := pg.db.ValidateApplicationForBlockchain(ctx, applicationID); err != nil {
		http.Error(w, fmt.Sprintf("Application validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Get application details from database
	details, err := pg.db.GetApplicationPaymentDetails(ctx, applicationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get application details: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify the request matches database data
	if details.ApplicantWalletAddress == nil || *details.ApplicantWalletAddress != req.FreelancerAddress {
		http.Error(w, "Freelancer address mismatch", http.StatusBadRequest)
		return
	}
	if details.PosterWalletAddress == nil || *details.PosterWalletAddress != req.ClientAddress {
		http.Error(w, "Client address mismatch", http.StatusBadRequest)
		return
	}

	// Parse addresses and amount
	freelancerAddr := common.HexToAddress(req.FreelancerAddress)
	clientAddr := common.HexToAddress(req.ClientAddress)
	usdAmount, ok := new(big.Int).SetString(req.USDAmount, 10)
	if !ok {
		http.Error(w, "Invalid USD amount", http.StatusBadRequest)
		return
	}

	// Post job to blockchain
	result, err := pg.client.PostJob(ctx, req.JobID, freelancerAddr, usdAmount, clientAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post job to blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	// Update database with transaction hash
	if err := pg.db.UpdatePaymentStatus(ctx, applicationID, "deposit_initiated", &result.TxHash, "deposit"); err != nil {
		log.Printf("Warning: Failed to update payment status in database: %v", err)
	}

	response := TransactionResponse{
		TxHash:      result.TxHash,
		BlockNumber: result.BlockNumber,
		GasUsed:     result.GasUsed,
		Success:     result.Success,
	}

	if result.Error != nil {
		response.Error = result.Error.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /complete-job?job_id=X - Called when poster approves work
func (pg *PaymentGateway) completeJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobIDStr := r.URL.Query().Get("job_id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	applicationID := int32(jobID) // application.id is used as escrow job_id

	// Get application details to verify payment status
	details, err := pg.db.GetApplicationPaymentDetails(ctx, applicationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get application details: %v", err), http.StatusInternalServerError)
		return
	}

	if details.PaymentStatus != "deposited" {
		http.Error(w, fmt.Sprintf("Cannot complete job: payment status is '%s', expected 'deposited'", details.PaymentStatus), http.StatusBadRequest)
		return
	}

	// Complete job on blockchain
	result, err := pg.client.MarkJobCompleted(ctx, jobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to complete job on blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	// Update database with release transaction hash
	if err := pg.db.UpdatePaymentStatus(ctx, applicationID, "release_initiated", &result.TxHash, "release"); err != nil {
		log.Printf("Warning: Failed to update payment status in database: %v", err)
	}

	response := TransactionResponse{
		TxHash:      result.TxHash,
		BlockNumber: result.BlockNumber,
		GasUsed:     result.GasUsed,
		Success:     result.Success,
	}

	if result.Error != nil {
		response.Error = result.Error.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /cancel-job?job_id=X - Called for refunds
func (pg *PaymentGateway) cancelJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobIDStr := r.URL.Query().Get("job_id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	applicationID := int32(jobID) // application.id is used as escrow job_id

	// Get application details to verify payment status
	details, err := pg.db.GetApplicationPaymentDetails(ctx, applicationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get application details: %v", err), http.StatusInternalServerError)
		return
	}

	if details.PaymentStatus != "deposited" {
		http.Error(w, fmt.Sprintf("Cannot cancel job: payment status is '%s', expected 'deposited'", details.PaymentStatus), http.StatusBadRequest)
		return
	}

	// Cancel job on blockchain
	result, err := pg.client.CancelJob(ctx, jobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel job on blockchain: %v", err), http.StatusInternalServerError)
		return
	}

	// Update database with refund transaction hash
	if err := pg.db.UpdatePaymentStatus(ctx, applicationID, "refund_initiated", &result.TxHash, "refund"); err != nil {
		log.Printf("Warning: Failed to update payment status in database: %v", err)
	}

	response := TransactionResponse{
		TxHash:      result.TxHash,
		BlockNumber: result.BlockNumber,
		GasUsed:     result.GasUsed,
		Success:     result.Success,
	}

	if result.Error != nil {
		response.Error = result.Error.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /job-status?job_id=X - Get application payment status
func (pg *PaymentGateway) getJobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobIDStr := r.URL.Query().Get("job_id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	applicationID := int32(jobID)

	// Get application details from database
	details, err := pg.db.GetApplicationPaymentDetails(ctx, applicationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get application details: %v", err), http.StatusInternalServerError)
		return
	}

	response := JobStatusResponse{
		JobID:             jobID,
		ApplicationID:     details.ApplicationID,
		FreelancerAddress: *details.ApplicantWalletAddress,
		ClientAddress:     *details.PosterWalletAddress,
		USDAmount:         fmt.Sprintf("%d", *details.AgreedUSDAmount),
		PaymentStatus:     details.PaymentStatus,
		ApplicationStatus: details.ApplicationStatus,
	}

	if details.EscrowTxHashDeposit != nil {
		response.TxHashDeposit = *details.EscrowTxHashDeposit
	}
	if details.EscrowTxHashRelease != nil {
		response.TxHashRelease = *details.EscrowTxHashRelease
	}
	if details.EscrowTxHashRefund != nil {
		response.TxHashRefund = *details.EscrowTxHashRefund
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /confirm-deposit?job_id=X - Called to confirm deposit (for polling/webhook)
func (pg *PaymentGateway) confirmDepositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobIDStr := r.URL.Query().Get("job_id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	applicationID := int32(jobID)

	// Update payment status to deposited
	if err := pg.db.UpdatePaymentStatus(ctx, applicationID, "deposited", nil, ""); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update payment status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// POST /confirm-release?job_id=X - Called to confirm release (for polling/webhook)
func (pg *PaymentGateway) confirmReleaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jobIDStr := r.URL.Query().Get("job_id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	applicationID := int32(jobID)

	// Update payment status to released
	if err := pg.db.UpdatePaymentStatus(ctx, applicationID, "released", nil, ""); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update payment status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// GET /eth-price - Get current ETH price
func (pg *PaymentGateway) getEthPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	price, err := pg.client.GetETHUSDPrice(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get ETH price: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"eth_usd_price": price.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.ContractAddress == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable is required")
	}
	if cfg.PrivateKey == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}
	if cfg.EthereumRPCURL == "https://sepolia.infura.io/v3/YOUR_INFURA_KEY" {
		log.Fatal("Please set a valid ETHEREUM_RPC_URL")
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize payment gateway
	gateway, err := NewPaymentGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize payment gateway: %v", err)
	}
	defer gateway.client.Close()
	defer gateway.db.Close()

	// Setup HTTP routes for your application flow
	http.HandleFunc("/post-job", gateway.postJobHandler)               // Offer accepted → fund escrow
	http.HandleFunc("/complete-job", gateway.completeJobHandler)       // Work approved → release payment
	http.HandleFunc("/cancel-job", gateway.cancelJobHandler)           // Cancel/refund
	http.HandleFunc("/job-status", gateway.getJobStatusHandler)        // Get payment status
	http.HandleFunc("/confirm-deposit", gateway.confirmDepositHandler) // Confirm deposit completion
	http.HandleFunc("/confirm-release", gateway.confirmReleaseHandler) // Confirm release completion
	http.HandleFunc("/eth-price", gateway.getEthPriceHandler)          // Current ETH price

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting payment gateway server on port %s", cfg.ServerPort)
	log.Printf("Contract address: %s", cfg.ContractAddress)
	log.Printf("Network ID: %d", cfg.NetworkID)
	log.Printf("Database connected successfully")

	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
