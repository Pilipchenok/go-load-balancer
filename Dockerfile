FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /balancer ./cmd/balancer


FROM alpine:3.19

COPY --from=builder /balancer /balancer
COPY configs/config.yaml /configs/config.yaml

EXPOSE 8080

CMD ["/balancer"]