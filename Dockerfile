# Build stage
FROM golang:1.24-alpine AS builder

# Install ANTLR4 and Task
RUN apk add --no-cache git make curl openjdk11-jre-headless

# Install Task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Install ANTLR4
RUN wget https://www.antlr.org/download/antlr-4.13.1-complete.jar -O /usr/local/lib/antlr-4.13.1-complete.jar
RUN echo 'java -jar /usr/local/lib/antlr-4.13.1-complete.jar "$@"' > /usr/local/bin/antlr4 && chmod +x /usr/local/bin/antlr4

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate ANTLR files and build
RUN ./bin/task generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qasmparser ./cmd/qasmparser

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/qasmparser .

# Expose port (if needed for future web interface)
# EXPOSE 8080

ENTRYPOINT ["./qasmparser"]