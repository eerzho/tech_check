FROM golang:1.23-alpine

WORKDIR /http

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY ./ ./

CMD ["go", "run", "./cmd/http"]