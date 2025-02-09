FROM golang:1.23-alpine

WORKDIR /app

COPY . .

COPY config.yml /app/config.yml

EXPOSE 8080

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o backend/cmd/main backend/cmd/main.go

CMD ["go", "run", "backend/cmd/main"]
