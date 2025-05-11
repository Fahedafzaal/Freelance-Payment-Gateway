// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test, console} from "forge-std/Test.sol";
import {MockV3Aggregator} from "lib/chainlink-brownie-contracts/contracts/src/v0.8/tests/MockV3Aggregator.sol";
import {EthJobEscrow} from "../src/PaymentGateway.sol"; // Adjust the path to your contract

contract EthJobEscrowTest is Test {
    EthJobEscrow escrow;
    MockV3Aggregator mockPriceFeed;
    address client = address(0x1);
    address freelancer = address(0x2);
    address Owner = address(0x3);
    uint jobId = 1;
    uint256 usdAmount = 1000; // USD amount to be converted to ETH
    address feeReceiver = Owner; // Use Owner as the fee receiver

    function setUp() public {
        mockPriceFeed = new MockV3Aggregator(8, 3000e8); // Simulates ETH/USD = $3000
        escrow = new EthJobEscrow(address(mockPriceFeed), Owner);
    }

    function testPostJob() public {
        uint256 requiredEth = escrow.convertUsdToEth(usdAmount);

        // Fund the client with enough ETH
        vm.deal(client, 10 ether);

        // Call postJob from client and send required ETH
        vm.prank(client);
        escrow.postJob{value: requiredEth}(
            jobId,
            freelancer,
            usdAmount,
            client
        );

        // Verify job details
        (
            address clientJob,
            address freelancerJob,
            uint256 usdJob,
            uint256 ethJob,
            bool isCompletedJob,
            bool isPaidJob
        ) = escrow.getJobDetails(jobId);

        assertEq(clientJob, client);
        assertEq(freelancerJob, freelancer);
        assertEq(usdJob, usdAmount);
        assertEq(ethJob, requiredEth);
        assertEq(isCompletedJob, false);
        assertEq(isPaidJob, false);

        // Verify that ETH is in the contract
        assertEq(address(escrow).balance, requiredEth);
    }

    function testMarkJobCompleted() public {
        uint256 requiredEth = escrow.convertUsdToEth(usdAmount);

        // Give client ETH and impersonate client
        vm.deal(client, requiredEth);
        vm.prank(client);
        escrow.postJob{value: requiredEth}(
            jobId,
            freelancer,
            usdAmount,
            client
        );

        // Verify job details
        (
            address clientJob,
            address freelancerJob,
            uint256 usdJob,
            uint256 ethJob,
            bool isCompletedJob,
            bool isPaidJob
        ) = escrow.getJobDetails(jobId);

        assertEq(clientJob, client);
        assertEq(freelancerJob, freelancer);
        assertEq(usdJob, usdAmount);
        assertEq(ethJob, requiredEth);
        assertEq(isCompletedJob, false);
        assertEq(isPaidJob, false);

        // Capture freelancer's balance before
        uint256 balanceBefore = freelancer.balance;

        // Mark job completed as client
        vm.prank(client);
        escrow.markJobCompleted(jobId);

        // Check if freelancer received the ETH (after fee)
        uint256 feeAmount = (requiredEth * 5) / 100; // 5% fee
        uint256 freelancerAmount = requiredEth - feeAmount; // 95% for freelancer
        uint256 balanceAfter = freelancer.balance;

        assertEq(balanceAfter, balanceBefore + freelancerAmount);

        // Check if fee receiver (Owner) got the fee
        uint256 feeReceiverBalance = feeReceiver.balance;
        assertEq(feeReceiverBalance, feeAmount);

        // Check updated job status
        (, , , , bool isCompletedAfter, bool isPaidAfter) = escrow
            .getJobDetails(jobId);
        assertTrue(isCompletedAfter);
        assertTrue(isPaidAfter);
    }

    function test_RevertWhen_NonClientMarksJobCompleted() public {
        uint256 requiredEth = escrow.convertUsdToEth(usdAmount);

        vm.deal(client, requiredEth);
        vm.prank(client);
        escrow.postJob{value: requiredEth}(
            jobId,
            freelancer,
            usdAmount,
            client
        );

        vm.prank(freelancer); // Non-client actor
        vm.expectRevert(
            abi.encodeWithSignature("OnlyClientCanMarkCompleted()")
        );

        escrow.markJobCompleted(jobId);
    }

    function testPostJobTwice() public {
        uint256 requiredEth = escrow.convertUsdToEth(usdAmount);

        // Fund and post the job the first time
        vm.deal(client, requiredEth * 2); // Give enough ETH for both calls, just in case
        vm.prank(client);
        escrow.postJob{value: requiredEth}(
            jobId,
            freelancer,
            usdAmount,
            client
        );

        // Expect revert on the second post with the same jobId
        vm.expectRevert(bytes("Job already exists"));
        vm.prank(client);
        escrow.postJob{value: requiredEth}(
            jobId,
            freelancer,
            usdAmount,
            client
        );
    }
}
