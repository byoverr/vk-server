FROM golang:1.23-alpine

WORKDIR /app

COPY . .

COPY config.yml /app/config.yml

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o app/pinger ./cmd/main.go

CMD ["go", "run", "pinger/cmd/main"]