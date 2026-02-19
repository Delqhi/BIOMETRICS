# Integration Workflow Template

## Overview

The Integration workflow template provides a comprehensive solution for implementing third-party service integrations. This template guides the process of integrating external APIs, libraries, and services into your application, ensuring proper implementation, testing, and documentation.

The workflow covers the full integration lifecycle: from initial assessment and design through implementation, testing, and documentation. It leverages AI agents to understand integration requirements, generate appropriate code, and create comprehensive tests.

This template is essential for organizations seeking to:
- Quickly integrate external services
- Ensure integration best practices
- Maintain code quality during integration
- Document integrated services
- Test integrations thoroughly

## Purpose

The primary purpose of the Integration template is to:

1. **Streamline Integration** - Provide a structured approach to integrations
2. **Ensure Best Practices** - Follow industry standards for integrations
3. **Automate Implementation** - Generate integration code automatically
4. **Comprehensive Testing** - Validate integration behavior
5. **Maintain Documentation** - Keep integration docs up to date

### Key Use Cases

- **API Integration** - Connect to external REST/GraphQL APIs
- **SDK Integration** - Integrate third-party SDKs
- **Service Integration** - Connect to microservices
- **Library Integration** - Add external libraries
- **Platform Integration** - Integrate with platforms (AWS, Azure, etc.)

## Input Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `integration_type` | string | Yes | - | Type: api, sdk, service, library |
| `service_name` | string | Yes | - | Name of service to integrate |
| `service_url` | string | No | - | Service URL or documentation link |
| `auth_type` | string | No | - | Authentication type |
| `features` | array | No | all | Specific features to implement |

### Input Examples

```yaml
# Example 1: API Integration
inputs:
  integration_type: api
  service_name: Stripe
  service_url: https://stripe.com/docs/api
  auth_type: bearer
  features:
    - payments
    - refunds
    - customers

# Example 2: SDK Integration
inputs:
  integration_type: sdk
  service_name: AWS-S3
  features:
    - upload
    - download
    - list_buckets
```

## Output Results

```json
{
  "integration": {
    "name": "stripe",
    "type": "api",
    "status": "completed"
  },
  "files": {
    "client": "src/integrations/stripe/client.ts",
    "tests": "src/integrations/stripe/client.test.ts",
    "types": "src/integrations/stripe/types.ts"
  },
  "coverage": 85
}
```

## Workflow Steps

### Step 1: Analyze Requirements

Analyzes integration requirements from documentation.

### Step 2: Create Client

Generates API client code.

### Step 3: Add Authentication

Implements authentication logic.

### Step 4: Generate Tests

Creates comprehensive integration tests.

### Step 5: Document

Generates integration documentation.

## Usage

```bash
biometrics workflow run integration \
  --integration_type api \
  --service_name Stripe \
  --service_url https://stripe.com/docs/api
```

## Troubleshooting

- **Auth Issues**: Verify credentials are correctly configured
- **Connection Errors**: Check network connectivity and firewall rules

## Related Templates

- **Config Validator** (`config-validator/`) - Validate integration configs
- **Code Review** (`code-review/`) - Review integration code

---

**Template Version:** 1.0.0  
**Author:** BIOMETRICS Team  
**Category:** Integration  

*Last Updated: February 2026*
