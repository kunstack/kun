/*
 * Copyright 2021 The KunStack Authors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import (
	_ "embed"
)

//go:generate protoc --proto_path=. --proto_path=../../../../third_party --go_opt=paths=source_relative --grpc-gateway_opt=logtostderr=true  --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=allow_delete_body=true --go-grpc_opt=paths=source_relative --validate_opt=paths=source_relative --validate_opt=lang=go --openapiv2_opt=logtostderr=true  --openapiv2_opt=allow_delete_body=true  --openapiv2_opt=use_go_templates=true --openapiv2_opt=allow_merge=true   --openapiv2_out=. --grpc-gateway_out=.  --go_out=. --go-grpc_out=.  --validate_out=.  peer.proto