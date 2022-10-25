FROM golang:1.18-alpine

WORKDIR /usr/src/app

ENV GIN_MODE=release

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apk add git

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server

CMD ["app"]
