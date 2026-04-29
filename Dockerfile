# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /app/server .

ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
