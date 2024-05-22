FROM golang:1.20-alpine as build

WORKDIR /app
COPY . .
RUN cd /app/cmd/network-dependencies-webhook && CGO_ENABLED=0 GOOS=linux go build -o /webhook 

FROM alpine:3.20 as runtime

COPY --from=build /webhook /usr/local/bin/webhook
RUN chmod +x /usr/local/bin/webhook

ENTRYPOINT ["webhook"]