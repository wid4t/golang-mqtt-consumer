FROM golang:1.22-alpine3.20 as builder
WORKDIR /golang-mqtt-consumer 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" .
FROM busybox
WORKDIR /golang-mqtt-consumer
COPY --from=builder /golang-mqtt-consumer  /usr/bin/
ENTRYPOINT ["golang-mqtt-consumer"]