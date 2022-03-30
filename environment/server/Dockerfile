# build stage
FROM golang:1.18 as builder

COPY server /go/src/nek0meshi/line-messaging-app

WORKDIR /go/src/nek0meshi/line-messaging-app

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o app .


# deploy stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go/src/nek0meshi/line-messaging-app/app ./

CMD ["./app"]
