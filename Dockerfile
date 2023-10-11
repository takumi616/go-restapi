#a container creating binary which is included in deploy container
FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o main

#-------------------------------------------------------------------------

#a container which is used to deploy
FROM alpine:latest as deploy

RUN apk update

COPY --from=builder /app/main .

CMD ["./main"]
