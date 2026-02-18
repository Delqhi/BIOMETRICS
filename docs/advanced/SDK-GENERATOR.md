# SDK-GENERATOR.md

**Project:** BIOMETRICS  
**Version:** 1.0.0  
**Last Updated:** 2026-02-18  
**Status:** Active

---

## Overview

The SDK Generator creates and maintains official client libraries for multiple programming languages. This document covers supported languages, generation process, and versioning.

## Supported Languages

| Language | Status | Package Manager |
|----------|--------|----------------|
| JavaScript/TypeScript | Stable | npm |
| Python | Stable | pip |
| Go | Stable | go get |
| Java | Beta | Maven |
| Ruby | Beta | gem |
| PHP | Beta | composer |
| Rust | Alpha | crates.io |
| C#/.NET | Alpha | NuGet |

## Generation Process

### OpenAPI Specification

```yaml
openapi: 3.0.0
info:
  title: BIOMETRICS API
  version: 1.0.0
  description: Biometric authentication API
servers:
  - url: https://api.biometrics.com/v1
paths:
  /users:
    get:
      summary: List users
      responses:
        '200':
          description: Success
```

### Code Generation

```bash
# Generate TypeScript SDK
openapi-generator generate \
  -i openapi.yaml \
  -g typescript \
  -o ./sdks/typescript

# Generate Python SDK
openapi-generator generate \
  -i openapi.yaml \
  -g python \
  -o ./sdks/python

# Generate Go SDK
openapi-generator generate \
  -i openapi.yaml \
  -g go \
  -o ./sdks/go
```

## SDK Structure

### TypeScript SDK

```
biometrics-sdk/
├── src/
│   ├── index.ts
│   ├── client.ts
│   ├── types/
│   │   ├── user.ts
│   │   └── scan.ts
│   ├── resources/
│   │   ├── users.ts
│   │   └── scans.ts
│   └── utils/
│       └── http.ts
├── package.json
└── tsconfig.json
```

### Usage Example

```typescript
import { BiometricsClient } from '@biometrics/sdk';

const client = new BiometricsClient({
  apiKey: process.env.BIOMETRICS_API_KEY
});

// List users
const users = await client.users.list({
  limit: 10
});

// Create scan
const scan = await client.scans.create({
  template: biometricTemplate,
  deviceId: 'device-123'
});
```

## Versioning

### Semantic Versioning

| Version | Type | Example |
|---------|------|---------|
| Major | Breaking changes | 1.0.0 → 2.0.0 |
| Minor | New features | 1.0.0 → 1.1.0 |
| Patch | Bug fixes | 1.0.0 → 1.0.1 |

### Changelog

```markdown
## 2.0.0 (2026-02-18)

### Breaking Changes
- Removed deprecated `/auth/login` endpoint
- Changed response format for `/scans`

### New Features
- Added biometric template v2
- New batch processing API

### Bug Fixes
- Fixed race condition in scan creation
```

---

**Document Status:** Complete  
**Next Review:** 2026-03-18
