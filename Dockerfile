FROM golang:latest as builder
WORKDIR /app
COPY . .

RUN GOOS=linux go build -ldflags="-w -s" -o webserver ./cmd/server/main.go

CMD ["./webserver"]