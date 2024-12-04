FROM golang:latest

WORKDIR /backtest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

EXPOSE 4242
# RUN go build -o main main.go

