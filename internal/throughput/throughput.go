// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program measures the throughput of various I/O setups.
// Run once with -initdb to store data in the databases.

/*

With io.Copy:

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
	"bufio"
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
	"strings"
	"time"

	"cloud.google.com/go/storage"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var initDB = flag.Bool("initdb", false, "init the DB")

func main() {
	flag.Parse()
	config, err := readConfig("config.gitignore")
	if err != nil {
		log.Fatal(err)
	}
	memory()
	file()
	gcs(config["bucket"])
	db(*initDB, "local DB", "pgx", config["localDB"])
	db(*initDB, "cloud DB", "cloudsqlpostgres", config["cloudDB"])
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
	fmt.Printf("%6.1f M/sec  %s\n", mbsec, msg)
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

func db(init bool, msg, driver, connString string) {
	db := connect(driver, connString)
	defer db.Close()
	if init {
		initializeDB(db)
	}
	var data []byte
	start := time.Now()
	if err := db.QueryRow(`SELECT data FROM throughput WHERE name = $1`, testFileName).Scan(&data); err != nil {
		log.Fatal(err)
	}
	throughput(msg, len(data), time.Since(start))
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

////////////////////////////////////////////////////////////////

func readConfig(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := map[string]string{}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		i := strings.IndexAny(line, " \t")
		if i < 0 {
			return nil, fmt.Errorf("bad line: %q", line)
		}
		m[strings.TrimSpace(line[:i])] = strings.TrimSpace(line[i:])
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func connect(driver, connString string) *sql.DB {
	db, err := sql.Open(driver, connString+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}
	return db
}

func exec(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
	if err != nil {
		log.Fatalf("%s: %v", query, err)
	}
}
