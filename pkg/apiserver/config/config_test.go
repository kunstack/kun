/*
 * Copyright 2021 Aapeli.Smith<aapeli.nian@gmail.com>.
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

package config_test

import (
	"fmt"
	"github.com/aapelismith/kun/pkg/apiserver/config"
	"sigs.k8s.io/yaml"
)

import "testing"

const data = `
log:
  level: debug
encoderConfig: 
  levelEncoder: capital
`

func TestConfiguration_Flags(t *testing.T) {
	cfg := &config.Configuration{}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		t.Fatal(err)
	}

	fmt.Println(cfg.Log.Level)
}
