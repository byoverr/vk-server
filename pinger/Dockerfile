FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o pinger ./cmd

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/pinger .
COPY config.yml .

RUN chmod +x pinger


EXPOSE 8080

CMD ["/app/pinger"]