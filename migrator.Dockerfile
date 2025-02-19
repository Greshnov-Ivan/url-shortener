FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_TAGS
RUN CGO_ENABLED=0 GOOS=linux go build -tags "${BUILD_TAGS}" -o migrator ./cmd/migrator

FROM scratch

WORKDIR /migrator

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/migrator .

ENTRYPOINT ["./migrator"]