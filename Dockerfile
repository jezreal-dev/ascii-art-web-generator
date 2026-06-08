# --- Stage 1: Build the Go application ---
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file
COPY go.mod ./

# Copy the Go source code files
COPY main.go server.go ./

# Compile the Go application to a static binary named 'server'
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- Stage 2: Create the final lightweight container ---
FROM alpine:latest

# Set working directory in the runner container
WORKDIR /app

# Copy the compiled binary from Stage 1
COPY --from=builder /app/server .

# Copy the asset directories needed at runtime
COPY templates/ ./templates/
COPY banners/ ./banners/

# Expose the default port
EXPOSE 8080

# Command to run the executable
CMD ["./server"]