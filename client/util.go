package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// ReadAndClose reads the input until EOF and then closes it.
// If reading results in an error, the error is returned and ethe input is not closed.
func ReadAndClose(r io.ReadCloser) error {
	b := make([]byte, bytes.MinRead)
	for {
		_, err := r.Read(b)
		if err == io.EOF {
			r.Close()
			return nil
		}

		if err != nil {
			return err
		}
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func mustRandomHex(n int) string {
	s, err := randomHex(n)
	if err != nil {
		panic(err)
	}

	return s
}
