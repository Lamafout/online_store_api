FROM golang:1.24 AS builder

WORKDIR /app
COPY go.work go.work
COPY go.mod go.sum ./
COPY libs libs
COPY services services

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o audit-consumer ./services/audit/cmd/audit_order_create
RUN CGO_ENABLED=0 GOOS=linux go build -o audit-migrator ./services/audit/cmd/migrator

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/audit-consumer /app/audit-consumer
COPY --from=builder /app/audit-migrator /app/audit-migrator
COPY --from=builder /app/services/audit/migrations /app/services/audit/migrations

CMD ["/app/audit-consumer"]
