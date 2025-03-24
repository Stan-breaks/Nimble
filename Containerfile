
#############################
# Stage 1: Build Frontend Assets
#############################
FROM docker.io/library/node:18-alpine AS frontend
WORKDIR /app

# Copy package files and install dependencies with pnpm
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install

# Copy the rest of the frontend source
COPY . .

#############################
# Stage 2: Build Backend (Go + Templ)
#############################
FROM docker.io/library/golang:1.24-alpine AS backend
WORKDIR /app

# Install GCC and musl-dev for CGO support
RUN apk add --no-cache gcc musl-dev

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source (including your views/ folder)
COPY . .

# Install Templ (adjust the import path if needed)
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate view templates from your views folder
RUN templ generate ./views/

# Build the Go binary with CGO enabled and fully static linking
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-s -w -extldflags \"-static\"" -o nimblestack .

#############################
# Stage 3: Final Image
#############################
FROM scratch
# Copy the compiled binary from the backend stage
COPY --from=backend /app/nimblestack /nimblestack
# Copy generated view templates
COPY --from=backend /app/views /views
# copy the schema for the db
COPY --from=backend /app/sqlc /sqlc
# Copy the static assets built in the frontend stage
COPY --from=frontend /app/public /public

# Expose the port your app listens on
EXPOSE 8080

# Set the entrypoint to run your binary
ENTRYPOINT ["/nimblestack"]
