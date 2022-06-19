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
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/spf13/pflag"

	"github.com/aapelismith/kun/pkg/client/mysql"
)

func Test_Options(t *testing.T) {
	t.Run("test NewOptions", func(t *testing.T) {
		options := mysql.NewOptions()
		if err := options.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("test SetDefaults", func(t *testing.T) {
		options := &mysql.Options{}
		if err := options.Validate(); err == nil {
			t.Fatalf("expected err, got nil")
		}
		options.SetDefaults()
		if err := options.Validate(); err != nil {
			t.Fatalf("expected no err, got %v", err)
		}
	})

	t.Run("test Flags", func(t *testing.T) {
		expectHost := "8.8.8.8:3306"
		expectDatabase := "test"
		expectUsername := "root"
		expectPassword := uuid.New().String()
		expectPoolLimit := 10
		expectMaxIdle := 11
		expectMaxLife := 20
		options := &mysql.Options{}
		args := []string{
			"--mysql.host", expectHost,
			"--mysql.database", expectDatabase,
			"--mysql.username", expectUsername,
			"--mysql.password", expectPassword,
			"--mysql.pool-limit", strconv.Itoa(expectPoolLimit),
			"--mysql.max-idle", strconv.Itoa(expectMaxIdle),
			"--mysql.max-life", strconv.Itoa(expectMaxLife),
		}
		cleanFlags := pflag.NewFlagSet("", pflag.ContinueOnError)
		cleanFlags.AddFlagSet(options.Flags())
		if err := cleanFlags.Parse(args); err != nil {
			t.Fatal(err)
		}
		if err := options.Validate(); err != nil {
			t.Fatal(err)
		}
		if options.Host != expectHost {
			t.Fatalf("expected %v; got %v", expectHost, options.Host)
		}
		if options.Database != expectDatabase {
			t.Fatalf("expected %v; got %v", expectDatabase, options.Database)
		}
		if options.Username != expectUsername {
			t.Fatalf("expected %v; got %v", expectUsername, options.Username)
		}
		if options.Password != expectPassword {
			t.Fatalf("expected %v; got %v", expectPassword, options.Password)
		}
		if options.PoolLimit != expectPoolLimit {
			t.Fatalf("expected %v; got %v", expectPoolLimit, options.PoolLimit)
		}
		if options.MaxIdle != expectMaxIdle {
			t.Fatalf("expected %v; got %v", expectMaxIdle, options.MaxIdle)
		}
		if options.MaxLife != expectMaxLife {
			t.Fatalf("expected %v; got %v", expectMaxLife, options.MaxLife)
		}
	})
}
