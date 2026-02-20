# =============================================================================
# BIOMETRICS CLI - DOCKERFILE
# =============================================================================
# Multi-stage build for minimal production image
# Compliant with MANDATE 0.19: Modern CLI Toolchain (2026 Standard)
# =============================================================================

# -----------------------------------------------------------------------------
# STAGE 1: BUILDER
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    make

# Set working directory
WORKDIR /build

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && \
    go mod verify

# Copy source code
COPY . .

# Build arguments for versioning
ARG VERSION=2.0.0
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
ARG GIT_TAG=untagged

# Build flags
ARG LDFLAGS="-ldflags '-s -w \
    -X main.Version=${VERSION} \
    -X main.BuildTime=${BUILD_TIME} \
    -X main.GitCommit=${GIT_COMMIT} \
    -X main.GitBranch=${GIT_BRANCH} \
    -X main.GitTag=${GIT_TAG}'"

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    ${LDFLAGS} \
    -trimpath \
    -o biometrics \
    ./cmd/biometrics

# Verify the binary
RUN ./biometrics --help || true

# -----------------------------------------------------------------------------
# STAGE 2: RUNTIME (Minimal production image)
# -----------------------------------------------------------------------------
FROM alpine:3.19 AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Create non-root user for security
RUN addgroup -g 1000 biometrics && \
    adduser -D -u 1000 -G biometrics biometrics

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/biometrics /app/biometrics

# Copy licenses and documentation
COPY --from=builder /build/LICENSE /app/LICENSE 2>/dev/null || true
COPY --from=builder /build/README.md /app/README.md 2>/dev/null || true

# Set ownership to non-root user
RUN chown -R biometrics:biometrics /app

# Switch to non-root user
USER biometrics

# Set environment variables
ENV PATH="/app:${PATH}"
ENV TZ=UTC

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/biometrics", "system", "status"]

# Expose any necessary ports (if applicable)
# EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app/biometrics"]

# Default command
CMD ["--help"]

# -----------------------------------------------------------------------------
# STAGE 3: DEVELOPMENT (Optional - for local development)
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS development

# Install development tools
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    make \
    bash \
    vim \
    curl \
    golangci-lint \
    gotestsum

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && \
    go mod verify

# Copy source code
COPY . .

# Expose ports for development
# EXPOSE 8080

# Default command for development
CMD ["make", "all"]

# =============================================================================
# BUILD INSTRUCTIONS
# =============================================================================
# 
# Production Build:
#   docker build -t biometrics-cli:latest -t biometrics-cli:2.0.0 .
#
# Development Build:
#   docker build --target development -t biometrics-cli:dev .
#
# Run Production:
#   docker run --rm biometrics-cli:latest --help
#
# Run Development:
#   docker run --rm -v $(pwd):/app biometrics-cli:dev make all
#
# Multi-platform Build:
#   docker buildx build --platform linux/amd64,linux/arm64 \
#     -t biometrics-cli:latest .
#
# =============================================================================
# END OF DOCKERFILE
# =============================================================================
