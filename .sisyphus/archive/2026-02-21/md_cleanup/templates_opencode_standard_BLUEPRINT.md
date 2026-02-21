# {{PROJECT_NAME}} - BLUEPRINT

> Comprehensive project blueprint following 22-pillar structure

---

## 1. Project Overview

### Purpose
{{PROJECT_DESCRIPTION}}

### Goals
- Provide a scalable, maintainable codebase
- Follow best practices 2026
- Enable rapid feature development

### Target Users
- Developers building {{USE_CASE}}
- DevOps teams deploying the application
- End users consuming the API

---

## 2. Architecture

### Tech Stack
- **Runtime:** Node.js 20+
- **Language:** TypeScript (Strict Mode)
- **Framework:** Express.js / Fastify
- **Database:** PostgreSQL 15+
- **ORM:** Prisma
- **Cache:** Redis 7+
- **Container:** Docker

### High-Level Architecture
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│   API GW    │────▶│  Service    │
│  (Browser)  │     │  (Express)  │     │  (Node.js)  │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                    ┌─────────────┐             │
                    │    Redis    │◀────────────┤
                    │   (Cache)   │             │
                    └─────────────┘             │
                                               │
                    ┌─────────────┐             │
                    │ PostgreSQL  │◀────────────┘
                    │   (Data)    │
                    └─────────────┘
```

---

## 3. Technology Decisions

### Why This Stack?
| Component | Choice | Rationale |
|-----------|--------|-----------|
| Runtime | Node.js 20 | LTS, fast, ecosystem |
| Language | TypeScript | Type safety, IDE support |
| Framework | Express | Maturity, flexibility |
| ORM | Prisma | Type safety, migrations |
| Database | PostgreSQL | Reliability, ACID |

---

## 4. Folder Structure

```
{{PROJECT_NAME}}/
├── src/
│   ├── config/
│   ├── controllers/
│   ├── services/
│   ├── repositories/
│   ├── models/
│   ├── middleware/
│   ├── utils/
│   └── app.ts
├── tests/
├── scripts/
├── docs/
│   ├── dev/
│   └── non-dev/
├── docker/
├── .env.example
├── package.json
├── tsconfig.json
└── docker-compose.yml
```

---

## 5. Database Schema

### Core Entities

```prisma
model User {
  id        String   @id @default(uuid())
  email     String   @unique
  name      String?
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
```

---

## 6. API Design

### RESTful Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/v1/users | List users |
| POST | /api/v1/users | Create user |
| GET | /api/v1/users/:id | Get user |
| PUT | /api/v1/users/:id | Update user |
| DELETE | /api/v1/users/:id | Delete user |

### Request/Response Format
- **Content-Type:** application/json
- **Authentication:** Bearer JWT

---

## 7. Configuration

### Environment Variables
```bash
# Required
DATABASE_URL=postgresql://user:pass@localhost:5432/db
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key

# Optional
PORT=50001
NODE_ENV=development
LOG_LEVEL=debug
```

---

## 8. Security

### Measures
- JWT authentication
- Input validation (Zod)
- SQL injection prevention (Prisma)
- CORS configuration
- Rate limiting
- Helmet headers

---

## 9. Performance

### Optimizations
- Redis caching for frequently accessed data
- Database indexing
- Pagination for list endpoints
- Connection pooling

---

## 10. Testing

### Strategy
- **Unit Tests:** Vitest
- **Integration Tests:** Supertest
- **Coverage Target:** 80%

---

## 11. CI/CD

### Pipeline
1. Lint & Type Check
2. Unit Tests
3. Build Docker Image
4. Deploy to Staging
5. Deploy to Production

---

## 12. Monitoring

### Metrics
- Response time (P95, P99)
- Error rate
- Request volume
- Database query time

### Tools
- Prometheus (metrics)
- Grafana (dashboards)

---

## 13. Deployment

### Docker
```bash
# Build
docker build -t {{PROJECT_NAME}}:latest .

# Run
docker-compose up -d
```

---

## 14. Error Handling

### Strategy
- Structured error responses (RFC 7807)
- Centralized error handler middleware
- Detailed logging
- User-friendly error messages

---

## 15. Logging

### Format
```json
{
  "timestamp": "2026-01-01T00:00:00Z",
  "level": "info",
  "message": "Request processed",
  "context": { "method": "GET", "path": "/api/v1/users" }
}
```

---

## 16. Dependencies

### Core Dependencies
- express
- @prisma/client
- ioredis
- jsonwebtoken
- zod

### Dev Dependencies
- typescript
- vitest
- eslint
- prettier
- prisma

---

## 17. Versioning

### Strategy
- API Versioning: URL-based (`/v1/`, `/v2/`)
- Semantic Versioning for releases

---

## 18. Documentation

### Resources
- README.md (this file)
- API Reference: `/docs/dev/api-reference.md`
- Deployment: `/docs/dev/deployment.md`

---

## 19. Contributing

### Workflow
1. Fork repository
2. Create feature branch
3. Write tests
4. Submit PR
5. Code review
6. Merge

---

## 20. Roadmap

### Phase 1 (Current)
- [x setup
- [x] Basic CRUD] Project
- [ ] Authentication
- [ ] Caching

### Phase 2
- [ ] Advanced filtering
- [ ] Rate limiting
- [ ] Webhooks

---

## 21. Glossary

| Term | Definition |
|------|------------|
| API | Application Programming Interface |
| JWT | JSON Web Token |
| ORM | Object-Relational Mapping |
| REST | Representational State Transfer |

---

## 22. References

- [TypeScript Strict Mode](https://www.typescriptlang.org/tsconfig#strict)
- [Prisma Documentation](https://www.prisma.io/docs/)
- [Express.js](https://expressjs.com/)
- [Best Practices 2026](/Users/jeremy/.config/opencode/AGENTS.md)

---

**Document Version:** 1.0.0  
**Last Updated:** {{DATE}}  
**Status:** Active
