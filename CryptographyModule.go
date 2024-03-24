package main

import (
	"crypto/sha512"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, chacha20.KeySize, sha512.New)
}
