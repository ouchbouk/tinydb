package tinydb

import (
	"encoding/binary"
	"io"
)

const (
	headerSize    = 9
	flagTombstone = byte(1)
)

type Record struct {
	key   []byte
	value []byte
	flags byte
}

func writeRecord(w io.Writer, r Record) error {
	header := make([]byte, headerSize)

	binary.LittleEndian.PutUint32(header[0:4], uint32(len(r.key)))
	binary.LittleEndian.PutUint32(header[4:8], uint32(len(r.value)))
	header[8] = r.flags

	if _, err := w.Write(header); err != nil {
		return err
	}
	if _, err := w.Write(r.key); err != nil {
		return err
	}
	if _, err := w.Write(r.value); err != nil {
		return err
	}

	return nil
}

func readRecord(r io.Reader) (Record, error) {
	header := make([]byte, headerSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return Record{}, err
	}

	keyLen := binary.LittleEndian.Uint32(header[0:4])
	valueLen := binary.LittleEndian.Uint32(header[4:8])
	flags := header[8]

	key := make([]byte, keyLen)
	value := make([]byte, valueLen)

	if _, err := io.ReadFull(r, key); err != nil {
		return Record{}, err
	}
	if _, err := io.ReadFull(r, value); err != nil {
		return Record{}, err
	}

	return Record{
		key:   key,
		value: value,
		flags: flags,
	}, nil
}

func readRecordAt(r io.ReaderAt, offset int64) (Record, error) {
	header := make([]byte, headerSize)
	if _, err := r.ReadAt(header, offset); err != nil {
		return Record{}, err
	}

	keyLen := binary.LittleEndian.Uint32(header[0:4])
	valueLen := binary.LittleEndian.Uint32(header[4:8])
	flags := header[8]

	key := make([]byte, keyLen)
	value := make([]byte, valueLen)

	pos := offset + headerSize
	if _, err := r.ReadAt(key, pos); err != nil {
		return Record{}, err
	}
	pos += int64(keyLen)
	if _, err := r.ReadAt(value, pos); err != nil {
		return Record{}, err
	}

	return Record{
		key:   key,
		value: value,
		flags: flags,
	}, nil
}
