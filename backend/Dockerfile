# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies (gcc, musl-dev for CGO)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the unified server
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/server/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/server /app/server

# Create directories
RUN mkdir -p /app/uploads /app/data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
CMD ["./server"]

