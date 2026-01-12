FROM golang:1.25-alpine
RUN apk add --no-cache git
COPY notification-service/go.mod notification-service/go.sum ./notification-service/
COPY auth/go.mod ./auth/

WORKDIR /app/notification-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/notification-service
RUN go build -o server ./cmd/main.go
EXPOSE 8080
CMD ["./server"]
