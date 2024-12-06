FROM golang:1.23-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM base AS dev

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM base AS prod


COPY . .
RUN go build -v -o ./bin/app .

EXPOSE 1323
CMD ["./bin/app"]
