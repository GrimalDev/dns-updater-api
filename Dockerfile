FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache dnsmasq

COPY main.go .
COPY dnsmasq.conf /etc/dnsmasq.conf
COPY dnsmasq.d/ /etc/dnsmasq.d/

# Build Go application
RUN go mod init dns-updater && \
    go get github.com/labstack/echo/v4 && \
    go build -o dns-updater

# Expose ports: 8080 for API, 53/udp for DNS
EXPOSE 8080
EXPOSE 53/udp

# Start dnsmasq in background and Go API
CMD ["/bin/sh", "-c", "dnsmasq -k & ./dns-updater"]
