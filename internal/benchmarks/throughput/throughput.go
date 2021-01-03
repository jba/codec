// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program measures the throughput of various I/O setups.
// Run once with -initdb to store data in the databases.

/*

With io.Copy to Discard:

69237429.0 M/sec  read memory
2591.9 M/sec  read file
 279.9 M/sec  read GCS
 177.0 M/sec  local DB
  58.4 M/sec  cloud DB

44419236.8 M/sec  read memory
3104.3 M/sec  read file
 248.7 M/sec  read GCS
 177.1 M/sec  local DB
  59.4 M/sec  cloud DB

With ioutil.ReadAll:

 914.3 M/sec  read memory
 833.9 M/sec  read file
 254.8 M/sec  read GCS
 177.1 M/sec  local DB
  60.9 M/sec  cloud DB

 881.5 M/sec  read memory
 876.7 M/sec  read file
 313.0 M/sec  read GCS
 171.3 M/sec  local DB
  58.0 M/sec  cloud DB

 335.5 M/sec  read memory
 205.4 M/sec  read file
 221.8 M/sec  read GCS
  97.9 M/sec  local DB
  56.0 M/sec  cloud DB
*/

package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"cloud.google.com/go/storage"
	"github.com/jba/codec/internal/benchmarks/config"
)

var initDB = flag.Bool("initdb", false, "init the DB")

func main() {
	flag.Parse()
	cfg, err := config.Read("../config.json")
	if err != nil {
		log.Fatal(err)
	}
	memory()
	file()
	gcs(cfg.GCSBucket)
	db(cfg, "local", *initDB)
	db(cfg, "cloud", *initDB)
}

const (
	testFileSize = 1053291048
	testFileName = "testfile.gitignore"
)

func readAll(r io.Reader) (int, time.Duration) {
	runtime.GC()
	start := time.Now()
	n, err := io.Copy(ioutil.Discard, r)
	if err != nil {
		log.Fatal(err)
	}
	return int(n), time.Since(start)
}

func throughput(msg string, size int, dur time.Duration) {
	mbsec := float64(size) / (1024 * 1024) / dur.Seconds()
	fmt.Printf("%6.2f M/sec  %s\n", mbsec, msg)
}

func memory() {
	bs := make([]byte, testFileSize)
	size, dur := readAll(bytes.NewReader(bs))
	throughput("read memory", size, dur)
}

func file() {
	f, err := os.Open(testFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	size, dur := readAll(f)
	throughput("read file", size, dur)
}

func gcs(bucket string) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	r, err := client.Bucket(bucket).Object(testFileName).NewReader(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	size, dur := readAll(r)
	throughput("read GCS", size, dur)
}

func db(cfg *config.Config, kind string, init bool) {
	db, err := cfg.Connect(kind)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if init {
		initializeDB(db)
	}
	var data []byte
	start := time.Now()
	if err := db.QueryRow(`SELECT data FROM throughput WHERE name = $1`, testFileName).Scan(&data); err != nil {
		log.Fatal(err)
	}
	throughput(kind+" DB", len(data), time.Since(start))
}

func initializeDB(db *sql.DB) {
	data, err := ioutil.ReadFile(testFileName)
	if err != nil {
		log.Fatal(err)
	}
	data = data[:len(data)/2]
	exec(db, `DROP TABLE IF EXISTS throughput`)
	exec(db, `CREATE TABLE throughput (name TEXT PRIMARY KEY, data BYTEA)`)
	exec(db, `INSERT INTO throughput (name, data) VALUES($1, $2)`, testFileName, data)
}

func exec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("%s: %v", query, err)
	}
}
