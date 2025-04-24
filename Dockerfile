# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o hustler ./cmd/e2e-test

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/hustler .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./hustler"]
