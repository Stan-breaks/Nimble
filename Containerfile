
#############################
# Stage 1: Build Go Binary
#############################
FROM docker.io/library/golang:1.24-alpine AS builder
WORKDIR /app

# Install GCC and musl-dev for CGO support (required by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary with CGO enabled and fully static linking
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w -extldflags \"-static\"" -o nimble .

#############################
# Stage 2: Final Image
#############################
FROM scratch

# Copy the compiled binary
COPY --from=builder /app/nimble /nimble

# Copy the schema for DB initialization
COPY --from=builder /app/sqlc /sqlc

# Expose the port
EXPOSE 8080

# Run
ENTRYPOINT ["/nimble"]
