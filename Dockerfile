FROM golang:1.21 as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o hustler-trading-bot ./cmd/hustler

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/hustler-trading-bot .

# Expose the UI port
EXPOSE 8080

# Run the application
CMD ["./hustler-trading-bot", "--config", "config.json"]
