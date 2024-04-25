# Use a Go image as the base
FROM golang:latest as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o your-service-name ./cmd/your-service-name

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the final container
COPY --from=builder /app/your-service-name .

# Expose the port your application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./your-service-name"]