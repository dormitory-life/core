FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o core-service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/core-service .

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/configs ./configs

EXPOSE 8082

CMD ["./core-service", "/app/configs/config.yaml"]