FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server-calendar ./cmd/app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server-calendar .
COPY cfg/config.yaml ./cfg/config.yaml
EXPOSE 8081
CMD ["./server-calendar"]
