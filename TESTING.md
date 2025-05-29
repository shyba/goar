# Testing Guide

This document provides information about testing the Goar library.

## Test Structure

The project includes comprehensive unit tests for all core functionality:

### Unit Tests (No Network Required)

- **`crypto/`** - Cryptographic functions (SHA256, base64url, deep hash)
- **`tag/`** - Tag encoding/decoding with Apache Avro format
- **`transaction/`** - Transaction creation, signing, verification, and Merkle trees
- **`signer/`** - Key management and wallet signing operations
- **`uploader/`** - Transaction upload logic (structure and validation only)
- **`transaction/bundle/`** - ANS-104 bundle functionality
- **`transaction/data_item/`** - ANS-104 data item functionality

### Integration Tests (Network Required)

- **`client/`** - HTTP API client (requires running Arweave node)
- **`wallet/`** - High-level wallet operations (requires running Arweave node)

## Running Tests

### Unit Tests Only
```bash
# Run all unit tests (no network required)
go test ./crypto ./tag ./transaction ./signer ./uploader ./transaction/bundle ./transaction/data_item -v

# Run tests in short mode (skips slow tests)
go test ./... -short
```

### All Tests (Including Integration)
```bash
# Run all tests (requires Arweave test node at localhost:1984)
go test ./... -v
```

### Individual Package Tests
```bash
# Test specific packages
go test ./transaction -v
go test ./crypto -v
go test ./signer -v
```

## Test Requirements

### For Unit Tests
- Go 1.18+ 
- Test files in `test/` directory (included in repository)

### For Integration Tests
- Local Arweave development node running on `localhost:1984`
- Test wallet file at `test/signer.json` (included)

### Setting Up Arweave Test Node

To run integration tests, you need a local Arweave node:

```bash
# Using Docker
docker run -p 1984:1984 arweave/arweave-testnet

# Or download and run Arweave node locally
# See: https://docs.arweave.org/developers/mining/mining-guide
```

## Test Coverage

### Crypto Package
- ✅ SHA256 hashing
- ✅ Base64URL encoding/decoding  
- ✅ Deep hash computation
- ✅ RSA-PSS signing and verification

### Transaction Package
- ✅ Transaction creation with data and tags
- ✅ Transaction signing and verification
- ✅ Merkle tree generation and validation
- ✅ Chunk preparation and retrieval
- ✅ Data integrity validation

### Signer Package
- ✅ Key generation (4096-bit RSA)
- ✅ JWK file loading and parsing
- ✅ Wallet address computation
- ✅ Private key management

### Tag Package
- ✅ Avro encoding/decoding
- ✅ Base64URL tag conversion
- ✅ Tag serialization/deserialization

### Uploader Package
- ✅ Uploader initialization
- ✅ Fatal error detection
- ✅ Configuration validation
- ✅ Transaction preparation

### Bundle & Data Item Packages
- ✅ ANS-104 data item creation
- ✅ Bundle generation
- ✅ Binary encoding/decoding
- ✅ Signature verification

## Test Data

The repository includes test data files:

- `test/signer.json` - Test wallet JWK file
- `test/1MB.bin` - 1MB test data for chunking tests
- `test/lotsofdata.bin` - Large test data file
- `test/rebar3` - Binary file with known Merkle root
- `test/1115BDataItem` - Pre-encoded ANS-104 data item

## Network-Dependent Tests

Some tests require a running Arweave node and are skipped in unit test runs:

### Client Tests
- Transaction submission and retrieval
- Wallet balance queries
- Network information retrieval
- Block data fetching

### Wallet Tests  
- End-to-end transaction workflows
- Network fee calculation
- Transaction anchoring

These tests will fail if no Arweave node is available but don't affect the core library functionality.

## Test Documentation Standards

All test functions include:
- Clear descriptive names
- Comprehensive test cases covering edge cases
- Proper error handling verification
- Documentation comments explaining test purpose
- Use of `require` for critical assertions
- Use of `assert` for validation assertions

## CI/CD Considerations

For continuous integration:

```bash
# Run only unit tests (recommended for CI)
go test ./crypto ./tag ./transaction ./signer ./uploader ./transaction/bundle ./transaction/data_item

# Run with coverage
go test -coverprofile=coverage.out ./crypto ./tag ./transaction ./signer ./uploader ./transaction/bundle ./transaction/data_item
go tool cover -html=coverage.out
```

Integration tests should be run in a separate CI stage with proper Arweave node setup. 