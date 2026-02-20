# API Documentation Template

## Overview

The `api-docs` template generates comprehensive API documentation from source code annotations. This template creates clean, readable API documentation with examples, parameter descriptions, and response schemas.

## Features

- **OpenAPI/Swagger Support**: Generate OpenAPI 3.0 specifications
- **Code Annotations**: Parse docstrings and annotations
- **Interactive Console**: Try API calls directly
- **Code Examples**: Auto-generate client examples
- **Version Support**: Document multiple API versions

## Usage

### Generate Documentation
```bash
biometrics-cli docs generate api-docs \
  --source ./src/api \
  --output ./docs/api
```

### Watch Mode
```bash
biometrics-cli docs generate api-docs \
  --source ./src/api \
  --watch
```

## Configuration

```yaml
api-docs:
  source_dir: "./src/api"
  output_dir: "./docs/api"
  
  openapi:
    version: "3.0.3"
    title: "Biometrics API"
    description: "Biometric Authentication API"
    version: "1.0.0"
  
  sections:
    - Authentication
    - Users
    - Biometrics
    - Analytics
  
  code_examples:
    languages:
      - curl
      - python
      - javascript
      - go
```

## Annotation Syntax

### Endpoint Documentation
```python
"""
Authentication API

Endpoints for user authentication and authorization.

## Endpoints

### POST /api/v1/auth/login
Authenticate user and receive access token.

**Parameters:**
| Name | Type | In | Description |
|------|------|-----|-------------|
| email | string | body | User email |
| password | string | body | User password |

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

**Example:**
```bash
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secret"}'
```
"""
```

### Schema Documentation
```python
class UserSchema:
    """
    User data schema.
    
    Attributes:
        id: Unique user identifier
        email: User email address
        name: Full name
        created_at: Account creation timestamp
    """
    id: str
    email: str
    name: str
    created_at: datetime
```

## Output Structure

```
docs/api/
├── index.html              # Main documentation
├── openapi.json           # OpenAPI specification
├── endpoints/
│   ├── authentication.md
│   ├── users.md
│   └── biometrics.md
├── schemas/
│   ├── UserSchema.md
│   └── AuthResponse.md
└── examples/
    ├── curl/
    ├── python/
    ├── javascript/
    └── go/
```

## Interactive Documentation

Access interactive API console at:
```
https://api.example.com/docs
```

Features:
- Try API calls
- View real responses
- Generate code samples
- Download OpenAPI spec

## OpenAPI Specification

### Basic Info
```yaml
openapi: 3.0.3
info:
  title: Biometrics API
  version: 1.0.0
  description: Biometric Authentication API
  
servers:
  - url: https://api.example.com
    description: Production
  - url: https://staging-api.example.com
    description: Staging
```

### Paths
```yaml
paths:
  /api/v1/auth/login:
    post:
      summary: User Login
      operationId: login
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Successful login
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'
```

### Components
```yaml
components:
  schemas:
    LoginRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
        password:
          type: string
          format: password
          
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

## Versioning

### Version Headers
```bash
# Request specific version
Accept: application/vnd.biometrics.v1+json
```

### Versioned Documentation
```bash
# Generate docs for specific version
biometrics-cli docs generate api-docs --version v1
biometrics-cli docs generate api-docs --version v2
```

## Customization

### Custom Templates
Create custom templates in `templates/api-docs/`:

```html
<!-- Custom endpoint template -->
<div class="endpoint">
  <h2>{{.Title}}</h2>
  <p>{{.Description}}</p>
  {{range .Methods}}
    <div class="method {{.Verb}}">{{.Verb}} {{.Path}}</div>
  {{end}}
</div>
```

## Integration

### CI/CD
```yaml
- name: Generate API Docs
  run: biometrics-cli docs generate api-docs
  artifacts:
    paths:
      - docs/api/
```

### GitHub Pages
```yaml
- name: Deploy to Pages
  uses: peaceiris/actions-gh-pages@v3
  with:
    publish_dir: docs/api
```

## Performance

| Metric | Value |
|--------|-------|
| Generation Speed | 1000 endpoints/minute |
| HTML Size | ~50KB per endpoint |
| Load Time | < 1 second |

## Troubleshooting

### Missing Documentation
1. Check annotations are in correct format
2. Verify source files are included
3. Check file path configuration

### Generation Errors
1. Validate OpenAPI schema
2. Check for circular references
3. Review error logs

## See Also

- [CLI Commands](../cmd/README.md)
- [Documentation](./README.md)
- [Configuration](../docs/configuration.md)
