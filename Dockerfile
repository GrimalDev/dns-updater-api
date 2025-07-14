FROM golang:1.23-alpine

WORKDIR /app

COPY main.go .

RUN go mod init dns-updater && \
    go get github.com/labstack/echo/v4 && \
    go build -o dns-updater

EXPOSE 8080

CMD ["./dns-updater"]
