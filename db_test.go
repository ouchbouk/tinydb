package tinydb

import "testing"

func TestSetGet(t *testing.T) {
	db, err := Open("")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Set([]byte("name"), []byte("Alice"))
	if err != nil {
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
	db, _ := Open("")
	defer db.Close()

	_, err := db.Get([]byte("missing"))
	if err != ErrKeyNotFound {
		t.Fatalf("got %v, want %v", err, ErrKeyNotFound)
	}
}

func TestDelete(t *testing.T) {
	db, _ := Open("")
	defer db.Close()

	_ = db.Set([]byte("name"), []byte("Alice"))
	_ = db.Delete([]byte("name"))

	_, err := db.Get([]byte("name"))
	if err != ErrKeyNotFound {
		t.Fatalf("got %v, want %v", err, ErrKeyNotFound)
	}
}
