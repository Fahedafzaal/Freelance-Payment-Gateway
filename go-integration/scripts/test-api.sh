#!/bin/bash

# API Testing Script for Freelance Payment Gateway

set -e

BASE_URL="http://localhost:8081"

echo "🧪 Testing Freelance Payment Gateway API..."
echo "Make sure the server is running with: go run cmd/main.go"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test health endpoint
echo -e "${YELLOW}1. Testing Health Check...${NC}"
curl -s "$BASE_URL/health" && echo ""
echo ""

# Test ETH price endpoint
echo -e "${YELLOW}2. Getting Current ETH Price...${NC}"
curl -s "$BASE_URL/eth-price" | jq '.' 2>/dev/null || curl -s "$BASE_URL/eth-price"
echo ""

# Test job posting (requires actual addresses and contract)
echo -e "${YELLOW}3. Testing Job Posting (Example)...${NC}"
echo "Note: This will fail without valid addresses and deployed contract"

# Example job posting payload
JOB_PAYLOAD='{
  "job_id": 123,
  "freelancer_address": "0x742C4356e2B18C51EB9D0CbaF6A1A6c0C8c7DBCE",
  "usd_amount": "100",
  "client_address": "0x8ba1f109551bD432803012645Hac136c4Ce7"
}'

echo "Payload: $JOB_PAYLOAD"

# Uncomment to test actual job posting
# curl -X POST "$BASE_URL/post-job" \
#   -H "Content-Type: application/json" \
#   -d "$JOB_PAYLOAD" | jq '.' 2>/dev/null || echo "Job posting test skipped"

echo ""

# Test job status (requires existing job)
echo -e "${YELLOW}4. Testing Job Status Query...${NC}"
echo "Note: This will fail without an existing job"

# Uncomment to test actual job status
# curl -s "$BASE_URL/job-status?job_id=123" | jq '.' 2>/dev/null || echo "Job status test skipped"

echo ""

echo -e "${GREEN}✅ API Testing Complete!${NC}"
echo ""
echo "To run actual tests:"
echo "1. Deploy your smart contract"
echo "2. Update .env with contract address and private key"
echo "3. Start the server: go run cmd/main.go"
echo "4. Uncomment the test calls in this script"
echo "5. Run this script again" 