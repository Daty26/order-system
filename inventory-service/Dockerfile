FROM golang:1.25-alpine
RUN apk add --no-cache git
COPY inventory-service/go.mod inventory-service/go.sum ./inventory-service/
COPY auth/go.mod ./auth/

WORKDIR /app/inventory-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/inventory-service
RUN go build -o server ./cmd/main.go
EXPOSE 8080
CMD ["./server"]
