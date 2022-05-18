FROM golang:1.18.2-bullseye as builder
WORKDIR /traffox
COPY . .
RUN go build -o traffox main.go

FROM debian:bullseye
RUN apt-get update && apt-get install ca-certificates -y && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/lib/traffox
COPY --from=builder /traffox/traffox /usr/bin/traffox
CMD ["/usr/bin/traffox"]
