package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/chacha20"
	"io"
	"os"
)

var ENC_BUFFERSIZE = 1000000

type Encryptor interface {
	encrypt(password string)
	decrypt(password string)
}

type ChaChaStreamingEncryptor struct {
	chacha  chacha20.Cipher
	encData *EncData
	worker  *Worker
	index   int
}

func generateChaChaStreamingEncryptor(encData EncData, worker *Worker, index int) *ChaChaStreamingEncryptor {
	var result ChaChaStreamingEncryptor
	result.encData = &encData
	result.worker = worker
	result.index = index
	return &result
}

func (this_ptr *ChaChaStreamingEncryptor) encrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	mac := hmac.New(sha256.New, key)
	defer in.Close()
	defer out.Close()

	var readBuffer = this_ptr.worker.readBuffer
	var encBuffer = this_ptr.worker.encBuffer
	for {
		n, err := in.Read(readBuffer)
		mac.Write(readBuffer[:n])
		chacha.XORKeyStream(encBuffer, readBuffer)
		out.Write(encBuffer[:n])
		if err == io.EOF {
			break
		}
	}
	this_ptr.worker.pool.addMacAt(this_ptr.index, mac.Sum(nil))
	safeDelete(this_ptr.encData.oldPath)
}

func (this_ptr *ChaChaStreamingEncryptor) decrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	mac := hmac.New(sha256.New, key)
	defer in.Close()
	defer out.Close()

	var readBuffer = this_ptr.worker.readBuffer
	var encBuffer = this_ptr.worker.encBuffer
	for {
		n, err := in.Read(readBuffer)
		chacha.XORKeyStream(encBuffer, readBuffer)
		mac.Write(encBuffer[:n])
		out.Write(encBuffer[:n])
		if err == io.EOF {
			break
		}
	}

	if !hmac.Equal(mac.Sum(nil), this_ptr.encData.mac) {
		fmt.Println(this_ptr.encData.newPath + ": Message Authentication failed!")
	}
	os.Remove(this_ptr.encData.oldPath)
}

func safeDelete(path string) {
	fileSize := getFileSize(path)

	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	file.Seek(0, io.SeekStart)
	writer := io.WriterAt(file)

	zeros := make([]byte, ENC_BUFFERSIZE)
	for i := 0; i < len(zeros); i++ {
		zeros[i] = 0
	}
	for i := 0; i < int(fileSize); i += ENC_BUFFERSIZE {
		writer.WriteAt(zeros, int64(i))
	}
	writer.WriteAt(zeros, 0)
	file.Close()
	os.Remove(path)
}

func getFileSize(path string) int64 {
	size, _ := os.Open(path)
	info, _ := size.Stat()
	fileSize := info.Size()
	err := size.Close()
	if err != nil {
		fmt.Println("Error while reading fileSize of " + path)
	}
	return fileSize
}

type ChaChaInPlaceEncryptor struct {
	chacha  chacha20.Cipher
	encData *EncData
	worker  *Worker
	index   int
}

func generateChaChaInPlaceEncryptor(encData EncData, worker *Worker, index int) *ChaChaInPlaceEncryptor {
	var result ChaChaInPlaceEncryptor
	result.encData = &encData
	result.worker = worker
	result.index = index
	return &result
}

func (this_ptr *ChaChaInPlaceEncryptor) encrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	mac := hmac.New(sha256.New, key)

	inputData, _ := os.ReadFile(this_ptr.encData.oldPath)
	outputData := make([]byte, len(inputData))
	mac.Write(inputData)
	chacha.XORKeyStream(outputData, inputData)
	err := os.WriteFile(this_ptr.encData.oldPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error overwriting: " + this_ptr.encData.oldPath)
	}
	os.Rename(this_ptr.encData.oldPath, this_ptr.encData.newPath)
	this_ptr.worker.pool.addMacAt(this_ptr.index, mac.Sum(nil))
}

func (this_ptr *ChaChaInPlaceEncryptor) decrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	mac := hmac.New(sha256.New, key)

	inputData, _ := os.ReadFile(this_ptr.encData.oldPath)
	outputData := make([]byte, len(inputData))
	chacha.XORKeyStream(outputData, inputData)
	mac.Write(outputData)
	if !hmac.Equal(mac.Sum(nil), this_ptr.encData.mac) {
		fmt.Println(this_ptr.encData.newPath + ": Message Authentication failed!")
	}
	err := os.WriteFile(this_ptr.encData.oldPath, outputData, 0644)
	if err != nil {
		fmt.Println("Error overwriting: " + this_ptr.encData.oldPath)
	}
	os.Rename(this_ptr.encData.oldPath, this_ptr.encData.newPath)
}
