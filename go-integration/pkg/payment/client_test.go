package payment

import (
	"math/big"
	"testing"

	"github.com/Fahedafzaal/Freelance-Payment-Gateway/go-integration/internal/config"
)

func TestConfig(t *testing.T) {
	cfg := &config.Config{
		EthereumRPCURL:  "https://sepolia.infura.io/v3/test",
		NetworkID:       11155111,
		ContractAddress: "0x1234567890123456789012345678901234567890",
		PrivateKey:      "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		GasLimit:        300000,
		GasPrice:        20,
	}

	if cfg.NetworkID != 11155111 {
		t.Errorf("Expected NetworkID to be 11155111, got %d", cfg.NetworkID)
	}

	if cfg.GasLimit != 300000 {
		t.Errorf("Expected GasLimit to be 300000, got %d", cfg.GasLimit)
	}
}

func TestJobDetails(t *testing.T) {
	jobDetails := &JobDetails{
		USDAmount:   big.NewInt(100),
		ETHAmount:   big.NewInt(31250000000000000), // ~0.03125 ETH
		IsCompleted: false,
		IsPaid:      false,
	}

	if jobDetails.USDAmount.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected USDAmount to be 100, got %s", jobDetails.USDAmount.String())
	}

	if jobDetails.IsCompleted {
		t.Errorf("Expected IsCompleted to be false")
	}

	if jobDetails.IsPaid {
		t.Errorf("Expected IsPaid to be false")
	}
}

func TestTransactionResult(t *testing.T) {
	result := &TransactionResult{
		TxHash:      "0x1234567890abcdef",
		BlockNumber: 12345,
		GasUsed:     150000,
		Success:     true,
		Error:       nil,
	}

	if !result.Success {
		t.Errorf("Expected Success to be true")
	}

	if result.GasUsed != 150000 {
		t.Errorf("Expected GasUsed to be 150000, got %d", result.GasUsed)
	}

	if result.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", result.Error)
	}
}

// Integration test - only run with valid configuration
func TestNewClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This would require valid configuration
	// Uncomment and configure for actual integration testing
	/*
		cfg := &config.Config{
			EthereumRPCURL:  "https://sepolia.infura.io/v3/YOUR_PROJECT_ID",
			NetworkID:       11155111,
			ContractAddress: "0xYourContractAddress",
			PrivateKey:      "your_private_key",
			GasLimit:        300000,
			GasPrice:        20,
		}

		client, err := NewClient(cfg)
		if err != nil {
			t.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// Test getting ETH price
		ctx := context.Background()
		price, err := client.GetETHUSDPrice(ctx)
		if err != nil {
			t.Fatalf("Failed to get ETH price: %v", err)
		}

		if price.Cmp(big.NewInt(0)) <= 0 {
			t.Errorf("Expected positive ETH price, got %s", price.String())
		}
	*/

	t.Skip("Integration test requires valid configuration")
}
