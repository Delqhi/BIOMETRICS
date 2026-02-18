# STAKING.md - Staking System

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the staking system for BIOMETRICS token, enabling token holders to earn rewards while securing the network.

## 2. Staking Contract

### 2.1 Main Staking

```solidity
// contracts/staking/BIOStaking.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract BIOStaking is ReentrancyGuard, AccessControl {
    using SafeERC20 for IERC20;
    
    IERC20 public stakingToken;
    IERC20 public rewardToken;
    
    struct Stake {
        uint256 amount;
        uint256 startTime;
        uint256 lockPeriod;
        uint256 rewards;
    }
    
    mapping(address => Stake[]) public stakes;
    mapping(address => uint256) public totalStaked;
    mapping(address => uint256) public totalRewards;
    
    uint256 public constant LOCK_30_DAYS = 30 days;
    uint256 public constant LOCK_60_DAYS = 60 days;
    uint256 public constant LOCK_90_DAYS = 90 days;
    uint256 public constant LOCK_180_DAYS = 180 days;
    uint256 public constant LOCK_365_DAYS = 365 days;
    
    uint256[5] public rewardRates = [110, 125, 150, 200, 300]; // APY * 10
    uint256[5] public lockPeriods = [30 days, 60 days, 90 days, 180 days, 365 days];
    
    event Staked(
        address indexed user,
        uint256 amount,
        uint256 lockPeriod,
        uint256 stakeId
    );
    event Unstaked(
        address indexed user,
        uint256 amount,
        uint256 stakeId
    );
    event RewardClaimed(
        address indexed user,
        uint256 reward
    );
    
    constructor(IERC20 _stakingToken, IERC20 _rewardToken) {
        stakingToken = _stakingToken;
        rewardToken = _rewardToken;
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    function stake(uint256 amount, uint256 lockPeriodIndex) 
        external 
        nonReentrant 
    {
        require(amount > 0, "Cannot stake 0");
        require(lockPeriodIndex < 5, "Invalid lock period");
        
        stakingToken.safeTransferFrom(msg.sender, address(this), amount);
        
        uint256 stakeId = stakes[msg.sender].length;
        
        stakes[msg.sender].push(Stake({
            amount: amount,
            startTime: block.timestamp,
            lockPeriod: lockPeriods[lockPeriodIndex],
            rewards: 0
        }));
        
        totalStaked[msg.sender] += amount;
        
        emit Staked(msg.sender, amount, lockPeriodIndex, stakeId);
    }
    
    function unstake(uint256 stakeId) external nonReentrant {
        require(stakes[msg.sender].length > stakeId, "Invalid stake");
        
        Stake storage stake_ = stakes[msg.sender][stakeId];
        require(stake_.amount > 0, "Already unstaked");
        require(
            block.timestamp >= stake_.startTime + stake_.lockPeriod,
            "Lock period not over"
        );
        
        // Calculate pending rewards
        uint256 pending = calculatePendingRewards(msg.sender, stakeId);
        stake_.rewards += pending;
        
        // Transfer staked tokens
        stakingToken.safeTransfer(msg.sender, stake_.amount);
        
        // Transfer rewards
        if (stake_.rewards > 0) {
            rewardToken.safeTransfer(msg.sender, stake_.rewards);
            totalRewards[msg.sender] += stake_.rewards;
            emit RewardClaimed(msg.sender, stake_.rewards);
        }
        
        totalStaked[msg.sender] -= stake_.amount;
        stake_.amount = 0;
        
        emit Unstaked(msg.sender, stake_.amount, stakeId);
    }
    
    function calculatePendingRewards(
        address user,
        uint256 stakeId
    ) public view returns (uint256) {
        Stake storage stake_ = stakes[user][stakeId];
        
        if (stake_.amount == 0) return 0;
        
        uint256 lockIndex = getLockIndex(stake_.lockPeriod);
        uint256 apy = rewardRates[lockIndex];
        
        uint256 timeStaked = block.timestamp - stake_.startTime;
        uint256 rewards = (stake_.amount * apy * timeStaked) / (1000 * 365 days);
        
        return rewards;
    }
    
    function getLockIndex(uint256 lockPeriod) internal pure returns (uint256) {
        for (uint256 i = 0; i < lockPeriods.length; i++) {
            if (lockPeriods[i] == lockPeriod) return i;
        }
        return 0;
    }
}
```

## 2. Reward Distribution

### 2.1 Reward Pool

```typescript
class StakingRewards {
  async distributeRewards() {
    const rewardDistribution = await ethers.getContractFactory('RewardDistribution');
    const contract = await rewardDistribution.deploy(
      TREASURY_ADDRESS,
      REWARD_TOKEN_ADDRESS,
      STAKING_TOKEN_ADDRESS
    );
    
    return contract.address;
  }
}
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
