FROM golang:1.18-alpine

ENV APP_ENV=production
ENV GIN_MODE=release

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag fmt
RUN swag init \
  -d "./internal/server,./internal/controller/,./pkg/musicresearcher" \
  -g server.go \
  --parseInternal

# build
RUN apk add git
RUN go build -v -o /usr/local/bin/app ./cmd/server

RUN echo "building for ${APP_ENV}"

CMD app --env ${APP_ENV}
