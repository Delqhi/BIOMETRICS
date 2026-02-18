# BLOCKCHAIN.md - Blockchain Infrastructure

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the blockchain infrastructure for BIOMETRICS, enabling secure health data verification, identity management, and tokenized incentives.

## 2. Architecture

### 2.1 Network

| Component | Network | Purpose |
|-----------|---------|---------|
| Mainnet | Ethereum L1 | Production assets |
| L2 | Polygon | Scalable transactions |
| L2 | Arbitrum | DeFi integration |
| Testnet | Sepolia | Development |

### 2.2 Components

| Component | Technology | Purpose |
|-----------|------------|---------|
| Node | Geth + Lighthouse | Network access |
| Smart Contracts | Solidity | Business logic |
| Indexer | The Graph | Event indexing |
| Wallet | MetaMask/Web3.js | User interaction |

## 3. Smart Contract Development

### 3.1 Contract Structure

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract BiometricsHealthNFT is ReentrancyGuard, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant VERIFIER_ROLE = keccak256("VERIFIER_ROLE");
    
    struct HealthRecord {
        string ipfsHash;
        uint256 timestamp;
        address verifier;
        bool verified;
    }
    
    mapping(uint256 => HealthRecord) public healthRecords;
    mapping(address => uint256[]) public userRecords;
    uint256 public tokenCounter;
    
    event HealthRecordCreated(
        uint256 indexed tokenId,
        string ipfsHash,
        address indexed user,
        address verifier
    );
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    function createHealthRecord(
        address user,
        string memory ipfsHash,
        bytes memory signature
    ) public onlyRole(VERIFIER_ROLE) returns (uint256) {
        uint256 tokenId = tokenCounter++;
        
        healthRecords[tokenId] = HealthRecord({
            ipfsHash: ipfsHash,
            timestamp: block.timestamp,
            verifier: msg.sender,
            verified: true
        });
        
        userRecords[user].push(tokenId);
        
        emit HealthRecordCreated(tokenId, ipfsHash, user, msg.sender);
        
        return tokenId;
    }
    
    function getUserRecords(address user) public view returns (uint256[] memory) {
        return userRecords[user];
    }
}
```

### 3.2 Contract Deployment

```bash
# Deploy to testnet
npx hardhat run scripts/deploy.js --network sepolia

# Deploy to mainnet
npx hardhat run scripts/deploy.js --network polygon

# Verify contracts
npx hardhat verify --network polygon <CONTRACT_ADDRESS>
```

## 4. Web3 Integration

### 4.1 Wallet Connection

```typescript
import Web3 from 'web3';
import detectEthereumProvider from '@metamask/detect-provider';

class Web3Manager {
  private web3: Web3;
  private provider: any;
  
  async connect(): Promise<boolean> {
    this.provider = await detectEthereumProvider();
    
    if (!this.provider) {
      throw new Error('MetaMask not installed');
    }
    
    this.web3 = new Web3(this.provider);
    
    // Request account access
    const accounts = await this.provider.request({
      method: 'eth_requestAccounts',
    });
    
    return accounts.length > 0;
  }
  
  async signMessage(message: string): Promise<string> {
    const accounts = await this.web3.eth.getAccounts();
    
    return await this.web3.eth.personal.sign(
      message,
      accounts[0],
      ''
    );
  }
  
  async sendTransaction(tx: {
    to: string;
    value?: string;
    data?: string;
  }): Promise<string> {
    const accounts = await this.web3.eth.getAccounts();
    
    const txObj = {
      from: accounts[0],
      to: tx.to,
      value: tx.value || '0',
      data: tx.data || '0x',
      gas: await this.estimateGas(tx),
    };
    
    return await this.provider.request({
      method: 'eth_sendTransaction',
      params: [txObj],
    });
  }
}
```

### 4.2 Contract Interaction

```typescript
import HealthNFTArtifact from './artifacts/HealthNFT.json';

class ContractService {
  private contract: ethers.Contract;
  private signer: ethers.Signer;
  
  async initialize() {
    const provider = new ethers.providers.Web3Provider(window.ethereum);
    this.signer = provider.getSigner();
    
    this.contract = new ethers.Contract(
      HealthNFTArtifact.address,
      HealthNFTArtifact.abi,
      this.signer
    );
  }
  
  async createHealthRecord(user: string, ipfsHash: string): Promise<number> {
    const tx = await this.contract.createHealthRecord(user, ipfsHash);
    const receipt = await tx.wait();
    
    const event = receipt.events.find((e: any) => e.event === 'HealthRecordCreated');
    return event.args.tokenId.toNumber();
  }
  
  async getUserRecords(user: string): Promise<number[]> {
    return await this.contract.getUserRecords(user);
  }
  
  async getHealthRecord(tokenId: number): Promise<any> {
    return await this.contract.healthRecords(tokenId);
  }
}
```

## 5. IPFS Storage

### 5.1 Data Storage

```typescript
import { create } from 'ipfs-http-client';

class IPFSStorage {
  private ipfs: any;
  
  constructor() {
    this.ipfs = create({
      url: 'https://ipfs.infura.io:5001/api/v0',
      headers: {
        authorization: `Basic ${Buffer.from(
          `${process.env.IPFS_PROJECT_ID}:${process.env.IPFS_PROJECT_SECRET}`
        ).toString('base64')}`,
      },
    });
  }
  
  async upload(data: any): Promise<string> {
    const buffer = Buffer.from(JSON.stringify(data));
    
    const result = await this.ipfs.add(buffer);
    return result.path;
  }
  
  async download(hash: string): Promise<any> {
    const stream = this.ipfs.cat(hash);
    let data = '';
    
    for await (const chunk of stream) {
      data += chunk.toString();
    }
    
    return JSON.parse(data);
  }
  
  async uploadEncrypted(data: any, publicKey: string): Promise<string> {
    // Encrypt data before uploading
    const encrypted = await this.encrypt(data, publicKey);
    return this.upload(encrypted);
  }
}
```

## 6. Events & Indexing

### 6.1 The Graph

```yaml
# subgraph.yaml
specVersion: 0.0.5
schema:
  file: ./schema.graphql
dataSources:
  - kind: ethereum/contract
    name: BiometricsHealthNFT
    network: polygon
    source:
      address: "0x..."
      abi: HealthNFT
    mapping:
      kind: ethereum/events
      apiVersion: 0.0.7
      language: wasm/assemblyscript
      entities:
        - HealthRecord
        - User
      abis:
        - name: HealthNFT
          file: ./abis/HealthNFT.json
      eventHandlers:
        - event: HealthRecordCreated(indexed uint256,string,indexed address,address)
          handler: handleHealthRecordCreated
```

### 6.2 GraphQL Queries

```graphql
query GetUserHealthRecords($user: Bytes!) {
  healthRecords(where: { user: $user }) {
    id
    tokenId
    ipfsHash
    timestamp
    verified
    verifier
  }
}

query GetHealthStats($user: Bytes!, $from: Int!, $to: Int!) {
  healthRecords(
    where: { user: $user, timestamp_gte: $from, timestamp_lte: $to }
  ) {
    id
    timestamp
    verified
  }
}
```

## 7. Security

### 7.1 Security Practices

| Practice | Implementation |
|----------|---------------|
| Access Control | OpenZeppelin AccessControl |
| Reentrancy | ReentrancyGuard |
| Overflow | SafeMath / Solidity 0.8+ |
| Upgradeability | UUPS Proxy |
| Multi-sig | Gnosis Safe |

### 7.2 Audit Checklist

- [ ] Reentrancy protection
- [ ] Access control
- [ ] Input validation
- [ ] Integer overflow
- [ ] Front-running
- [ ] Flash loan attacks

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
