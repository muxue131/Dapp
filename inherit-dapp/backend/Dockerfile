# Backend API Server - Multi-stage build
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /legacy-api ./main.go

# Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /legacy-api /usr/local/bin/legacy-api

EXPOSE 8080

ENTRYPOINT ["legacy-api"]
