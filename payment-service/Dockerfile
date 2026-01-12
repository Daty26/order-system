FROM golang:1.25-alpine
RUN apk add --no-cache git
COPY payment-service/go.mod payment-service/go.sum ./payment-service/
COPY auth/go.mod ./auth/

WORKDIR /app/payment-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/payment-service
RUN go build -o server ./cmd/main.go
EXPOSE 8080
CMD ["./server"]
