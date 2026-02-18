# NFT-MARKETPLACE.md - NFT Marketplace

**Version:** 1.0  
**Date:** 2026-02-18  
**Status:** Phase 3 - Advanced Tech (Blockchain)  
**Author:** BIOMETRICS Blockchain Team

---

## 1. Overview

This document describes the NFT marketplace for BIOMETRICS health achievements and data tokens, enabling trading and monetization of verified health records.

## 2. Marketplace Architecture

### 2.1 Smart Contracts

| Contract | Purpose |
|----------|---------|
| HealthNFT | Health achievement tokens |
| Marketplace | Trading marketplace |
| Royalties | Fee distribution |

### 2.2 Marketplace Contract

```solidity
// contracts/marketplace/Marketplace.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract Marketplace is ReentrancyGuard, AccessControl {
    using SafeERC20 for IERC20;
    
    struct Listing {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 price;
        address paymentToken;
        bool active;
    }
    
    mapping(bytes32 => Listing) public listings;
    mapping(address => uint256) public earnings;
    
    uint256 public constant ROYALTY_FEE = 250; // 2.5%
    uint256 public constant MARKETPLACE_FEE = 250; // 2.5%
    
    event Listed(
        bytes32 listingId,
        address indexed seller,
        uint256 tokenId,
        uint256 price
    );
    event Sold(
        bytes32 listingId,
        address indexed buyer,
        uint256 price
    );
    event Cancelled(bytes32 listingId);
    
    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }
    
    function list(
        address nftContract,
        uint256 tokenId,
        uint256 price,
        address paymentToken
    ) external nonReentrant {
        require(price > 0, "Price must be > 0");
        
        IERC721 nft = IERC721(nftContract);
        require(
            nft.ownerOf(tokenId) == msg.sender,
            "Not the owner"
        );
        require(
            nft.getApproved(tokenId) == address(this) ||
            nft.isApprovedForAll(msg.sender, address(this)),
            "Not approved"
        );
        
        bytes32 listingId = keccak256(
            abi.encodePacked(nftContract, tokenId, msg.sender, block.timestamp)
        );
        
        listings[listingId] = Listing({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            price: price,
            paymentToken: paymentToken,
            active: true
        });
        
        emit Listed(listingId, msg.sender, tokenId, price);
    }
    
    function buy(bytes32 listingId) external nonReentrant {
        Listing storage listing = listings[listingId];
        require(listing.active, "Not active");
        
        listing.active = false;
        
        IERC20 paymentToken = IERC20(listing.paymentToken);
        
        // Calculate fees
        uint256 price = listing.price;
        uint256 royaltyFee = (price * ROYALTY_FEE) / 10000;
        uint256 marketplaceFee = (price * MARKETPLACE_FEE) / 10000;
        uint256 sellerAmount = price - royaltyFee - marketplaceFee;
        
        // Transfer payment
        paymentToken.transferFrom(msg.sender, address(this), price);
        
        // Pay marketplace
        earnings[address(this)] += marketplaceFee;
        
        // Pay royalty to original minter
        IERC721 nft = IERC721(listing.nftContract);
        address originalMinter = nft.ownerOf(listing.tokenId);
        paymentToken.transfer(originalMinter, royaltyFee);
        
        // Pay seller
        paymentToken.transfer(listing.seller, sellerAmount);
        
        // Transfer NFT
        nft.safeTransferFrom(
            listing.seller,
            msg.sender,
            listing.tokenId
        );
        
        emit Sold(listingId, msg.sender, price);
    }
    
    function cancel(bytes32 listingId) external nonReentrant {
        Listing storage listing = listings[listingId];
        require(listing.seller == msg.sender, "Not seller");
        require(listing.active, "Not active");
        
        listing.active = false;
        
        emit Cancelled(listingId);
    }
}
```

## 3. NFT Types

### 3.1 Achievement NFTs

```solidity
// Achievement types
enum AchievementType {
    STREAK_7_DAYS,
    STREAK_30_DAYS,
    STREAK_100_DAYS,
    GOAL_REACHED,
    PERSONAL_BEST,
    CHALLENGE_COMPLETED,
    VERIFIED_DATA,
    COMMUNITY_HELP
}

struct Achievement {
    AchievementType achievementType;
    uint256 timestamp;
    string metadata;
    uint256 score;
}
```

### 3.2 Data NFTs

```typescript
interface HealthDataNFT {
  // Medical data ownership
  struct HealthRecord {
    string ipfsHash;
    uint256 timestamp;
    address verifier;
    bytes signature;
  }
  
  // Mint data as NFT
  function mintHealthData(
    address to,
    string calldata tokenURI,
    string calldata dataHash
  ) external returns (uint256);
  
  // Verify data authenticity
  function verifyData(uint256 tokenId, string calldata dataHash) external view returns (bool);
}
```

## 4. Frontend Integration

### 4.1 Marketplace UI

```typescript
import { useWeb3React } from '@web3-react/core';

const Marketplace = () => {
  const { account, library } = useWeb3React();
  const [listings, setListings] = useState([]);
  
  useEffect(() => {
    loadListings();
  }, []);
  
  const loadListings = async () => {
    const marketplace = new ethers.Contract(
      MARKETPLACE_ADDRESS,
      MARKETPLACE_ABI,
      library.getSigner()
    );
    
    // Load all active listings
    // ... implement with events
  };
  
  const buyListing = async (listingId: string) => {
    const marketplace = new ethers.Contract(
      MARKETPLACE_ADDRESS,
      MARKETPLACE_ABI,
      library.getSigner()
    );
    
    const listing = await marketplace.listings(listingId);
    
    const tx = await marketplace.buy(listingId);
    await tx.wait();
    
    // Show success
  };
  
  return (
    <div className="marketplace">
      <h1>Health NFT Marketplace</h1>
      
      <div className="listings-grid">
        {listings.map(listing => (
          <NFTCard
            key={listing.id}
            listing={listing}
            onBuy={buyListing}
          />
        ))}
      </div>
    </div>
  );
};
```

## 5. Trading Features

### 5.1 Features

| Feature | Description |
|---------|-------------|
| Fixed Price | List at fixed price |
| Auction | Time-based auction |
| Offer | Accept offers |
| Bundle | Bundle multiple NFTs |

### 5.2 Auction Contract

```solidity
contract Auction {
    struct AuctionData {
        address seller;
        address nftContract;
        uint256 tokenId;
        uint256 startingPrice;
        uint256 highestBid;
        address highestBidder;
        uint256 endTime;
        bool active;
    }
    
    mapping(bytes32 => AuctionData) public auctions;
    
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startingPrice,
        uint256 duration
    ) external {
        bytes32 auctionId = keccak256(
            abi.encodePacked(nftContract, tokenId, msg.sender, block.timestamp)
        );
        
        auctions[auctionId] = AuctionData({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startingPrice: startingPrice,
            highestBid: startingPrice,
            highestBidder: address(0),
            endTime: block.timestamp + duration,
            active: true
        });
        
        IERC721(nftContract).transferFrom(
            msg.sender,
            address(this),
            tokenId
        );
    }
    
    function bid(bytes32 auctionId, uint256 bidAmount) external {
        AuctionData storage auction = auctions[auctionId];
        require(auction.active, "Auction not active");
        require(block.timestamp < auction.endTime, "Auction ended");
        require(bidAmount > auction.highestBid, "Bid too low");
        
        // Refund previous highest bidder
        if (auction.highestBidder != address(0)) {
            // Refund logic
        }
        
        auction.highestBid = bidAmount;
        auction.highestBidder = msg.sender;
    }
}
```

---

**Last Updated:** 2026-02-18  
**Next Review:** 2026-03-18  
**Version:** 1.0
