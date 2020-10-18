package main

import (
	"encoding/hex"
	"fmt"
	"os"
)

// A ciphertext matrix with a single key
type Crib struct {
	hexCiphertexts []string
	ciphertexts    [][]byte
	key            []byte
}

func NewCrib(hexCiphertexts []string) *Crib {
	if len(hexCiphertexts) < 2 {
		fmt.Println("You must provide at least two ciphertexts. Example: \"crib-drag 315c4e 234c02 ...\"")
		os.Exit(1)
	}

	var ciphertexts [][]byte
	for _, hexCiphertext := range hexCiphertexts {
		ciphertext, err := hex.DecodeString(hexCiphertext)
		if err != nil {
			fmt.Printf("%s is not valid hexidecimal.\n", hexCiphertext)
			os.Exit(1)
		}
		ciphertexts = append(ciphertexts, ciphertext)
	}

	var key []byte

	return &Crib{hexCiphertexts, ciphertexts, key}
}

func (c *Crib) guess(pos Position, b byte) bool {
	if pos.j < len(c.ciphertexts[pos.i]) {
		k := c.ciphertexts[pos.i][pos.j] ^ b

		if pos.j < len(c.key) {
			c.key[pos.j] = k
			return true
		}

		if pos.j == len(c.key) {
			c.key = append(c.key, k)
			return true
		}
	}

	return false
}

func (c *Crib) remove(pos Position) {
	if pos.j == len(c.key) - 1 {
		c.key = c.key[:len(c.key)-1]
	}
}

func (c Crib) get(pos Position) byte {
	return c.ciphertexts[pos.i][pos.j] ^ c.key[pos.j]
}
