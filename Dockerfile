# Use the official Go image as the build environment
FROM golang:1.21 as builder

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io,direct"

# Set the working directory
WORKDIR /app

# Copy the go module and sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the project files into the container
COPY . .

# Install protoc compiler and Go plugins for protobuf
RUN apt-get update && apt-get install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure the target directory for protobuf files exists
RUN mkdir -p internal/pkg/stub

# Ensure PATH includes the Go bin directory
ENV PATH="$PATH:/go/bin"

# Generate protobuf files
RUN protoc --proto_path=protos/qq_guild --go_out=internal/pkg/stub --go-grpc_out=internal/pkg/stub --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/qq_guild/*.proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qq-guild-bot ./cmd/qq-guild-bot/main.go

# Use alpine as the final base image to create a minimal container
FROM alpine:latest

# Copy the built application from the builder image
COPY --from=builder /app/qq-guild-bot .

# Run the application
CMD ["./qq-guild-bot"]
