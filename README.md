# tinydb

A tiny append-only key/value database in Go.

## Features
- `Set(key, value)` writes records to a log file
- `Get(key)` reads the latest value using an in-memory index
- `Delete(key)` writes a tombstone record
- `Close()` closes the file handle safely
- Data persists across process restarts (close + reopen)

## Install

```bash
go test ./...
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/ouchbouk/tinydb"
)

func main() {
	db, err := tinydb.Open("./tiny.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Set([]byte("name"), []byte("Alice")); err != nil {
		log.Fatal(err)
	}

	v, err := db.Get([]byte("name"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(v)) // Alice
}
```

## Notes
- Keys must be non-empty (`ErrEmptyKey`).
- Missing keys return `ErrKeyNotFound`.
- Operations on a closed DB return `ErrClosed`.
