FROM golang:1.25 AS builder

WORKDIR /app

COPY services/auth/go.mod services/auth/go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o auth-service \
    ./services/auth/cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/auth-service .

EXPOSE 50051

CMD ["./auth-service"]