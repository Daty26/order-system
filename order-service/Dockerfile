FROM golang:1.23-alpine
RUN apk add --no-cache git ca-certificate
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o order-service ./cmd/main.go
EXPOSE 8080
CMD ["./order-service"]