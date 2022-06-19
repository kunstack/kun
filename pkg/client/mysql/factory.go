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
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewClient Create a database instance
func NewClient(ctx context.Context, options *Options) (*gorm.DB, error) {
	var dsn = fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		options.Username, options.Password, options.Host, options.Database,
	)
	orm, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	db, err := orm.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(options.MaxIdle)
	db.SetMaxOpenConns(options.PoolLimit)
	db.SetConnMaxLifetime(time.Duration(options.MaxLife) * time.Second)
	return orm, nil
}
