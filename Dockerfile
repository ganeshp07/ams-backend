## Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder ./seed/seed.go

## Runtime stage
FROM alpine:3.18
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server     ./server
COPY --from=builder /app/seeder     ./seeder
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./server"]
