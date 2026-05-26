# tinydb

A tiny append-only key/value database in Go, plus a minimal TCP server/client demo.

## Features
- `Open(path)` opens (or creates) a file-backed DB
- `Set(key, value)` appends a record and updates index
- `Get(key)` reads latest value by index offset
- `Delete(key)` appends a tombstone record
- `Compact()` rewrites only live records into a new file
- `Close()` closes the DB file handle
- Data persists across close/reopen

## Project Layout
- `db.go`, `record.go`, `errors.go`: core DB
- `cmd/server/main.go`: TCP server exposing simple commands
- `cmd/client/main.go`: interactive TCP client

## Run Tests

```bash
go test ./...
```

## Library Usage

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

## TCP Demo

Start server:

```bash
go run ./cmd/server
```

Start client (new terminal):

```bash
go run ./cmd/client
```

Example commands:

```text
PING
SET name Alice
GET name
DELETE name
GET name
```

## Error Notes
- Empty keys return `ErrEmptyKey`
- Missing keys return `ErrKeyNotFound`
- Operations after close return `ErrClosed`
