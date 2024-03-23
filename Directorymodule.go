package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getAllPaths(path string, pathsBuffer []string) {
	var localPaths = getAllLocalPaths(path)
	for _, p := range localPaths {
		pathsBuffer = append(pathsBuffer, p)
	}
	for _, p := range localPaths {
		if isDirectory(p) {
			getAllPaths(p, pathsBuffer)
		}
	}
}

func isDirectory(path string) bool {
	file, err := os.Open(path)
	var result bool
	if err != nil {
		fmt.Println("Path doesnt exist:")
		fmt.Println(path)
		result = false
	} else {
		fileInfo, _ := file.Stat()
		result = fileInfo.IsDir()
	}
	file.Close()
	return result
}

func getAllLocalPaths(dir string) []string {
	var result []string
	files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	for _, file := range files {
		result = append(result, filepath.Join(path, file.Name()))
	}

	return result
}
