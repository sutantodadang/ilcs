FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -ldflags="-s -w" -o app cmd/main.go

# Stage 2: Create a minimal image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 6565

CMD ["./app"]