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
)

type PaymentGateway struct {
	client *blockchain.Client
	config *config.Config
}

type PostJobRequest struct {
	JobID      uint64 `json:"job_id"`
	Freelancer string `json:"freelancer_address"`
	USDAmount  string `json:"usd_amount"`
	ClientAddr string `json:"client_address"`
}

type JobStatusResponse struct {
	JobID       uint64 `json:"job_id"`
	Client      string `json:"client"`
	Freelancer  string `json:"freelancer"`
	USDAmount   string `json:"usd_amount"`
	ETHAmount   string `json:"eth_amount"`
	IsCompleted bool   `json:"is_completed"`
	IsPaid      bool   `json:"is_paid"`
}

type TransactionResponse struct {
	TxHash      string `json:"tx_hash"`
	BlockNumber uint64 `json:"block_number"`
	GasUsed     uint64 `json:"gas_used"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

func NewPaymentGateway(cfg *config.Config) (*PaymentGateway, error) {
	client, err := blockchain.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &PaymentGateway{
		client: client,
		config: cfg,
	}, nil
}

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

	// Parse addresses
	freelancerAddr := common.HexToAddress(req.Freelancer)
	clientAddr := common.HexToAddress(req.ClientAddr)

	// Parse USD amount (assume it's in dollars, convert to wei-like format)
	usdAmount, ok := new(big.Int).SetString(req.USDAmount, 10)
	if !ok {
		http.Error(w, "Invalid USD amount", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := pg.client.PostJob(ctx, req.JobID, freelancerAddr, usdAmount, clientAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to post job: %v", err), http.StatusInternalServerError)
		return
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

	result, err := pg.client.MarkJobCompleted(ctx, jobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to complete job: %v", err), http.StatusInternalServerError)
		return
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

	result, err := pg.client.CancelJob(ctx, jobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel job: %v", err), http.StatusInternalServerError)
		return
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

	jobDetails, err := pg.client.GetJobDetails(ctx, jobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get job details: %v", err), http.StatusInternalServerError)
		return
	}

	response := JobStatusResponse{
		JobID:       jobID,
		Client:      jobDetails.Client.Hex(),
		Freelancer:  jobDetails.Freelancer.Hex(),
		USDAmount:   jobDetails.USDAmount.String(),
		ETHAmount:   jobDetails.ETHAmount.String(),
		IsCompleted: jobDetails.IsCompleted,
		IsPaid:      jobDetails.IsPaid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

	// Initialize payment gateway
	gateway, err := NewPaymentGateway(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize payment gateway: %v", err)
	}
	defer gateway.client.Close()

	// Setup HTTP routes
	http.HandleFunc("/post-job", gateway.postJobHandler)
	http.HandleFunc("/complete-job", gateway.completeJobHandler)
	http.HandleFunc("/cancel-job", gateway.cancelJobHandler)
	http.HandleFunc("/job-status", gateway.getJobStatusHandler)
	http.HandleFunc("/eth-price", gateway.getEthPriceHandler)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting payment gateway server on port %s", cfg.ServerPort)
	log.Printf("Contract address: %s", cfg.ContractAddress)
	log.Printf("Network ID: %d", cfg.NetworkID)

	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
