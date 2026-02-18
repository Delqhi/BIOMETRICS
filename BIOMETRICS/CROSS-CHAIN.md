# CROSS-CHAIN.md - Cross-Chain Bridge Integration

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the cross-chain bridge integration for BIOMETRICS, enabling asset transfers across multiple blockchain networks.

## 2. Supported Networks

### 2.1 Bridge Networks

| Network | Bridge | Native Token |
|---------|--------|-------------|
| Ethereum | Native Bridge | ETH |
| Polygon | Bridge | MATIC |
| Arbitrum | Bridge | ETH |
| Optimism | Bridge | ETH |
| BSC | Bridge | BNB |

### 2.2 Bridge Architecture

```typescript
// Cross-chain transfer
class CrossChainBridge {
  async bridge(
    fromChain: Chain,
    toChain: Chain,
    token: string,
    amount: bigint
  ) {
    const bridge = this.getBridge(fromChain, toChain);
    
    // Approve tokens
    const tokenContract = new ethers.Contract(token, ERC20ABI, signer);
    await tokenContract.approve(bridge.address, amount);
    
    // Initiate bridge
    const tx = await bridge.send(
      token,
      amount,
      signer.address,
      destinationChains[toChain]
    );
    
    return tx.wait();
  }
}
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
