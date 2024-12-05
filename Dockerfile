FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./bin/app .

EXPOSE 1323
CMD ["./bin/app"]
