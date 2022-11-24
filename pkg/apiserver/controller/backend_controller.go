/*
Copyright 2021 The KunStack Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	v1 "github.com/aapelismith/kun/pkg/apiserver/apis/v1"
)

var _ v1.BackendControllerServer = (*BackendController)(nil)

type BackendController struct {
	v1.BackendControllerServer
}

func (b *BackendController) WatchTunnels(request *v1.WatchTunnelsRequest, server v1.BackendController_WatchTunnelsServer) error {
	//TODO implement me
	panic("implement me")
}

func (b *BackendController) ConnectTunnel(server v1.BackendController_ConnectTunnelServer) error {
	//TODO implement me
	panic("implement me")
}

func (b *BackendController) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}
