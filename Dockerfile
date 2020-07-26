FROM golang:1.14

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["make","run"]

EXPOSE 5000