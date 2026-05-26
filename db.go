package tinydb

import (
	"errors"
	"io"
	"os"
	"sync"
)

type DB struct {
	mu     sync.RWMutex
	file   *os.File
	data   map[string][]byte
	index  map[string]int64
	closed bool
}

func Open(path string) (*DB, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}

	db := &DB{
		file:  file,
		index: make(map[string]int64),
	}

	if err := db.loadIndex(); err != nil {
		_ = file.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) loadIndex() error {
	if _, err := db.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	for {
		offset, err := db.file.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		rec, err := readRecord(db.file)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return err
		}

		key := string(rec.key)

		if rec.flags&flagTombstone != 0 {
			delete(db.index, key)
		} else {
			db.index[key] = offset
		}
	}

	_, err := db.file.Seek(0, io.SeekEnd)
	return err
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
