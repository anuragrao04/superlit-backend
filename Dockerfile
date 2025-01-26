# ---- Base Stage ----
FROM golang:1.23 AS base
WORKDIR /superlit-backend

# Install firejail
RUN apt-get update && apt-get install -y firejail

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend source code
COPY . .

EXPOSE 6969

# ---- Development Stage ----
FROM base AS dev
CMD ["go", "run", "main.go"]

# ---- Production Stage ----
FROM base AS prod
# Build and run the Go application
RUN go build .
CMD ["./superlit-backend"]
