FROM golang:1.25-alpine as builder
WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:3.23.3
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
