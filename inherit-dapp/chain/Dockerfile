# Legacy Chain Node - Multi-stage build
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY chain/go.mod chain/go.sum ./
RUN go mod download

# Copy source
COPY chain/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /legacyd ./cmd/legacyd/

# Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /legacyd /usr/local/bin/legacyd

EXPOSE 26657 1317 9090

ENTRYPOINT ["legacyd"]
