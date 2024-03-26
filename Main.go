package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	var path, _ = filepath.Abs(os.Args[1])
	var numThreads = runtime.NumCPU() / 2
	showBanner()

	for {
		fmt.Println("======================info======================")
		printPathInfo(path)
		printThreadInfo(numThreads)
		printNumOfFiles(path)
		printBufferSizeInfo()
		choice := showMenuAndGetChoice(path)
		if choice == 3 {
			password := getPasswordFromUser()
			var pool = generateWorkerpool(numThreads, path, password)
			var start = time.Now()
			fmt.Println("Running workers")
			pool.run()
			fmt.Println("Done after " + time.Since(start).String())
			break
		}
		if choice == 2 {
			numThreads = getNumOfThreadsFromUser()
		}
		if choice == 1 {
			ENC_BUFFERSIZE = getBufferSizeFromUser()
		}
		if choice != 1 && choice != 2 && choice != 3 {
			os.Exit(0)
		}
	}

}

func getBufferSizeFromUser() int {
	fmt.Print("Please set the desired buffer size in bytes> ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	line = strings.ReplaceAll(line, "\n", "")
	if err != nil {
		log.Fatal(err)
	}
	num, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func printBufferSizeInfo() {
	fmt.Println("Buffersize: " + strconv.Itoa(ENC_BUFFERSIZE) + " bytes")
}

func printNumOfFiles(path string) {
	var paths []string
	paths = getAllPaths(path, paths)
	fmt.Print("Number of files: ")
	fmt.Print(len(paths))
	fmt.Println()
}

func getNumOfThreadsFromUser() int {
	fmt.Print("Please set the desired number of threads> ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	line = strings.ReplaceAll(line, "\n", "")
	if err != nil {
		log.Fatal(err)
	}
	num, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func getPasswordFromUser() string {
	fmt.Print("Please set the desired password> ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	return line
}

func printThreadInfo(numThreads int) {
	fmt.Print("Number of Threads: ")
	fmt.Print(numThreads)
	fmt.Println()
}
func showMenuAndGetChoice(path string) int {
	fmt.Println("======================Menu======================")
	if pathFolderContainsEncData(path) {
		fmt.Println("3 => Decrypt")
	} else {
		fmt.Println("3 => Encrypt")
	}
	fmt.Println("2 => Change number of threads")
	fmt.Println("1 => Change buffer size")
	fmt.Println("0 => Exit")
	fmt.Println()
	fmt.Print("Please set the desired option> ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	line = strings.ReplaceAll(line, "\n", "")
	if err != nil {
		log.Fatal(err)
	}
	num, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func printPathInfo(path string) {
	if pathFolderContainsEncData(path) {
		fmt.Print("Folder to decrypt: ")
		if len(path) > 20 {
			fmt.Println("...." + path[len(path)-20:])
		} else {
			fmt.Println(path)
		}
	} else {
		fmt.Print("Folder to encrypt: ")
		if len(path) > 20 {
			fmt.Println("...." + path[len(path)-20:])
		} else {
			fmt.Println(path)
		}
	}
}

func showBanner() {
	fmt.Println("______    _     _             _____                            _")
	fmt.Println("|  ___|  | |   | |           |  ___|                          | |")
	fmt.Println("| |_ ___ | | __| | ___ _ __  | |__ _ __   ___ _ __ _   _ _ __ | |_ ___  _ __")
	fmt.Println("|  _/ _ \\| |/ _` |/ _ \\ '__| |  __| '_ \\ / __| '__| | | | '_ \\| __/ _ \\| '__|")
	fmt.Println("| || (_) | | (_| |  __/ |    | |__| | | | (__| |  | |_| | |_) | || (_) | |")
	fmt.Println("\\_| \\___/|_|\\__,_|\\___|_|    \\____/_| |_|\\___|_|   \\__, | .__/ \\__\\___/|_| ")
	fmt.Println("V2 by Robin K.                                      |__/|_| ")
}
