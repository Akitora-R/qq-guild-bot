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

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qq-guild-bot ./cmd/qq-guild-bot/main.go

# Use scratch as the final base image to create a minimal container
FROM scratch

# Copy the built application from the builder image
COPY --from=builder /app/qq-guild-bot .

# Run the application
CMD ["./qq-guild-bot"]
