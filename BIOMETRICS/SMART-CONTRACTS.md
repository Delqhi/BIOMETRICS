# SMART-CONTRACTS.md - Smart Contract Development

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the smart contract development practices for BIOMETRICS, enabling secure, auditable, and upgradeable blockchain applications.

## 2. Contract Architecture

### 2.1 Contract Structure

```
contracts/
├── interfaces/
│   ├── IHealthNFT.sol
│   └── IToken.sol
├── tokens/
│   ├── BiometricsToken.sol
│   └── RewardToken.sol
├── nft/
│   └── HealthNFT.sol
├── staking/
│   └── StakingPool.sol
├── governance/
│   └── Governance.sol
├── mocks/
│   └── MockERC20.sol
└── utils/
    └── MerkleDistributor.sol
```

### 2.2 Base Contracts

```solidity
// contracts/utils/BaseContract.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

abstract contract BaseContract is ReentrancyGuard, Pausable, AccessControl {
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");
    
    event Paused(address account);
    event Unpaused(address account);
    
    modifier whenNotPaused() {
        require(!paused(), "Pausable: paused");
        _;
    }
    
    modifier whenPaused() {
        require(paused(), "Pausable: not paused");
        _;
    }
    
    function pause() public onlyRole(PAUSER_ROLE) whenNotPaused {
        _pause();
        emit Paused(msg.sender);
    }
    
    function unpause() public onlyRole(PAUSER_ROLE) whenPaused {
        _unpause();
        emit Unpaused(msg.sender);
    }
}
```

## 3. Health NFT Contract

### 3.1 NFT Implementation

```solidity
// contracts/nft/HealthNFT.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "./BaseContract.sol";

contract HealthNFT is ERC721, ERC721URIStorage, ERC721Burnable, BaseContract {
    using Counters for Counters.Counter;
    
    Counters.Counter private _tokenIdCounter;
    
    struct HealthData {
        string dataHash;
        uint256 timestamp;
        address verifier;
        uint8 dataType;
    }
    
    mapping(uint256 => HealthData) public healthData;
    mapping(address => uint256[]) public userTokens;
    mapping(bytes32 => bool) public dataHashes;
    
    event HealthDataStored(
        uint256 indexed tokenId,
        bytes32 dataHash,
        address indexed user,
        address verifier
    );
    
    constructor() ERC721("BIOMETRICS Health", "BHT") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(PAUSER_ROLE, msg.sender);
    }
    
    function storeHealthData(
        address to,
        string memory _tokenURI,
        string memory _dataHash,
        uint8 _dataType
    ) public onlyRole(OPERATOR_ROLE) returns (uint256) {
        require(bytes(_dataHash).length > 0, "Invalid data hash");
        
        uint256 tokenId = _tokenIdCounter.current();
        _tokenIdCounter.increment();
        
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, _tokenURI);
        
        healthData[tokenId] = HealthData({
            dataHash: _dataHash,
            timestamp: block.timestamp,
            verifier: msg.sender,
            dataType: _dataType
        });
        
        dataHashes[keccak256(bytes(_dataHash))] = true;
        userTokens[to].push(tokenId);
        
        emit HealthDataStored(tokenId, keccak256(bytes(_dataHash)), to, msg.sender);
        
        return tokenId;
    }
    
    function verifyHealthData(
        uint256 tokenId,
        string memory _dataHash
    ) public onlyRole(OPERATOR_ROLE) returns (bool) {
        require(ownerOf(tokenId) != address(0), "Invalid token");
        
        bytes32 inputHash = keccak256(bytes(_dataHash));
        
        return dataHashes[inputHash];
    }
    
    function getUserTokens(address user) public view returns (uint256[] memory) {
        return userTokens[user];
    }
    
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721, AccessControl)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
```

## 4. Token Contract

### 4.1 Reward Token

```solidity
// contracts/tokens/RewardToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract RewardToken is ERC20, ERC20Burnable, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant BURNER_ROLE = keccak256("BURNER_ROLE");
    
    mapping(address => bool) public minters;
    uint256 public constant MAX_SUPPLY = 1000000000 * 10**18;
    
    event MinterAdded(address indexed minter);
    event MinterRemoved(address indexed minter);
    
    constructor() ERC20("BIOMETRICS Reward", "BIO") {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
    }
    
    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        require(
            totalSupply() + amount <= MAX_SUPPLY,
            "Max supply exceeded"
        );
        _mint(to, amount);
    }
    
    function burn(address from, uint256 amount) public onlyRole(BURNER_ROLE) {
        _burn(from, amount);
    }
    
    function addMinter(address minter) public onlyRole(DEFAULT_ADMIN_ROLE) {
        grantRole(MINTER_ROLE, minter);
        emit MinterAdded(minter);
    }
    
    function removeMinter(address minter) public onlyRole(DEFAULT_ADMIN_ROLE) {
        revokeRole(MINTER_ROLE, minter);
        emit MinterRemoved(minter);
    }
}
```

## 5. Staking Contract

### 5.1 Staking Pool

```solidity
// contracts/staking/StakingPool.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract StakingPool is ReentrancyGuard, AccessControl {
    using SafeERC20 for IERC20;
    
    IERC20 public stakingToken;
    IERC20 public rewardToken;
    
    struct Stake {
        uint256 amount;
        uint256 startTime;
        uint256 rewards;
        uint256 lastClaimTime;
    }
    
    mapping(address => Stake[]) public stakes;
    mapping(address => uint256) public totalStaked;
    
    uint256 public rewardRate = 100; // 100 tokens per second per 1000 staked
    uint256 public constant REWARD_RATE_BASE = 1000;
    
    event Staked(address indexed user, uint256 amount, uint256 stakeId);
    event Unstaked(address indexed user, uint256 amount, uint256 stakeId);
    event RewardClaimed(address indexed user, uint256 reward);
    
    constructor(IERC20 _stakingToken, IERC20 _rewardToken) {
        stakingToken = _stakingToken;
        rewardToken = _rewardToken;
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    function stake(uint256 amount) external nonReentrant {
        require(amount > 0, "Cannot stake 0");
        
        stakingToken.safeTransferFrom(msg.sender, address(this), amount);
        
        stakes[msg.sender].push(Stake({
            amount: amount,
            startTime: block.timestamp,
            rewards: 0,
            lastClaimTime: block.timestamp
        }));
        
        totalStaked[msg.sender] += amount;
        
        emit Staked(msg.sender, amount, stakes[msg.sender].length - 1);
    }
    
    function unstake(uint256 stakeId) external nonReentrant {
        require(stakes[msg.sender].length > stakeId, "Invalid stake");
        
        Stake storage stake_ = stakes[msg.sender][stakeId];
        require(stake_.amount > 0, "Already unstaked");
        
        // Calculate pending rewards
        uint256 pending = calculatePendingRewards(msg.sender, stakeId);
        stake_.rewards += pending;
        
        // Transfer tokens
        stakingToken.safeTransfer(msg.sender, stake_.amount);
        
        if (stake_.rewards > 0) {
            rewardToken.safeTransfer(msg.sender, stake_.rewards);
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
        
        uint256 timeStaked = block.timestamp - stake_.lastClaimTime;
        uint256 reward = (stake_.amount * rewardRate * timeStaked) / REWARD_RATE_BASE;
        
        return reward;
    }
}
```

## 6. Testing

### 6.1 Test Structure

```typescript
// test/HealthNFT.test.ts
import { expect } from 'chai';
import { ethers } from 'hardhat';

describe('HealthNFT', function () {
  let healthNFT: any;
  let owner: any;
  let user: any;
  
  beforeEach(async function () {
    [owner, user] = await ethers.getSigners();
    
    const HealthNFT = await ethers.getContractFactory('HealthNFT');
    healthNFT = await HealthNFT.deploy();
  });
  
  it('should store health data', async function () {
    await healthNFT.storeHealthData(
      user.address,
      'ipfs://QmHash',
      'dataHash123',
      1
    );
    
    const balance = await healthNFT.balanceOf(user.address);
    expect(balance).to.equal(1);
  });
  
  it('should allow verification', async function () {
    await healthNFT.storeHealthData(
      user.address,
      'ipfs://QmHash',
      'dataHash123',
      1
    );
    
    const isValid = await healthNFT.verifyHealthData(1, 'dataHash123');
    expect(isValid).to.equal(true);
  });
});
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
