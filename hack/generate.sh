#!/usr/bin/env bash
set -e
#
# Copyright 2022 Aapeli.Smith<aapeli.nian@gmail.com>.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

PROJECT_DIR="$(dirname "$(readlink -f "$0")")"

GO_BIN="$(go env GOBIN)"

export PATH="$PATH:${GO_BIN:-$(go env GOPATH)/bin}"

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/envoyproxy/protoc-gen-validate

go generate "${PROJECT_DIR}/../..."
