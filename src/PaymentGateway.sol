// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "lib/chainlink-brownie-contracts/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract EthJobEscrow {
    error InsufficientEthSent();
    error NotJobClient();
    error JobAlreadyCompleted();
    error PaymentAlreadyReleased();
    error JobNotCompleted();
    error OnlyClientCanMarkCompleted();

    event JobPosted(
        uint jobId,
        address indexed client,
        address indexed freelancer,
        uint256 usdAmount,
        uint256 ethAmount
    );
    event JobCompleted(uint jobId);
    event PaymentReleased(
        uint jobId,
        address indexed freelancer,
        uint256 ethAmount
    );

    AggregatorV3Interface internal priceFeed;

    constructor(address _ethUsdPriceFeed) {
        priceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
    }

    struct JobDetails {
        address client;
        address freelancer;
        uint256 usdAmount;
        uint256 ethAmount;
        bool isCompleted;
        bool isPaid;
    }

    mapping(uint => JobDetails) public jobs;

    // Get the latest ETH/USD conversion rate
    function getLatestEthUsd() public view returns (uint256) {
        (, int price, , , ) = priceFeed.latestRoundData();
        return uint256(price); // 8 decimals, e.g., 3000 * 10^8
    }

    // Convert USD to ETH using the latest conversion rate
    function convertUsdToEth(uint256 usdAmount) public view returns (uint256) {
        uint256 ethUsdPrice = getLatestEthUsd(); // 8 decimals
        return (usdAmount * 1e18) / ethUsdPrice; // result in wei
    }

    // Post a job with freelancer, client, and USD price
    function postJob(
        uint jobId,
        address freelancer,
        uint256 usdAmount,
        address client
    ) external payable {
        uint256 requiredEth = convertUsdToEth(usdAmount);

        if (msg.value < requiredEth) {
            revert InsufficientEthSent();
        }

        // ✅ Add this check to prevent posting the same job ID twice
        if (jobs[jobId].client != address(0)) {
            revert("Job already exists");
        }

        jobs[jobId] = JobDetails({
            client: client,
            freelancer: freelancer,
            usdAmount: usdAmount,
            ethAmount: requiredEth,
            isCompleted: false,
            isPaid: false
        });

        emit JobPosted(jobId, client, freelancer, usdAmount, requiredEth);
    }

    // Mark the job as completed and release payment to freelancer
    function markJobCompleted(uint jobId) external {
        JobDetails storage job = jobs[jobId];

        // ✅ Only client can mark job as completed
        if (msg.sender != job.client) revert OnlyClientCanMarkCompleted();

        if (job.isCompleted) revert JobAlreadyCompleted();

        job.isCompleted = true;
        emit JobCompleted(jobId);

        payable(job.freelancer).transfer(job.ethAmount);
        job.isPaid = true;

        emit PaymentReleased(jobId, job.freelancer, job.ethAmount);
    }

    function getJobDetails(
        uint jobId
    )
        external
        view
        returns (
            address client,
            address freelancer,
            uint256 usdAmount,
            uint256 ethAmount,
            bool isCompleted,
            bool isPaid
        )
    {
        return (
            jobs[jobId].client,
            jobs[jobId].freelancer,
            jobs[jobId].usdAmount,
            jobs[jobId].ethAmount,
            jobs[jobId].isCompleted,
            jobs[jobId].isPaid
        );
    }
}
