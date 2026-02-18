# CRYPTO-PAYMENT.md - Cryptocurrency Payment Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the cryptocurrency payment integration for BIOMETRICS, enabling users to pay with crypto for subscriptions and services.

## 2. Supported Currencies

### 2.1 Supported Tokens

| Token | Network | Symbol | Type |
|-------|---------|--------|------|
| Ethereum | Mainnet | ETH | Native |
| USDT | Polygon | USDT | ERC-20 |
| USDC | Polygon | USDC | ERC-20 |
| DAI | Polygon | DAI | ERC-20 |
| Wrapped ETH | Polygon | WETH | ERC-20 |

### 2.2 Payment Contract

```solidity
// contracts/payments/CryptoPayment.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract CryptoPayment is AccessControl {
    using SafeERC20 for IERC20;
    
    struct Payment {
        bytes32 paymentId;
        address payer;
        address token;
        uint256 amount;
        uint256 amountInUSD;
        PaymentStatus status;
        uint256 timestamp;
    }
    
    enum PaymentStatus { PENDING, COMPLETED, REFUNDED, FAILED }
    
    mapping(bytes32 => Payment) public payments;
    mapping(address => uint256) public processedPayments;
    
    uint256 public usdAmount;
    address public treasury;
    IERC20 public usdt;
    
    event PaymentInitiated(
        bytes32 indexed paymentId,
        address indexed payer,
        address token,
        uint256 amount
    );
    event PaymentCompleted(
        bytes32 indexed paymentId,
        uint256 amountInUSD
    );
    event PaymentRefunded(
        bytes32 indexed paymentId
    );
    
    constructor(address _treasury, address _usdt) {
        treasury = _treasury;
        usdt = IERC20(_usdt);
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    function payWithToken(
        address token,
        uint256 amount,
        bytes32 paymentId
    ) external {
        require(payments[paymentId].timestamp == 0, "Payment exists");
        
        IERC20 tokenContract = IERC20(token);
        
        // Transfer tokens from user
        tokenContract.safeTransferFrom(msg.sender, address(this), amount);
        
        // Calculate USD value (using oracle)
        uint256 amountInUSD = getTokenValue(token, amount);
        
        // Record payment
        payments[paymentId] = Payment({
            paymentId: paymentId,
            payer: msg.sender,
            token: token,
            amount: amount,
            amountInUSD: amountInUSD,
            status: PaymentStatus.PENDING,
            timestamp: block.timestamp
        });
        
        emit PaymentInitiated(paymentId, msg.sender, token, amount);
        
        // Confirm payment
        confirmPayment(paymentId);
    }
    
    function payWithETH(bytes32 paymentId) external payable {
        require(payments[paymentId].timestamp == 0, "Payment exists");
        
        uint256 amountInUSD = getEthValue(msg.value);
        
        require(amountInUSD >= usdAmount, "Insufficient payment");
        
        // Record payment
        payments[paymentId] = Payment({
            paymentId: paymentId,
            payer: msg.sender,
            token: address(0),
            amount: msg.value,
            amountInUSD: amountInUSD,
            status: PaymentStatus.PENDING,
            timestamp: block.timestamp
        });
        
        emit PaymentInitiated(paymentId, msg.sender, address(0), msg.value);
        
        // Send excess back
        uint256 excess = msg.value - (amountInUSD / getEthPrice());
        if (excess > 0) {
            payable(msg.sender).transfer(excess);
        }
        
        confirmPayment(paymentId);
    }
    
    function confirmPayment(bytes32 paymentId) internal {
        Payment storage payment = payments[paymentId];
        
        payment.status = PaymentStatus.COMPLETED;
        
        if (payment.token != address(0)) {
            IERC20(payment.token).safeTransfer(
                treasury,
                payment.amount
            );
        } else {
            payable(treasury).transfer(payment.amount);
        }
        
        processedPayments[payment.payer] += payment.amountInUSD;
        
        emit PaymentCompleted(paymentId, payment.amountInUSD);
    }
    
    function getTokenValue(
        address token,
        uint256 amount
    ) public view returns (uint256) {
        // Use Chainlink or other oracle
    }
}
```

## 3. Payment Gateway

### 3.1 Backend Service

```typescript
import { ethers } from 'ethers';
import CryptoPaymentABI from './abis/CryptoPayment.json';

class CryptoPaymentService {
  private contract: ethers.Contract;
  
  constructor() {
    const provider = new ethers.providers.JsonRpcProvider(process.env.RPC_URL);
    const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);
    
    this.contract = new ethers.Contract(
      CRYPTO_PAYMENT_ADDRESS,
      CryptoPaymentABI,
      wallet
    );
  }
  
  async createPayment(
    userId: string,
    amountUSD: number,
    token: string
  ): Promise<PaymentRequest> {
    // Generate payment ID
    const paymentId = ethers.utils.keccak256(
      ethers.utils.toUtf8Bytes(`${userId}-${Date.now()}`)
    );
    
    // Get token price
    const tokenPrice = await this.getTokenPrice(token);
    const tokenAmount = (amountUSD / tokenPrice) * 1e8; // With decimals
    
    return {
      paymentId,
      amountUSD,
      token,
      tokenAmount: tokenAmount.toString(),
      tokenDecimal: await this.getTokenDecimals(token),
      contractAddress: CRYPTO_PAYMENT_ADDRESS,
    };
  }
  
  async checkPaymentStatus(paymentId: string): Promise<PaymentStatus> {
    const payment = await this.contract.payments(paymentId);
    
    if (payment.status === 2) { // COMPLETED
      return PaymentStatus.COMPLETED;
    } else if (payment.status === 3) { // REFUNDED
      return PaymentStatus.REFUNDED;
    } else if (payment.status === 4) { // FAILED
      return PaymentStatus.FAILED;
    }
    
    return PaymentStatus.PENDING;
  }
}
```

### 3.2 Frontend Integration

```typescript
import { useWeb3React } from '@web3-react/core';
import { ethers } from 'ethers';

const CryptoCheckout = ({ amount, onSuccess }) => {
  const { account, library } = useWeb3React();
  
  const handlePayment = async () => {
    const paymentRequest = await api.createPayment(amount, 'USDT');
    
    const token = new ethers.Contract(
      USDT_ADDRESS,
      ERC20_ABI,
      library.getSigner()
    );
    
    // Approve tokens
    const decimals = await token.decimals();
    await token.approve(
      CRYPTO_PAYMENT_ADDRESS,
      ethers.utils.parseUnits(paymentRequest.tokenAmount, decimals)
    );
    
    // Make payment
    const contract = new ethers.Contract(
      CRYPTO_PAYMENT_ADDRESS,
      CryptoPaymentABI,
      library.getSigner()
    );
    
    const tx = await contract.payWithToken(
      USDT_ADDRESS,
      ethers.utils.parseUnits(paymentRequest.tokenAmount, decimals),
      paymentRequest.paymentId
    );
    
    await tx.wait();
    
    // Verify payment
    const status = await api.checkPaymentStatus(paymentRequest.paymentId);
    
    if (status === 'COMPLETED') {
      onSuccess();
    }
  };
  
  return (
    <div>
      <h2>Pay with Crypto</h2>
      <p>Amount: ${amount}</p>
      <button onClick={handlePayment}>
        Pay with USDT
      </button>
    </div>
  );
};
```

## 4. Price Oracle

### 4.1 Chainlink Integration

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";

contract PriceOracle {
    mapping(address => address) public priceFeeds;
    
    constructor() {
        // Set price feeds for supported tokens
        priceFeeds[0x...USDT] = 0x...; // Chainlink feed
        priceFeeds[0x...USDC] = 0x...;
        priceFeeds[0x...DAI] = 0x...;
        priceFeeds[0x...ETH] = 0x...;
    }
    
    function getTokenPrice(address token) public view returns (int256) {
        AggregatorV3Interface feed = AggregatorV3Interface(
            priceFeeds[token]
        );
        
        (, int256 price, , , ) = feed.latestRoundData();
        
        return price;
    }
    
    function getTokenValue(
        address token,
        uint256 amount
    ) external view returns (uint256) {
        int256 price = getTokenPrice(token);
        
        return (amount * uint256(price)) / 1e8;
    }
}
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
