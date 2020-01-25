FROM golang:1.13 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o antibrutefroce .

FROM alpine:latest
RUN apk --no-cache  add \
    ca-certificates \
    tzdata \
    && cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
WORKDIR /root/
COPY --from=builder /app/antibrutefroce .