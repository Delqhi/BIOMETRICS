# TOKEN-ECONOMY.md - Token Economy Design

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Token Overview

### 1.1 Token Utility

| Token | Symbol | Purpose | Type |
|-------|--------|---------|------|
| BIOMETRICS | BIO | Governance | ERC-20 |
| Health Points | HP | Rewards | ERC-20 |

### 1.2 Token Distribution

| Category | Allocation | Vesting |
|----------|------------|---------|
| Community | 40% | 4 years |
| Team | 20% | 3 years |
| Investors | 15% | 2 years |
| Treasury | 15% | 4 years |
| Rewards | 10% | 3 years |

## 2. Token Contract

### 2.1 BIO Token

```solidity
// contracts/tokens/BIOToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Votes.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract BIOToken is ERC20, ERC20Burnable, ERC20Votes, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
    
    uint256 public constant MAX_SUPPLY = 1000000000e18; // 1B tokens
    
    mapping(address => bool) public minters;
    mapping(address => bool) public blacklisted;
    
    event MinterAdded(address indexed minter);
    event MinterRemoved(address indexed minter);
    event Blacklisted(address indexed account);
    event Unblacklisted(address indexed account);
    
    constructor() ERC20("BIOMETRICS Token", "BIO") ERC20Permit("BIOMETRICS Token") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
    }
    
    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        require(totalSupply() + amount <= MAX_SUPPLY, "Max supply exceeded");
        _mint(to, amount);
    }
    
    function burn(uint256 amount) public override onlyRole(BURNER_ROLE) {
        _burn(msg.sender, amount);
    }
    
    function burnFrom(address from, uint256 amount) public override {
        require(!blacklisted[from], "Blacklisted");
        super.burnFrom(from, amount);
    }
    
    function transfer(address to, uint256 amount) public override returns (bool) {
        require(!blacklisted[msg.sender], "Blacklisted");
        return super.transfer(to, amount);
    }
    
    function _afterTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal override(ERC20, ERC20Votes) {
        super._afterTokenTransfer(from, to, amount);
    }
}
```

## 3. Reward System

### 3.1 Health Points

```solidity
// contracts/tokens/HealthPoints.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract HealthPoints is ERC20, AccessControl {
    bytes32 public constant DISTRIBUTOR_ROLE = keccak256("DISTRIBUTOR_ROLE");
    
    uint256 public constant MAX_SUPPLY = 10000000000e18; // 10B HP
    
    mapping(address => uint256) public lastClaimTime;
    uint256 public constant CLAIM_INTERVAL = 1 days;
    uint256 public dailyRewardRate = 100e18; // 100 HP per day
    
    event RewardsClaimed(address indexed user, uint256 amount);
    event RewardRateUpdated(uint256 newRate);
    
    constructor() ERC20("Health Points", "HP") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(DISTRIBUTOR_ROLE, msg.sender);
    }
    
    function claimRewards() external {
        require(
            block.timestamp >= lastClaimTime[msg.sender] + CLAIM_INTERVAL,
            "Claim not available yet"
        );
        
        uint256 rewards = calculateRewards(msg.sender);
        
        require(totalSupply() + rewards <= MAX_SUPPLY, "Max supply exceeded");
        
        lastClaimTime[msg.sender] = block.timestamp;
        _mint(msg.sender, rewards);
        
        emit RewardsClaimed(msg.sender, rewards);
    }
    
    function calculateRewards(address user) public view returns (uint256) {
        uint256 daysSinceLastClaim = (block.timestamp - lastClaimTime[user]) / CLAIM_INTERVAL;
        return daysSinceLastClaim * dailyRewardRate;
    }
    
    function distributeRewards(address[] calldata users, uint256[] calldata amounts)
        external onlyRole(DISTRIBUTOR_ROLE)
    {
        require(users.length == amounts.length, "Length mismatch");
        
        for (uint256 i = 0; i < users.length; i++) {
            _mint(users[i], amounts[i]);
        }
    }
}
```

## 4. Incentive Programs

### 4.1 Staking Rewards

```typescript
const STAKING_REWARDS = {
  '30_days': 1.1,   // 10% APY
  '60_days': 1.25,   // 25% APY
  '90_days': 1.5,    // 50% APY
  '180_days': 2.0,   // 100% APY
  '365_days': 3.0,   // 200% APY
};
```

### 4.2 Activity Rewards

| Activity | HP Reward |
|----------|-----------|
| Daily login | 10 HP |
| Track health data | 5 HP |
| Complete challenge | 50 HP |
| Refer friend | 100 HP |
| Community contribution | 200 HP |

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
