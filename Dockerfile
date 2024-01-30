# Use the official Go image with version 1.21 as the base image
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

ENV GOPROXY=direct
ENV AWS_DEFAULT_REGION=""
ENV DDB_TABLE_NAME=""
ENV PORT="8080"


# Copy the Go module files
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the Go application
# RUN CGO_ENABLED=0 GOOS=linux go build -o applicationset-dynamodb-plugin
RUN CGO_ENABLED=0 go build -o applicationset-dynamodb-plugin

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder image
COPY --from=builder /app/applicationset-dynamodb-plugin .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["./applicationset-dynamodb-plugin"]
