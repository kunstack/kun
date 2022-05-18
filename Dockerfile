FROM golang:1.18.2-bullseye as builder
WORKDIR /traffox
COPY . .
RUN go build -o traffox main.go

FROM debian:bullseye
WORKDIR /usr/lib/traffox
COPY --from=builder /traffox/traffox /usr/bin/traffox
CMD ["/usr/bin/traffox"]
