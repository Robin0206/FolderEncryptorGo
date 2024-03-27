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

type TestEncryptor struct {
	encData EncData
}

type FlipEncryptor struct {
	encData EncData
}

func (this_ptr *FlipEncryptor) encrypt(password string) {
	fmt.Println(this_ptr.encData.toString())
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	defer in.Close()
	defer out.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := in.Read(buffer)
		for i := 0; i < n; i++ {
			buffer[i] ^= 255
		}

		out.Write(buffer[:n])
		if err == io.EOF {
			break
		}
	}
	os.Remove(this_ptr.encData.oldPath)
}

func (this_ptr *FlipEncryptor) decrypt(password string) {
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	defer in.Close()
	defer out.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := in.Read(buffer)
		for i := 0; i < n; i++ {
			buffer[i] ^= 255
		}

		out.Write(buffer[:n])
		if err == io.EOF {
			break
		}
	}
	os.Remove(this_ptr.encData.oldPath)
}

func (this_ptr *TestEncryptor) encrypt(password string) {
	fmt.Println(this_ptr.encData.toString())
	fmt.Println("=======================================================================================================================================================")
}

func (this_ptr *TestEncryptor) decrypt(password string) {
	fmt.Println(this_ptr.encData.toString())
}

type ChaChaEncryptor struct {
	chacha  chacha20.Cipher
	encData *EncData
	worker  *Worker
	index   int
}

func generateChaChaEncryptor(encData EncData, worker *Worker, index int) *ChaChaEncryptor {
	var result ChaChaEncryptor
	result.encData = &encData
	result.worker = worker
	result.index = index
	return &result
}

func (this_ptr *ChaChaEncryptor) encrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	mac := hmac.New(sha256.New, key)
	defer in.Close()
	defer out.Close()

	readBuffer := make([]byte, ENC_BUFFERSIZE)
	encBuffer := make([]byte, ENC_BUFFERSIZE)
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

func (this_ptr *ChaChaEncryptor) decrypt(password string) {
	key := deriveKey(password, this_ptr.encData.salt)
	chacha, _ := chacha20.NewUnauthenticatedCipher(key, this_ptr.encData.nonce)
	chacha.SetCounter(0)
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	mac := hmac.New(sha256.New, key)
	defer in.Close()
	defer out.Close()

	readBuffer := make([]byte, ENC_BUFFERSIZE)
	encBuffer := make([]byte, ENC_BUFFERSIZE)
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
	size, _ := os.Open(path)
	info, _ := size.Stat()
	fileSize := info.Size()

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
