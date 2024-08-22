FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -v -o /app/main

CMD ["main"]