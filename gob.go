package std

import (
	"bytes"
	"encoding/gob"
)

// GobEncodeBytes
//
// See bytes.NewBuffer()
// See gob.NewEncoder().Encode().
func GobEncodeBytes(v any) ([]byte, error) {
	bb := bytes.NewBuffer(nil)

	err := gob.NewEncoder(bb).Encode(v)
	if err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

// GobDecodeBytes
//
// See bytes.NewBuffer()
// See gob.NewDecoder().Decode().
func GobDecodeBytes[T any](data []byte) (T, error) {
	p := new(T)

	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(p)
	if err != nil {
		return *p, err
	}

	return *p, nil
}
