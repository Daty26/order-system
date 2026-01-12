FROM golang:1.25-alpine
RUN apk add --no-cache git
COPY order-service/go.mod order-service/go.sum ./order-service/
COPY auth/go.mod ./auth/

WORKDIR /app/order-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/order-service
RUN go build -o server ./cmd/main.go
EXPOSE 8080
CMD ["./server"]
