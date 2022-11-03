FROM golang:1.18.2-bullseye as builder
WORKDIR /kun
COPY . .
RUN go build -ldflags "-X 'github.com/aapelismith/kuntunnel/pkg/version.Semver=$(git describe --tags --always)'  \
    -X  'github.com/aapelismith/kuntunnel/pkg/version.BuildDate=$(git show -s --format=%cd)'  \
    -X  'github.com/aapelismith/kuntunnel/pkg/version.GitCommit=$(git show -s --format=%H)'" -o client cmd/client/main.go

FROM debian:bullseye
RUN apt-get update && apt-get install ca-certificates -y && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/lib/kun
COPY --from=builder /kun/kun /usr/bin/kun
CMD ["/usr/bin/kun"]
