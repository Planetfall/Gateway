FROM golang:1.21-alpine

ENV GIN_MODE=release

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag fmt
RUN swag init \
  -d "./internal/service,./internal/controller" \
  -g service.go \
  --parseInternal

# build
RUN apk add git
RUN go build -v -o /usr/local/bin/app ./cmd/server

CMD app
