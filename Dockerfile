# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
# Point to your specific main file location
RUN go build -o main ./cmd/api/main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# If you have a .env for local dev, Render will use Environment Variables instead
EXPOSE 8080
CMD ["./main"]