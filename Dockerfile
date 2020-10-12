FROM golang:1.15.2-alpine3.12 AS builder

WORKDIR /go/src/github.com/sasalatart/batcoms/

COPY . ./
RUN GOOS=linux go build -o api cmd/api/main.go
RUN GOOS=linux go build -o seeder cmd/seeder/main.go

###

FROM alpine:latest

LABEL maintainer="Sebastian Salata R-T <sa.salatart@gmail.com>"

WORKDIR /root/

COPY --from=builder /go/src/github.com/sasalatart/batcoms/config/config.yaml ./config/config.yaml
COPY --from=builder /go/src/github.com/sasalatart/batcoms/api .
COPY --from=builder /go/src/github.com/sasalatart/batcoms/seeder .

CMD ["./api"]
