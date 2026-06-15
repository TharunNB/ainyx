# Stage 1: Build binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server/main.go

# Stage 2: Runtime image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations ./db/migrations
EXPOSE 8080
ENTRYPOINT ["./server"]