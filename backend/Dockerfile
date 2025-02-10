FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o backend ./cmd

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/backend .
COPY config.yml .

RUN chmod +x backend


EXPOSE 8080

CMD ["/app/backend"]
