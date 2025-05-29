#!/bin/bash

# Smart contract deployment script for Go integration

set -e

echo "🚀 Deploying EthJobEscrow Smart Contract..."

# Check if environment variables are set
if [ -z "$SEPOLIA_RPC_URL" ] || [ -z "$PRIVATE_KEY" ]; then
    echo "❌ Error: Please set SEPOLIA_RPC_URL and PRIVATE_KEY environment variables"
    echo "Example:"
    echo "export SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
    echo "export PRIVATE_KEY=your_private_key_without_0x"
    exit 1
fi

# Go to parent directory (where foundry project is)
cd ../

echo "📝 Compiling smart contract..."
forge build

echo "🔄 Deploying to Sepolia testnet..."
DEPLOYMENT_OUTPUT=$(forge script script/DeployPaymentGateway.s.sol:DeployEthJobEscrow \
    --rpc-url $SEPOLIA_RPC_URL \
    --private-key $PRIVATE_KEY \
    --broadcast \
    --verify \
    2>&1)

echo "$DEPLOYMENT_OUTPUT"

# Extract contract address from deployment output
CONTRACT_ADDRESS=$(echo "$DEPLOYMENT_OUTPUT" | grep -o "EthJobEscrow deployed to: 0x[a-fA-F0-9]*" | cut -d' ' -f4)

if [ -z "$CONTRACT_ADDRESS" ]; then
    echo "❌ Failed to extract contract address from deployment output"
    exit 1
fi

echo "✅ Contract deployed successfully!"
echo "📍 Contract Address: $CONTRACT_ADDRESS"
echo "🔗 Sepolia Explorer: https://sepolia.etherscan.io/address/$CONTRACT_ADDRESS"

# Go back to go-integration directory
cd go-integration/

# Update .env file if it exists
if [ -f ".env" ]; then
    # Check if CONTRACT_ADDRESS line exists, update it or add it
    if grep -q "CONTRACT_ADDRESS=" .env; then
        # Update existing line
        sed -i "s/CONTRACT_ADDRESS=.*/CONTRACT_ADDRESS=$CONTRACT_ADDRESS/" .env
        echo "📝 Updated CONTRACT_ADDRESS in .env file"
    else
        # Add new line
        echo "CONTRACT_ADDRESS=$CONTRACT_ADDRESS" >> .env
        echo "📝 Added CONTRACT_ADDRESS to .env file"
    fi
else
    echo "⚠️  No .env file found. Please create one based on env.example"
    echo "CONTRACT_ADDRESS=$CONTRACT_ADDRESS"
fi

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "Next steps:"
echo "1. Update your .env file with the contract address above"
echo "2. Ensure you have testnet ETH in your wallet"
echo "3. Run 'go run cmd/main.go' to start the server"
echo ""
echo "Test the deployment with:"
echo "curl http://localhost:8080/eth-price" 