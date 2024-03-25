package main

import (
	"fmt"
	"golang.org/x/crypto/chacha20"
	"io"
	"os"
)

const ENC_BUFFERSIZE = 100000

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
}

func generateChaChaEncryptor(encData EncData) *ChaChaEncryptor {
	var result ChaChaEncryptor
	result.encData = &encData
	return &result
}

func (this_ptr *ChaChaEncryptor) encrypt(password string) {
	chacha, _ := chacha20.NewUnauthenticatedCipher(deriveKey(password, this_ptr.encData.salt), this_ptr.encData.nonce)
	chacha.SetCounter(0)
	out, _ := os.OpenFile(this_ptr.encData.newPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	in, _ := os.Open(this_ptr.encData.oldPath)
	defer in.Close()
	defer out.Close()

	readBuffer := make([]byte, ENC_BUFFERSIZE)
	encBuffer := make([]byte, ENC_BUFFERSIZE)
	for {
		n, err := in.Read(readBuffer)
		chacha.XORKeyStream(encBuffer, readBuffer)
		out.Write(encBuffer[:n])
		if err == io.EOF {
			break
		}
	}
	safeDelete(this_ptr.encData.oldPath)
}

func (this_ptr *ChaChaEncryptor) decrypt(password string) {
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
