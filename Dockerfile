FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/url-shortener

FROM alpine:latest

WORKDIR /url-shortener

COPY --from=builder /app/url-shortener .

CMD ["./url-shortener"]