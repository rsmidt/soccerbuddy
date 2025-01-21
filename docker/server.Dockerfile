# Stage 1: Builder.
# Use the official Golang image as the build environment.
FROM golang:1.23 AS builder

RUN BIN="/usr/local/bin" && \
VERSION="1.46.0" && \
curl -sSL \
"https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
-o "${BIN}/buf" && \
chmod +x "${BIN}/buf"

# Set environment variables.
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the entire project into the container.
COPY . .

# Compile all protobuf defintions.
RUN buf generate --clean

# Generate all code.
RUN go generate ./...

# Build the Go binary.
# -o specifies the output binary name.
RUN go build -o soccerbuddy ./cmd/start

# Stage 2: Final Image.
# Use a minimal Alpine image for the runtime.
FROM alpine:latest

# Install CA certificates (optional, if your app requires HTTPS).
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory.
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/soccerbuddy .
COPY --from=builder /app/migrations ./migrations

# Expose the application's port.
EXPOSE 4488

# Command to run the executable.
CMD ["./soccerbuddy"]
