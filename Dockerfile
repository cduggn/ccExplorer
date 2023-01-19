# dockerfile for golang CLI app ccExplorer with root command in ./cmd/ccExplorer

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Sergey Kuznetsov <"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o ccexplorer ./cmd/ccexplorer

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./ccexplorer"]





