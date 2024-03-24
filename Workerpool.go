package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type Worker struct {
	pool      *Workerpool
	encryptor Encryptor
	password  string
	wg        *sync.WaitGroup
}

type Workerpool struct {
	workers            []Worker
	encData            []EncData
	alreadyEnDeCrypted []bool
	mutex              sync.Mutex
	wg                 sync.WaitGroup
	path               string
	encrypting         bool
}

func generateEncryptingWorkerpool(numThreads int, path string, password string) *Workerpool {
	var result Workerpool
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
		result.workers[i] = generateTestWorker(&result, password)
		result.wg.Add(1)
	}
	return &result
}

func generateDecryptingWorkerpool(numThreads int, path string, password string) *Workerpool {
	var result Workerpool
	result.encrypting = false
	result.path = path
	result.encData = getEncDataArrFromFile(path)
	result.alreadyEnDeCrypted = make([]bool, len(result.encData))
	for i := 0; i < len(result.alreadyEnDeCrypted); i++ {
		result.alreadyEnDeCrypted[i] = false
	}
	result.mutex = sync.Mutex{}
	result.workers = make([]Worker, numThreads)
	for i := 0; i < numThreads; i++ {
		result.workers[i] = generateTestWorker(&result, password)
		result.wg.Add(1)
	}
	return &result
}

func getEncDataArrFromFile(path string) []EncData {
	var result []EncData
	file, _ := os.Open(path + "/EncData")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, encDataFromString(scanner.Text()))
	}
	return result
}

func generateTestWorker(w *Workerpool, password string) Worker {
	var result Worker
	result.pool = w
	result.password = password
	result.wg = &(w.wg)
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
}

func (this_ptr *Worker) run() {
	var index = this_ptr.pool.getNewEncDataIndex()
	var encryptor Encryptor
	for index != -1 {
		encryptor = generateChaChaEncryptor(this_ptr.pool.encData[index])
		encryptor.encrypt(this_ptr.password)
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
