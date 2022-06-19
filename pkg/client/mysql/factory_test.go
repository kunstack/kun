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

package mysql_test

import (
	"context"
	"testing"

	"github.com/aapelismith/kun/pkg/client/mysql"
)

func Test_NewClient(t *testing.T) {
	options := &mysql.Options{
		Host:     "localhost",
		Username: "root",
		Database: "test",
		Password: "",
	}
	if err := options.Validate(); err != nil {
		t.Fatal(err)
	}
	db, err := mysql.NewClient(context.TODO(), options)
	if err != nil {
		t.Fatal(err)
	}
	if result := db.Exec("SELECT 1"); result.Error != nil {
		t.Fatal(result.Error)
	}
}
