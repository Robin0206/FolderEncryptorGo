package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Worker struct {
	pool       *Workerpool
	encryptor  Encryptor
	password   string
	wg         *sync.WaitGroup
	readBuffer []byte
	encBuffer  []byte
}

type Workerpool struct {
	workers            []Worker
	encData            []EncData
	alreadyEnDeCrypted []bool
	mutex              sync.Mutex
	wg                 sync.WaitGroup
	path               string
	encrypting         bool
	password           string
}

func generateEncryptingWorkerpool(numThreads int, path string, password string) *Workerpool {
	var result Workerpool
	result.password = password
	result.encrypting = true
	result.path = path
	result.encData = getEncDataArr(path)
	result.alreadyEnDeCrypted = make([]bool, len(result.encData))
	for i := 0; i < len(result.alreadyEnDeCrypted); i++ {
		result.alreadyEnDeCrypted[i] = false
	}
	result.mutex = sync.Mutex{}
	result.workers = make([]Worker, numThreads)
	for i := 0; i < numThreads; i++ {
		result.workers[i] = generateWorker(&result, password)
		result.wg.Add(1)
	}
	return &result
}

func generateDecryptingWorkerpool(numThreads int, path string, password string) *Workerpool {
	var result Workerpool
	result.password = password
	result.encrypting = false
	result.path = path
	result.encData = result.getEncDataArrFromFile(path)
	result.alreadyEnDeCrypted = make([]bool, len(result.encData))
	for i := 0; i < len(result.alreadyEnDeCrypted); i++ {
		result.alreadyEnDeCrypted[i] = false
	}
	result.mutex = sync.Mutex{}
	result.workers = make([]Worker, numThreads)
	for i := 0; i < numThreads; i++ {
		result.workers[i] = generateWorker(&result, password)
		result.wg.Add(1)
	}
	return &result
}

func (this_ptr *Workerpool) getEncDataArrFromFile(path string) []EncData {
	var result []EncData
	err := this_ptr.decryptEncDatatable()
	file, _ := os.Open(path + "/EncData")
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			result = append(result, encDataFromString(scanner.Text()))
		}
	} else {
		fmt.Println(err.Error() + ": " + "Wrong password?")
		os.Exit(1)
	}
	file.Close()
	return result
}

func generateWorker(w *Workerpool, password string) Worker {
	var result Worker
	result.pool = w
	result.password = password
	result.wg = &(w.wg)
	result.readBuffer = make([]byte, ENC_BUFFERSIZE)
	result.encBuffer = make([]byte, ENC_BUFFERSIZE)
	return result
}

func (this_ptr *Workerpool) getNewEncDataIndex() int {
	this_ptr.mutex.Lock()
	var index = -1
	for i := 0; i < len(this_ptr.encData); i++ {
		if !this_ptr.alreadyEnDeCrypted[i] {
			index = i
			this_ptr.alreadyEnDeCrypted[i] = true
			break
		}
	}
	this_ptr.mutex.Unlock()
	return index
}

func (this_ptr *Workerpool) run() {
	for _, worker := range this_ptr.workers {
		go worker.run()
	}
	this_ptr.wg.Wait()
	if this_ptr.encrypting {
		this_ptr.writeEncData()
	} else {
		os.Remove(this_ptr.path + "/EncData")
	}
}

func (this_ptr *Workerpool) writeEncData() {
	out, _ := os.OpenFile(this_ptr.path+"/EncData", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	for i := 0; i < len(this_ptr.encData); i++ {
		out.WriteString(this_ptr.encData[i].toString() + "\n")
	}
	this_ptr.encryptEncDataTable()
}

func (this_ptr *Workerpool) decryptEncDatatable() error {

	//read the EncData file
	var path = this_ptr.path + "/EncData"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}

	//extract salt, nonce and content
	var salt = make([]byte, 16)
	var nonce = make([]byte, chacha20poly1305.NonceSize)
	var msgBytes = make([]byte, len(content)-len(salt)-len(nonce))
	for i := 0; i < len(salt); i++ {
		salt[i] = content[i]
	}
	for i := 0; i < len(nonce); i++ {
		nonce[i] = content[len(salt)+i]
	}
	for i := 0; i < len(msgBytes); i++ {
		msgBytes[i] = content[len(salt)+len(nonce)+i]
	}
	//generate key
	key := deriveKey(this_ptr.password, salt)

	//decrypt the content
	var encryptor, _ = chacha20poly1305.New(key)
	var decrypted []byte
	decrypted, err = encryptor.Open(decrypted, nonce, msgBytes, nil)
	if err == nil {
		//delete old encData file
		safeDelete(path)
		//write content
		//write out content
		out, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		out.Write(decrypted)
		//close output stream
		out.Close()
	} else {
		fmt.Println(err.Error())
		fmt.Println("Maybe wrong password?")
		os.Exit(1)
	}

	return err
}

func (this_ptr *Workerpool) encryptEncDataTable() {
	//generate salt nonce and key
	var salt = make([]byte, 16)
	var nonce = make([]byte, chacha20poly1305.NonceSize)
	rand.Read(nonce)
	rand.Read(salt)
	var key = deriveKey(this_ptr.password, salt)

	//read the encdata file
	var path = this_ptr.path + "/EncData"
	content, _ := ioutil.ReadFile(path)

	//encrypt the content
	var encryptor, _ = chacha20poly1305.New(key)
	var encrypted []byte
	encrypted = encryptor.Seal(encrypted, nonce, content, nil)
	safeDelete(path)

	//write out content
	out, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	out.Write(salt)
	out.Write(nonce)
	out.Write(encrypted)

	//close output stream
	out.Close()
}

func (this_ptr *Worker) run() {
	var index = this_ptr.pool.getNewEncDataIndex()
	var encryptor Encryptor
	for index != -1 {
		encryptor = generateChaChaEncryptor(this_ptr.pool.encData[index], this_ptr, index)
		if this_ptr.pool.encrypting {
			encryptor.encrypt(this_ptr.password)
		} else {
			encryptor.decrypt(this_ptr.password)
		}
		index = this_ptr.pool.getNewEncDataIndex()
	}
	this_ptr.wg.Done()
}

func generateWorkerpool(numThreads int, path string, password string) *Workerpool {
	if pathFolderContainsEncData(path) {
		return generateDecryptingWorkerpool(numThreads, path, password)
	} else {
		return generateEncryptingWorkerpool(numThreads, path, password)
	}
}
func (this_ptr *Workerpool) addMacAt(index int, mac []byte) {
	this_ptr.encData[index].mac = mac
}

func pathFolderContainsEncData(path string) bool {
	var buffer []string
	var paths = getAllPaths(path, buffer)
	for _, filePath := range paths {
		split := strings.Split(filePath, "/")
		if strings.Contains(split[len(split)-1], "EncData") {
			return true
		}
	}
	return false
}
