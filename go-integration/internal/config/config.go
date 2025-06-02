package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	// Ethereum network configuration
	EthereumRPCURL  string
	NetworkID       int64
	ContractAddress string
	PrivateKey      string

	// Chainlink price feed addresses
	ETHUSDPriceFeed string

	// Application settings
	FeePercentage int
	GasLimit      uint64
	GasPrice      int64 // in Gwei

	// Database settings
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string // Constructed from individual settings

	// Server settings
	ServerPort string
}

func Load() *Config {
	cfg := &Config{
		// Default to Sepolia testnet
		EthereumRPCURL:  getEnv("ETHEREUM_RPC_URL", "https://sepolia.infura.io/v3/YOUR_INFURA_KEY"),
		NetworkID:       getEnvAsInt64("NETWORK_ID", 11155111), // Sepolia
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),

		// Sepolia ETH/USD price feed
		ETHUSDPriceFeed: getEnv("ETH_USD_PRICE_FEED", "0x694AA1769357215DE4FAC081bf1f309aDC325306"),

		FeePercentage: getEnvAsInt("FEE_PERCENTAGE", 5),
		GasLimit:      getEnvAsUint64("GAS_LIMIT", 300000),
		GasPrice:      getEnvAsInt64("GAS_PRICE", 20), // 20 Gwei

		// Database settings
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "fahed"),
		DBPassword: getEnv("DB_PASSWORD", "junglebook"),
		DBName:     getEnv("DB_NAME", "fyp-go"),

		ServerPort: getEnv("SERVER_PORT", "8081"),
	}

	// Construct database URL
	cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsUint64(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Network configurations
var Networks = map[int64]NetworkConfig{
	1: { // Mainnet
		Name:            "ethereum",
		ChainID:         1,
		ETHUSDPriceFeed: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",
		ExplorerURL:     "https://etherscan.io",
	},
	11155111: { // Sepolia
		Name:            "sepolia",
		ChainID:         11155111,
		ETHUSDPriceFeed: "0x694AA1769357215DE4FAC081bf1f309aDC325306",
		ExplorerURL:     "https://sepolia.etherscan.io",
	},
}

type NetworkConfig struct {
	Name            string
	ChainID         int64
	ETHUSDPriceFeed string
	ExplorerURL     string
}
