FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gobirthdays ./cmd/mybot/main.go

CMD ["./gobirthdays"]
