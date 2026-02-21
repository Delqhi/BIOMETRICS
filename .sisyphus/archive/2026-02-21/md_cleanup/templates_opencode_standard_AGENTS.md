# {{PROJECT_NAME}} - AGENTS.md

## Projekt-Übersicht

- **Tech Stack:** TypeScript, Node.js, Express, PostgreSQL, Redis
- **Architektur:** REST API mit Modularer Struktur
- **Datenbank:** PostgreSQL mit Prisma ORM
- **Cache:** Redis für Session/Response Caching

## Konventionen

### Naming
- **Dateien:** kebab-case (`user-service.ts`)
- **Funktionen:** camelCase (`getUserById`)
- **Klassen:** PascalCase (`UserService`)
- ** Konstanten:** SCREAMING_SNAKE_CASE (`MAX_RETRY_COUNT`)

### Folder Structure
```
src/
├── controllers/    # HTTP Request Handler
├── services/       # Business Logic
├── repositories/   # Data Access Layer
├── models/        # TypeScript Interfaces
├── middleware/    # Express Middleware
├── utils/         # Utility Functions
├── config/        # Configuration
└── app.ts         # Main Application
```

### API-Standards
- **Base URL:** `http://localhost:50001/api`
- **Version:** `v1` (URL: `/api/v1/`)
- **Auth:** Bearer Token (JWT)
- **Response Format:** JSON
- **Errors:** RFC 7807 Problem Details

## API-Standards

### Request/Response
```typescript
// Request
interface ApiRequest<T> {
  headers: { Authorization: string }
  body: T
  params: { id: string }
  query: { page?: number; limit?: number }
}

// Response
interface ApiResponse<T> {
  data: T
  meta?: { page: number; limit: number; total: number }
}
```

### Error Handling
```typescript
// Standard Error Response
{
  "type": "https://api.example.com/errors/not-found",
  "title": "Not Found",
  "status": 404,
  "detail": "User with ID 123 not found"
}
```

## Testing
- **Framework:** Vitest
- **Coverage Target:** 80%
- **Test Location:** `src/**/*.test.ts`

## Troubleshooting

### Bekannte Probleme
1. **Port Conflict:** Port 50001 bereits belegt → `.env` PORT ändern
2. **Database Connection:** Prisma migrate ausführen `npm run db:migrate`

---

**Letzte Änderung:** {{DATE}}
