# Multistage Dockerfile to build a go application and run it

# Stage 1: Build the application
FROM golang:1.19-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run the application
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the current working directory inside the container
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]