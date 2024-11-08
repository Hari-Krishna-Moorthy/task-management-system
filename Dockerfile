# Stage 1: Build the Go application
### sha256:48eab5e3505d8c8b42a06fe5f1cf4c346c167cc6a89e772f31cb9e5c301dcf60
FROM golang:1.22.7-alpine3.20 AS base 

# Using multi stage build to reduce image size & security risk
FROM base AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install necessary packages
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Install the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install the dependencies without checking for go code
RUN go mod vendor

# Build the application
RUN go build -o server cmd/server/server.go

# Stage 2: Create a lightweight container for running the application
FROM base

# Set working directory
WORKDIR /root

# Install certificates (required if you're connecting to MongoDB Atlas or HTTPS resources)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the build stage
COPY --from=builder /app/server .
COPY --from=builder /app/.env .env

# Change the ownership of the application files to the 'tms' user

RUN addgroup -g 10000 runners
RUN adduser --uid 10000 --ingroup runners -S tms
RUN chown -R tms:runners /root
USER tms

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./server"]
