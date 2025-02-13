# Use the official Go image
FROM golang:latest

# Set working directory
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the source code
COPY . .

# Run tests
RUN go test ./tests -v

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Run the executable
CMD ["./main"]
