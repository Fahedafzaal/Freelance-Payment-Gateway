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
    error JobNotCancelable();

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
    event JobCancelled(uint jobId, address indexed client, uint256 ethAmount);

    AggregatorV3Interface internal priceFeed;
    address public Owner;
    uint256 public constant FEE_PERCENT = 5;

    constructor(address _ethUsdPriceFeed, address owner) {
        priceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
        Owner = owner;
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

    // Mark job as completed and release payment
    function markJobCompleted(uint jobId) external {
        JobDetails storage job = jobs[jobId];

        if (msg.sender != job.client) revert OnlyClientCanMarkCompleted();
        if (job.isCompleted) revert JobAlreadyCompleted();

        job.isCompleted = true;
        emit JobCompleted(jobId);

        uint256 feeAmount = (job.ethAmount * FEE_PERCENT) / 100; // Already in wei
        uint256 freelancerAmount = job.ethAmount - feeAmount; // Already in wei

        payable(Owner).transfer(feeAmount); // No extra multiplication
        payable(job.freelancer).transfer(freelancerAmount);

        job.isPaid = true;
        emit PaymentReleased(jobId, job.freelancer, freelancerAmount);
    }

    // Cancel the job and refund ETH to the client
    function cancelJob(uint jobId) external {
        JobDetails storage job = jobs[jobId];

        if (msg.sender != job.client) revert OnlyClientCanMarkCompleted();
        if (job.isCompleted) revert JobAlreadyCompleted();
        if (job.isPaid) revert PaymentAlreadyReleased();

        uint256 refundAmount = job.ethAmount;

        // Reset the job details
        delete jobs[jobId];

        // Refund the ETH to the client
        payable(job.client).transfer(refundAmount);

        emit JobCancelled(jobId, job.client, refundAmount);
    }

    // Get job details
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
