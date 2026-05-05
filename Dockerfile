# ---------- Build Stage ----------
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ---------- Runtime Stage ----------
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose app port
EXPOSE 8080

# Run application
CMD ["./main"]