# Docker Hub Registry Documentation

**Project:** BIOMETRICS  
**Last Updated:** 2026-02-18  
**Maintainer:** DevOps Team

---

## Overview

This document describes the Docker Hub container registry setup for the BIOMETRICS project, including image management, automated builds, and deployment workflows.

## Registry Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      DOCKER REGISTRY ARCHITECTURE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐             │
│  │   GitHub     │────▶│ Docker Hub   │────▶│  Production  │             │
│  │   (Source)   │     │   (Build)    │     │   (Deploy)   │             │
│  └──────────────┘     └──────────────┘     └──────────────┘             │
│         │                     │                     │                      │
│         ▼                     ▼                     ▼                      │
│   Auto-builds            Image Tags           Pull & Run                    │
│   on push               & Versions           in Docker                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Docker Hub Configuration

### Repository Setup

**Repository Name:** `biometrics/app`  
**Visibility:** Private (requires Pro plan for private repos)  
**Organization:** `biometrics-org`

### Image Tags Strategy

| Tag | Description | When to Use |
|-----|-------------|-------------|
| `latest` | Most recent build | Development |
| `stable` | Production-ready | Production deployments |
| `X.Y.Z` | Semantic version | Specific releases |
| `X.Y` | Minor version | Feature branches |
| `sha-abc123` | Git commit SHA | Reproducible builds |
| `branch-name` | Feature branch | Testing PRs |

## Dockerfile Examples

### Backend Application

Create `Dockerfile`:

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Copy source code
COPY src/ ./src/
COPY tsconfig.json ./
COPY next.config.js ./

# Build application
RUN npm run build

# Production stage
FROM node:20-alpine AS production

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy built artifacts
COPY --from=builder --chown=nodejs:nodejs /app/dist ./dist
COPY --from=builder --chown=nodejs:nodejs /app/node_modules ./node_modules
COPY --from=builder --chown=nodejs:nodejs /app/package*.json ./

# Set environment
ENV NODE_ENV=production
ENV PORT=3000

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3000/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

# Start application
CMD ["node", "dist/server.js"]
```

### Frontend Application

Create `frontend/Dockerfile`:

```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source
COPY . .

# Build arguments
ARG NEXT_PUBLIC_API_URL=https://api.biometrics.example.com
ARG NEXT_PUBLIC_ENV=production

ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_ENV=$NEXT_PUBLIC_ENV

# Build Next.js application
RUN npm run build

# Production stage with Nginx
FROM nginx:alpine AS production

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built artifacts
COPY --from=builder --chown=nginx:nginx /app/.next /usr/share/nginx/html/next

# Copy static assets
COPY --from=builder --chown=nginx:nginx /app/public /usr/share/nginx/html

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
```

### Worker Service

Create `worker/Dockerfile`:

```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY src/ ./src/
COPY tsconfig.json ./

RUN npm run build

FROM node:20-alpine AS production

WORKDIR /app

RUN addgroup -g 1001 -S worker && \
    adduser -S worker -u 1001

COPY --from=builder --chown=worker:worker /app/dist ./dist
COPY --from=builder --chown=worker:worker /app/node_modules ./node_modules
COPY --from=builder --chown=worker:worker /app/package*.json ./

USER worker

ENV NODE_ENV=production

CMD ["node", "dist/worker.js"]
```

## Docker Compose

### Development Setup

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
ports:
- "53050:3000" # Port Sovereignty Compliance (Rule -9): 3000→53050
environment:
- NODE_ENV=development
- DATABASE_URL=postgresql://postgres:postgres@db:5432/biometrics
- REDIS_URL=redis://redis:6379
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./src:/app/src
    command: npm run dev

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=biometrics
ports:
- "51003:5432" # Port Sovereignty Compliance (Rule -9): 5432→51003
volumes:
- postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

redis:
image: redis:7-alpine
ports:
- "51004:6379" # Port Sovereignty Compliance (Rule -9): 6379→51004
volumes:
- redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

volumes:
  postgres_data:
  redis_data:
```

### Production Setup

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
app:
image: biometrics-org/app:latest
restart: unless-stopped
ports:
- "53050:3000" # Port Sovereignty Compliance (Rule -9): 3000→53050
environment:
- NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - API_SECRET=${API_SECRET}
healthcheck:
test: ["CMD", "curl", "-f", "http://localhost:3000/health"] # Internal container port unchanged (3000)
interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 512M

  worker:
    image: biometrics-org/app:worker
    restart: unless-stopped
    environment:
      - NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
    depends_on:
      - app
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 1G

networks:
  default:
    name: biometrics_network
```

## GitHub Actions for Docker

### Build and Push Workflow

Create `.github/workflows/docker.yml`:

```yaml
name: Build and Push Docker Image

on:
  push:
    branches: [main]
    tags: ['v*']
  pull_request:
    branches: [main]

env:
  REGISTRY: docker.io
  IMAGE_NAME: biometrics-org/app

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha

      - name: Build and push backend
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./Dockerfile
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          build-args: |
            NODE_ENV=production

      - name: Build and push frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          file: ./frontend/Dockerfile
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}-frontend
          labels: ${{ steps.meta.outputs.labels }}-frontend
          cache-from: type=gha

      - name: Build and push worker
        uses: docker/build-push-action@v5
        with:
          context: ./worker
          file: ./worker/Dockerfile
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}-worker
          labels: ${{ steps.meta.outputs.labels }}-worker
          cache-from: type=gha

  security-scan:
    runs-on: ubuntu-latest
    needs: build-and-push
    steps:
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: '${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest'
          format: 'sarif'
          exit-code: '1'
          severity: 'CRITICAL,HIGH'

      - name: Upload Trivy results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy-results.sarif
```

## Multi-Platform Builds

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        platform:
          - linux/amd64
          - linux/arm64
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: ${{ matrix.platform }}
          push: true
          tags: |
            biometrics-org/app:latest
            biometrics-org/app:${{ matrix.platform }}
```

## Image Management

### Pulling Images

```bash
# Pull latest
docker pull biometrics-org/app:latest

# Pull specific version
docker pull biometrics-org/app:1.2.3

# Pull for specific architecture
docker pull biometrics-org/app:1.2.3-arm64

# Pull all tags
docker pull biometrics-org/app --all-tags
```

### Cleanup Old Images

```bash
# Remove dangling images
docker image prune -f

# Remove unused images (older than 30 days)
docker image prune -a --filter "until=720h"

# Remove specific tag
docker rmi biometrics-org/app:old-tag
```

## Deployment

### Docker Swarm Deploy

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.prod.yml biometrics

# Scale service
docker service scale biometrics_app=5

# Check service status
docker service ls
docker service ps biometrics_app
```

### Kubernetes Deploy (Optional)

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: biometrics-app
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: biometrics-app
  template:
    metadata:
      labels:
        app: biometrics-app
    spec:
      containers:
        - name: app
          image: biometrics-org/app:latest
          ports:
            - containerPort: 3000
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: biometrics-secrets
                  key: database-url
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "2000m"
          readinessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 10
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
```

## Related Documentation

- [CI-CD-PIPELINE.md](./CI-CD-PIPELINE.md)
- [TESTING-SUITE.md](./TESTING-SUITE.md)
- [MONITORING-SETUP.md](./MONITORING-SETUP.md)

---

**End of Docker Hub Documentation**
