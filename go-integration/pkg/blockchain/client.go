package blockchain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/fahedafzaal/freelance-payment-gateway/contracts"
	"github.com/fahedafzaal/freelance-payment-gateway/internal/config"
)

type Client struct {
	ethClient       *ethclient.Client
	contract        *contracts.EthJobEscrow
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	publicAddress   common.Address
	config          *config.Config
}

type JobDetails struct {
	Client      common.Address
	Freelancer  common.Address
	USDAmount   *big.Int
	ETHAmount   *big.Int
	IsCompleted bool
	IsPaid      bool
}

type TransactionResult struct {
	TxHash      string
	BlockNumber uint64
	GasUsed     uint64
	Success     bool
	Error       error
}

// NewClient creates a new blockchain client instance
func NewClient(cfg *config.Config) (*Client, error) {
	// Connect to Ethereum client
	ethClient, err := ethclient.Dial(cfg.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	// Get public address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Connect to smart contract
	contractAddress := common.HexToAddress(cfg.ContractAddress)
	contract, err := contracts.NewEthJobEscrow(contractAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		ethClient:       ethClient,
		contract:        contract,
		contractAddress: contractAddress,
		privateKey:      privateKey,
		publicAddress:   publicAddress,
		config:          cfg,
	}, nil
}

// GetAuth creates a new transactor for sending transactions
func (c *Client) GetAuth(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, c.publicAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	chainID, err := c.ethClient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = c.config.GasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// PostJob creates a new job on the blockchain
func (c *Client) PostJob(ctx context.Context, jobID uint64, freelancer common.Address, usdAmount *big.Int, client common.Address) (*TransactionResult, error) {
	// Get current ETH price and calculate required ETH
	ethAmount, err := c.contract.ConvertUsdToEth(&bind.CallOpts{Context: ctx}, usdAmount)
	if err != nil {
		return nil, err
	}

	// Get transaction options
	auth, err := c.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	// Set the value to send (ETH amount)
	auth.Value = ethAmount

	// Execute transaction
	tx, err := c.contract.PostJob(auth, big.NewInt(int64(jobID)), freelancer, usdAmount, client)
	if err != nil {
		return &TransactionResult{
			Success: false,
			Error:   err,
		}, err
	}

	// Wait for transaction confirmation
	return c.waitForTransaction(ctx, tx)
}

// MarkJobCompleted marks a job as completed and releases payment
func (c *Client) MarkJobCompleted(ctx context.Context, jobID uint64) (*TransactionResult, error) {
	auth, err := c.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := c.contract.MarkJobCompleted(auth, big.NewInt(int64(jobID)))
	if err != nil {
		return &TransactionResult{
			Success: false,
			Error:   err,
		}, err
	}

	return c.waitForTransaction(ctx, tx)
}

// CancelJob cancels a job and refunds the client
func (c *Client) CancelJob(ctx context.Context, jobID uint64) (*TransactionResult, error) {
	auth, err := c.GetAuth(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := c.contract.CancelJob(auth, big.NewInt(int64(jobID)))
	if err != nil {
		return &TransactionResult{
			Success: false,
			Error:   err,
		}, err
	}

	return c.waitForTransaction(ctx, tx)
}

// GetJobDetails retrieves job information from the blockchain
func (c *Client) GetJobDetails(ctx context.Context, jobID uint64) (*JobDetails, error) {
	result, err := c.contract.GetJobDetails(
		&bind.CallOpts{Context: ctx},
		big.NewInt(int64(jobID)),
	)
	if err != nil {
		return nil, err
	}

	return &JobDetails{
		Client:      result.Client,
		Freelancer:  result.Freelancer,
		USDAmount:   result.UsdAmount,
		ETHAmount:   result.EthAmount,
		IsCompleted: result.IsCompleted,
		IsPaid:      result.IsPaid,
	}, nil
}

// GetETHUSDPrice gets the current ETH/USD price from Chainlink
func (c *Client) GetETHUSDPrice(ctx context.Context) (*big.Int, error) {
	return c.contract.GetLatestEthUsd(&bind.CallOpts{Context: ctx})
}

// ConvertUSDToETH converts USD amount to ETH using current price
func (c *Client) ConvertUSDToETH(ctx context.Context, usdAmount *big.Int) (*big.Int, error) {
	return c.contract.ConvertUsdToEth(&bind.CallOpts{Context: ctx}, usdAmount)
}

// waitForTransaction waits for transaction confirmation and returns result
func (c *Client) waitForTransaction(ctx context.Context, tx *types.Transaction) (*TransactionResult, error) {
	log.Printf("Transaction sent: %s", tx.Hash().Hex())

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return &TransactionResult{
			TxHash:  tx.Hash().Hex(),
			Success: false,
			Error:   err,
		}, err
	}

	success := receipt.Status == types.ReceiptStatusSuccessful

	return &TransactionResult{
		TxHash:      tx.Hash().Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
		Success:     success,
		Error:       nil,
	}, nil
}

// GetBalance gets ETH balance for an address
func (c *Client) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.ethClient.BalanceAt(ctx, address, nil)
}

// Close closes the Ethereum client connection
func (c *Client) Close() {
	c.ethClient.Close()
}
