# Use an official Golang image as a build environment
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app and place the output binary in /app/main
WORKDIR /app/cmd
RUN go build -o /app/main .

# Use Alpine as a minimal base image to run the binary
FROM alpine:latest

# Install necessary runtime dependencies (like libc)
RUN apk --no-cache add libc6-compat

# Set the working directory in the runtime environment
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the .env file into the working directory
COPY --from=builder /app/.env /app/.env

# Expose the port that the app runs on
EXPOSE 50051

# Command to run the executable
CMD ["/app/main"]
