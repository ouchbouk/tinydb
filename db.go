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

	offset, err := db.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	rec := record{
		key:   append([]byte(nil), key...),
		value: append([]byte(nil), value...),
		flags: 0,
	}

	if err := writeRecord(db.file, rec); err != nil {
		return err
	}

	if err := db.file.Sync(); err != nil {
		return err
	}

	db.index[string(key)] = offset
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.closed {
		return nil, ErrClosed
	}

	offset, ok := db.index[string(key)]
	if !ok {
		return nil, ErrKeyNotFound
	}

	if _, err := db.file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	rec, err := readRecord(db.file)
	if err != nil {
		return nil, err
	}

	if rec.flags&flagTombstone != 0 {
		return nil, ErrKeyNotFound
	}

	return append([]byte(nil), rec.value...), nil
}

func (db *DB) Delete(key []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrClosed
	}
	if len(key) == 0 {
		return ErrEmptyKey
	}

	rec := record{
		key:   append([]byte(nil), key...),
		value: nil,
		flags: flagTombstone,
	}

	if err := writeRecord(db.file, rec); err != nil {
		return err
	}

	if err := db.file.Sync(); err != nil {
		return err
	}

	delete(db.index, string(key))
	return nil
}

func (db *DB) Compact() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrClosed
	}

	path := db.file.Name()
	tmpPath := path + ".compact"

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	newIndex := make(map[string]int64)

	for key, offset := range db.index {
		if _, err := db.file.Seek(offset, io.SeekStart); err != nil {
			_ = tmpFile.Close()
			return err
		}

		rec, err := readRecord(db.file)
		if err != nil {
			_ = tmpFile.Close()
			return err
		}

		newOffset, err := tmpFile.Seek(0, io.SeekEnd)
		if err != nil {
			_ = tmpFile.Close()
			return err
		}

		if err := writeRecord(tmpFile, rec); err != nil {
			_ = tmpFile.Close()
			return err
		}

		newIndex[key] = newOffset
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := db.file.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	db.file = newFile
	db.index = newIndex

	return nil
}
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrClosed
	}

	db.closed = true
	return db.file.Close()
}
