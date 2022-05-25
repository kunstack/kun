FROM golang:1.18.2-bullseye as builder
WORKDIR /traffox
COPY . .
RUN echo $(env)
RUN go build -ldflags "-X 'github.com/aapelismith/traffox/pkg/version.Semver=$(git describe --tags --always)'  \
    -X  'github.com/aapelismith/traffox/pkg/version.BuildDate=$(git show -s --format=%cd)'  \
    -X  'github.com/aapelismith/traffox/pkg/version.GitCommit=$(git show -s --format=%H)'" -o traffox cmd/traffox/main.go

FROM debian:bullseye
RUN apt-get update && apt-get install ca-certificates -y && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/lib/traffox
COPY --from=builder /traffox/traffox /usr/bin/traffox
CMD ["/usr/bin/traffox"]
