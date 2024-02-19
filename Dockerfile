FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY config /etc/conf
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main.go
RUN go build -o ./connect-emp ./cmd/main.go
RUN rm -rf ./cmd ./internal go

EXPOSE 8082
CMD ["/app/connect-emp", "serve"]

LABEL org.opencontainers.image.source https://github.com/mrexmelle/connect-emp