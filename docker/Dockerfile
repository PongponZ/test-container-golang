FROM golang:1.19 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/api/main.go

FROM alpine:3.13

WORKDIR /usr/src/app

COPY --from=builder /src/main /usr/src/app/main

EXPOSE 4000
CMD ["./main"]