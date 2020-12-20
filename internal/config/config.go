// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	GCSBucket string
	LocalDB   string
	CloudDB   string
}

func Read(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	d := json.NewDecoder(f)
	d.DisallowUnknownFields()
	if err := d.Decode(&c); err != nil {
		return nil, fmt.Errorf("decoding %q: %v", filename, err)
	}
	return &c, nil
}

func (c *Config) Connect(kind string) (*sql.DB, error) {
	var driver, connstr string

	switch kind {
	case "local":
		driver = "pgx"
		connstr = c.LocalDB
	case "cloud":
		driver = "cloudsqlpostgres"
		connstr = c.CloudDB
	default:
		return nil, fmt.Errorf("unknown DB kind %q", kind)
	}
	db, err := sql.Open(driver, connstr+" sslmode=disable")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
