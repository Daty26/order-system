FROM golang:1.23-alpine
RUN apk add --no-cache git ca-certificate
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment-service ./cmd/main.go
EXPOSE 8081
CMD ["./payment-service"]