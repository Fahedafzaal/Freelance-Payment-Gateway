// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Script.sol";
import "../src/PaymentGateway.sol";

contract DeployEthJobEscrow is Script {
    function run() external {
        address priceFeed;

        uint256 chainId = block.chainid;

        // Select price feed based on chain ID
        if (chainId == 1) {
            // Ethereum Mainnet
            priceFeed = 0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419;
        } else if (chainId == 11155111) {
            // Sepolia Testnet
            priceFeed = 0x694AA1769357215DE4FAC081bf1f309aDC325306;
        } else {
            revert("Unsupported network");
        }

        // Broadcast deployment
        vm.startBroadcast();

        EthJobEscrow escrow = new EthJobEscrow(priceFeed);

        vm.stopBroadcast();

        console.log("EthJobEscrow deployed to:", address(escrow));
    }
}
