# DEFI.md - Decentralized Finance Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the DeFi integration for BIOMETRICS, enabling decentralized finance features including yield generation, lending, and liquidity provision.

## 2. DeFi Protocols

### 2.1 Integration Architecture

```
┌─────────────────────────────────────────────┐
│            BIOMETRICS DeFi Layer            │
├─────────────────────────────────────────────┤
│  Yield    │  Lending  │  Swaps  │ Staking │
├───────────┼───────────┼─────────┼─────────┤
│  Aave     │  Compound │ Uniswap │ Lido    │
│  Yearn    │  Euler    │ Curve   │ Rocket  │
└───────────┴───────────┴─────────┴─────────┘
            │           │         │
            └───────────┴─────────┘
                        │
              Blockchain Network (Polygon)
```

### 2.2 Supported Protocols

| Protocol | Purpose | Network |
|----------|---------|---------|
| Aave | Lending/Borrowing | Polygon |
| Uniswap | Token Swaps | Polygon |
| Curve | Stablecoin Swaps | Polygon |
| Lido | ETH Staking | Ethereum |
| Yearn | Yield Optimization | Polygon |

## 3. Lending Integration

### 3.1 Aave V3

```typescript
import { ethers } from 'ethers';
import AaveV3PoolABI from './abis/AaveV3Pool.json';
import { AaveProtocolDataProvider } from './typechain/AaveProtocolDataProvider';

class LendingService {
  private pool: ethers.Contract;
  private dataProvider: AaveProtocolDataProvider;
  
  constructor() {
    const provider = new ethers.providers.JsonRpcProvider(POLYGON_RPC);
    
    this.pool = new ethers.Contract(
      '0x794a61358D6845594F94dc1DB02A252b5b4814aD', // Aave V3 Pool
      AaveV3PoolABI,
      provider
    );
  }
  
  async supply(asset: string, amount: ethers.BigNumber) {
    const signer = /* get signer */;
    
    // Approve token
    const token = new ethers.Contract(asset, ERC20ABI, signer);
    await token.approve(POOL_ADDRESS, amount);
    
    // Supply
    const tx = await this.pool.supply(asset, amount, signer.address, 0);
    await tx.wait();
    
    return tx.hash;
  }
  
  async borrow(asset: string, amount: ethers.BigNumber) {
    const signer = /* get signer */;
    
    const tx = await this.pool.borrow(asset, amount, 0, 0, signer.address);
    await tx.wait();
    
    return tx.hash;
  }
  
  async getUserAccountData(user: string) {
    return await this.pool.getUserAccountData(user);
  }
  
  async getUserReserveData(asset: string, user: string) {
    return await this.pool.getUserReserveData(asset, user);
  }
}
```

### 3.2 Supply Strategy

```typescript
// Auto-supply idle tokens to Aave
class AutoSupplyStrategy {
  private lendingService: LendingService;
  private threshold: ethers.BigNumber;
  
  constructor() {
    this.lendingService = new LendingService();
    this.threshold = ethers.utils.parseEther('100');
  }
  
  async checkAndSupply(user: string, token: string) {
    const balance = await this.getTokenBalance(user, token);
    
    if (balance.gt(this.threshold)) {
      const supplyAmount = balance.sub(this.threshold);
      await this.lendingService.supply(token, supplyAmount);
    }
  }
}
```

## 4. Token Swaps

### 4.1 Uniswap V3

```typescript
import { ethers } from 'ethers';
import { SwapRouter } from '@uniswap/v3-sdk';
import { Token, CurrencyAmount, TradeType, Percent } from '@uniswap/sdk-core';

class SwapService {
  private router: ethers.Contract;
  private quoter: ethers.Contract;
  
  constructor() {
    this.router = new ethers.Contract(
      '0xE592427A0AEce92De3Edee1F18E0157C05861564', // SwapRouter
      SwapRouterABI,
      provider
    );
    
    this.quoter = new ethers.Contract(
      '0xb27308f9F90D607463bb33eA1BeBb20C11cf7dd4', // QuoterV2
      QuoterV2ABI,
      provider
    );
  }
  
  async getQuote(
    tokenIn: string,
    tokenOut: string,
    amountIn: ethers.BigNumber
  ) {
    const [quote] = await this.quoter.quoteExactInputSingle([
      tokenIn,
      tokenOut,
      amountIn,
      3000, // fee tier
      0 // sqrtPriceLimitX96
    ]);
    
    return quote;
  }
  
  async swap(
    tokenIn: string,
    tokenOut: string,
    amountIn: ethers.BigNumber,
    minAmountOut: ethers.BigNumber
  ) {
    const signer = /* get signer */;
    
    // Approve
    const token = new ethers.Contract(tokenIn, ERC20ABI, signer);
    await token.approve(ROUTER_ADDRESS, amountIn);
    
    // Swap
    const params = {
      tokenIn,
      tokenOut,
      fee: 3000,
      recipient: signer.address,
      deadline: Math.floor(Date.now() / 1000) + 60 * 10,
      amountIn,
      minAmountOut,
      sqrtPriceLimitX96: 0
    };
    
    const tx = await this.router.exactInputSingle(params, {
      gasLimit: 300000
    });
    
    return tx.wait();
  }
}
```

## 5. Yield Optimization

### 5.1 Yearn Strategy

```typescript
// Yearn Vault interaction
class YearnStrategy {
  private vault: ethers.Contract;
  
  constructor(vaultAddress: string) {
    this.vault = new ethers.Contract(
      vaultAddress,
      YearnVaultABI,
      provider
    );
  }
  
  async deposit(amount: ethers.BigNumber) {
    const signer = /* get signer */;
    
    // Approve
    const token = new ethers.Contract(
      await this.vault.token(),
      ERC20ABI,
      signer
    );
    await token.approve(this.vault.address, amount);
    
    // Deposit
    const tx = await this.vault.deposit(amount, signer.address);
    await tx.wait();
  }
  
  async withdraw(shares: ethers.BigNumber) {
    const signer = /* get signer */;
    
    const tx = await this.vault.withdraw(shares, signer.address, 1);
    await tx.wait();
  }
  
  async getShareValue(shares: ethers.BigNumber): Promise<ethers.BigNumber> {
    return await this.vault.convertToAssets(shares);
  }
}
```

### 5.2 Auto-Compounding

```typescript
class AutoCompoundStrategy {
  private yearnVault: YearnStrategy;
  private harvestThreshold: ethers.BigNumber;
  
  async harvest() {
    const pendingReward = await this.getPendingRewards();
    
    if (pendingReward.gt(this.harvestThreshold)) {
      // Claim rewards
      // Swap to deposit token
      // Re-deposit
    }
  }
}
```

## 6. Dashboard

### 6.1 Portfolio View

```typescript
const PortfolioDashboard = () => {
  const [positions, setPositions] = useState([]);
  
  useEffect(() => {
    loadPositions();
  }, []);
  
  const loadPositions = async () => {
    const lending = await getLendingPositions();
    const staking = await getStakingPositions();
    const farming = await getFarmingPositions();
    
    setPositions([...lending, ...staking, ...farming]);
  };
  
  return (
    <div>
      <h1>DeFi Portfolio</h1>
      <TotalValue positions={positions} />
      <PositionList positions={positions} />
    </div>
  );
};
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
