FROM golang AS builder
WORKDIR /src
COPY main.go /src
COPY go.mod /src
COPY go.sum /src
RUN go build -o /auth-example -tags 'netgo osusergo static_build'

FROM alpine
COPY --from=builder /auth-example /auth-example 
RUN chmod +x /auth-example
ENTRYPOINT ["/auth-example"]
