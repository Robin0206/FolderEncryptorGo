package main

import (
	"crypto/rand"
	"encoding/binary"
	"golang.org/x/crypto/chacha20"
	"strconv"
	"strings"
)

type EncData struct {
	oldPath string
	newPath string
	nonce   []byte
	salt    []byte
	mac     []byte
}

func generateEncData(path string) EncData {
	var result EncData
	result.oldPath = path
	result.newPath = generateNewPath(path)
	result.nonce = make([]byte, chacha20.NonceSize)
	result.salt = make([]byte, 16)
	rand.Read(result.nonce)
	rand.Read(result.salt)
	return result
}

func (this_ptr *EncData) toString() string {
	var result = this_ptr.oldPath + ";" + this_ptr.newPath + ";"
	for _, b := range this_ptr.salt {
		result += strconv.Itoa(int(b)) + ";"
	}
	for _, b := range this_ptr.nonce {
		result += strconv.Itoa(int(b)) + ";"
	}
	for _, b := range this_ptr.mac {
		result += strconv.Itoa(int(b)) + ";"
	}
	return result
}

func generateNewPath(path string) string {
	splitOldPath := strings.Split(path, "/")
	splitNewPath := splitOldPath[:len(splitOldPath)-1]
	newFileName := generateNewFileName()
	splitNewPath = append(splitNewPath, string(newFileName))
	var result = ""
	for i := 0; i < len(splitNewPath); i++ {
		result += "/" + splitNewPath[i]
	}
	return result
}

func generateNewFileName() string {
	var resultBytes = make([]byte, 8)
	rand.Read(resultBytes)
	resultLong := binary.LittleEndian.Uint64(resultBytes)
	return strconv.FormatUint(resultLong, 16)
}
