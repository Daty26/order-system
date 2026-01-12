FROM golang:1.25-alpine
RUN apk add --no-cache git
COPY user-service/go.mod user-service/go.sum ./user-service/
COPY auth/go.mod ./auth/

WORKDIR /app/user-service
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/user-service
RUN go build -o server ./cmd/main.go
EXPOSE 8080
CMD ["./server"]
