package tinydb

import "sync"

type DB struct {
	mu     sync.RWMutex
	data   map[string][]byte
	closed bool
}

func Open(path string) (*DB, error) {
	return &DB{
		data: make(map[string][]byte),
	}, nil
}

func (db *DB) Set(key, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrClosed
	}
	if len(key) == 0 {
		return ErrEmptyKey
	}

	db.data[string(key)] = append([]byte(nil), value...)
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.closed {
		return nil, ErrClosed
	}

	value, ok := db.data[string(key)]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return append([]byte(nil), value...), nil
}

func (db *DB) Delete(key []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrClosed
	}

	delete(db.data, string(key))
	return nil
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.closed = true
	return nil
}
