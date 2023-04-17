FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN cd /app/cmd/network-resource-info && CGO_ENABLED=0 GOOS=linux go build -o /app/app

FROM alpine:latest

COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["/app/app"]