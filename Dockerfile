FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-crud .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/go-crud .

EXPOSE 8084

ENTRYPOINT ["./go-crud"]