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

package mysql

import (
	"errors"

	"github.com/spf13/pflag"
)

// Options related configuration
type Options struct {
	// The host name of the Database server
	Host string `yaml:"host,omitempty" json:"host,omitempty"`
	// The database username of the Database server
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	// Database database name
	Database string `yaml:"database,omitempty" json:"database,omitempty"`
	// Database server database password
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	// PoolLimit database connection pool size
	PoolLimit int `yaml:"pool_limit,omitempty" json:"pool_limit,omitempty"`
	// Maximum number of idle Database connections
	MaxIdle int `yaml:"max_idle,omitempty" json:"max_idle,omitempty"`
	// Maximum lifetime of Database connection
	MaxLife int `yaml:"max_life,omitempty" json:"max_life,omitempty"`
}

// Flags related command line parameters
func (o *Options) Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.StringVar(&o.Host, "mysql.host", o.Host, "The host name of the Database server")
	fs.StringVar(&o.Database, "mysql.database", o.Database, "The database name in mysql to be used by this service")
	fs.StringVar(&o.Username, "mysql.username", o.Username, "MySQL server username")
	fs.StringVar(&o.Password, "mysql.password", o.Password, "MySQL server database password")
	fs.IntVar(&o.PoolLimit, "mysql.pool-limit", o.PoolLimit, "MySQL database connection pool size")
	fs.IntVar(&o.MaxIdle, "mysql.max-idle", o.MaxIdle, "Maximum number of idle Database connections")
	fs.IntVar(&o.MaxLife, "mysql.max-life", o.MaxLife, "Maximum lifetime of Database connection")
	return fs
}

// SetDefaults sets the default values.
func (o *Options) SetDefaults() {
	o.Host = "localhost:3306"
	o.Database = "test"
	o.Username = "root"
}

// Validate verify the configuration and return an error if correct
func (o *Options) Validate() error {
	if o.Host == "" {
		return errors.New("host is required")
	}
	if o.Database == "" {
		return errors.New("database is required")
	}
	if o.Username == "" {
		return errors.New("username is required")
	}
	return nil
}

// NewOptions Create a new Options filled with default values
func NewOptions() *Options {
	opt := &Options{}
	opt.SetDefaults()
	return opt
}
