# 🔌 BIOMETRICS Documentation - API Reference

**Detailed API documentation, endpoints, and mappings.**

---

## 📁 API Documents

| File | Description |
|------|-------------|
| [MAPPING.md](MAPPING.md) | System mapping overview |
| [MAPPING-COMMANDS-ENDPOINTS.md](MAPPING-COMMANDS-ENDPOINTS.md) | Commands to endpoints |
| [MAPPING-DB-API.md](MAPPING-DB-API.md) | Database to API mapping |
| [MAPPING-FRONTEND-BACKEND.md](MAPPING-FRONTEND-BACKEND.md) | Frontend-backend mapping |
| [MAPPING-NLM-ASSETS.md](MAPPING-NLM-ASSETS.md) | NLM assets mapping |

---

## 🗺️ System Mappings

### Commands → Endpoints
Maps CLI commands to API endpoints for complete traceability.

### Database → API
Documents how database tables map to API resources.

### Frontend → Backend
Tracks frontend components to backend services.

### NLM Assets
Maps NotebookLM assets to system components.

---

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - User logout

### Users
- `GET /api/v1/users` - List users
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Products
- `GET /api/v1/products` - List products
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

### Orders
- `GET /api/v1/orders` - List orders
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/:id` - Get order details

---

## 🔍 API Documentation Standards

All APIs must have:
1. **OpenAPI/Swagger spec** - Machine-readable documentation
2. **Request/Response examples** - Clear usage examples
3. **Error codes** - Comprehensive error handling
4. **Rate limits** - Documented throttling
5. **Authentication** - Security requirements

---

**Last Updated:** 2026-02-18  
**Status:** ✅ Production-Ready
