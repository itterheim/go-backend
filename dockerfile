FROM golang:1.24.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]