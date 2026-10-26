FROM golang:1.25 AS builder

WORKDIR /app

COPY services/products/go.mod services/products/go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o products-service \
    ./services/products/cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/products-service .

EXPOSE 50051

CMD ["./products-service"]