# Use the official Golang image as a base
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Ensure migration files are copied into the image
COPY db/migrations /app/db/migrations

# Build the Go app
RUN go build -o main .

# Expose port 1323 to the outside world
EXPOSE 1323

# Command to run the executable
CMD ["./main"]
