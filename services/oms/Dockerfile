FROM golang:1.24 AS builder

WORKDIR /app
COPY go.work go.work
COPY go.mod go.sum ./
COPY core core
COPY libs libs
COPY services services

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /oms-migrator ./services/oms/cmd/migrator
RUN CGO_ENABLED=0 GOOS=linux go build -o /oms-api ./services/oms/cmd/api

FROM alpine:3.20
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /oms-migrator /app/oms-migrator
COPY --from=builder /oms-api /app/oms-api

COPY services/oms/migrations/ ./services/oms/migrations/

EXPOSE 8080
CMD ["/app/oms-api"]
