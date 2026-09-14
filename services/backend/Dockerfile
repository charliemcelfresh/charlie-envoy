FROM golang:1.27-alpine

WORKDIR /app

RUN go install github.com/go-delve/delve/cmd/dlv@master

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -gcflags="all=-N -l" -o /tmp/app .

EXPOSE 8080 2345

CMD ["dlv", "exec", "/tmp/app", "--headless", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--continue"]
