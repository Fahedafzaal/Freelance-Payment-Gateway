package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test default values
	cfg := Load()

	if cfg.NetworkID != 11155111 {
		t.Errorf("Expected default NetworkID to be 11155111, got %d", cfg.NetworkID)
	}

	if cfg.FeePercentage != 5 {
		t.Errorf("Expected default FeePercentage to be 5, got %d", cfg.FeePercentage)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("Expected default ServerPort to be 8080, got %s", cfg.ServerPort)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("NETWORK_ID", "1")
	os.Setenv("FEE_PERCENTAGE", "3")
	os.Setenv("SERVER_PORT", "9090")

	cfg := Load()

	if cfg.NetworkID != 1 {
		t.Errorf("Expected NetworkID to be 1, got %d", cfg.NetworkID)
	}

	if cfg.FeePercentage != 3 {
		t.Errorf("Expected FeePercentage to be 3, got %d", cfg.FeePercentage)
	}

	if cfg.ServerPort != "9090" {
		t.Errorf("Expected ServerPort to be 9090, got %s", cfg.ServerPort)
	}

	// Clean up
	os.Unsetenv("NETWORK_ID")
	os.Unsetenv("FEE_PERCENTAGE")
	os.Unsetenv("SERVER_PORT")
}

func TestNetworkConfigs(t *testing.T) {
	// Test Mainnet config
	mainnet := Networks[1]
	if mainnet.Name != "ethereum" {
		t.Errorf("Expected mainnet name to be 'ethereum', got %s", mainnet.Name)
	}

	// Test Sepolia config
	sepolia := Networks[11155111]
	if sepolia.Name != "sepolia" {
		t.Errorf("Expected sepolia name to be 'sepolia', got %s", sepolia.Name)
	}

	if sepolia.ChainID != 11155111 {
		t.Errorf("Expected sepolia ChainID to be 11155111, got %d", sepolia.ChainID)
	}
}
