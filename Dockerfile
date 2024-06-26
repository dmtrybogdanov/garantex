FROM golang:1.22-alpine AS builder

COPY . /github.com/dmtrybogdanov/garantex/source/
WORKDIR /github.com/dmtrybogdanov/garantex/source/

RUN go mod download
RUN go build -o ./bin/garantex cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/dmtrybogdanov/garantex/source/bin/garantex .

CMD ["./garantex"]