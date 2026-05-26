package tinydb

import (
	"path/filepath"
	"testing"
)

func openTempDB(t *testing.T) (*DB, string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "tinydb.data")
	db, err := Open(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	return db, path
}

func TestSetGet(t *testing.T) {
	db, _ := openTempDB(t)
	defer db.Close()

	if err := db.Set([]byte("name"), []byte("Alice")); err != nil {
		t.Fatal(err)
	}

	got, err := db.Get([]byte("name"))
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "Alice" {
		t.Fatalf("got %q, want %q", got, "Alice")
	}
}

func TestMissingKey(t *testing.T) {
	db, _ := openTempDB(t)
	defer db.Close()

	_, err := db.Get([]byte("missing"))
	if err != ErrKeyNotFound {
		t.Fatalf("got %v, want %v", err, ErrKeyNotFound)
	}
}

func TestDelete(t *testing.T) {
	db, _ := openTempDB(t)
	defer db.Close()

	if err := db.Set([]byte("name"), []byte("Alice")); err != nil {
		t.Fatal(err)
	}
	if err := db.Delete([]byte("name")); err != nil {
		t.Fatal(err)
	}

	_, err := db.Get([]byte("name"))
	if err != ErrKeyNotFound {
		t.Fatalf("got %v, want %v", err, ErrKeyNotFound)
	}
}

func TestPersistenceAcrossReopen(t *testing.T) {
	db, path := openTempDB(t)

	if err := db.Set([]byte("name"), []byte("Alice")); err != nil {
		t.Fatal(err)
	}
	if err := db.Set([]byte("city"), []byte("Casablanca")); err != nil {
		t.Fatal(err)
	}

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()

	name, err := reopened.Get([]byte("name"))
	if err != nil {
		t.Fatal(err)
	}
	if string(name) != "Alice" {
		t.Fatalf("got %q, want %q", name, "Alice")
	}

	city, err := reopened.Get([]byte("city"))
	if err != nil {
		t.Fatal(err)
	}
	if string(city) != "Casablanca" {
		t.Fatalf("got %q, want %q", city, "Casablanca")
	}
}
