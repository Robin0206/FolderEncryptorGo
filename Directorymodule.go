package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getAllPaths(dir string, buffer []string) []string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println(err.Error())
		}
		if !info.IsDir() {
			buffer = append(buffer, absPath)
		}
		return nil
	})
	if err == nil {

	} else {
		fmt.Println(err)
	}
	return buffer
}
func getEncDataArr(path string) []EncData {
	var paths []string
	paths = getAllPaths(path, paths)
	var result []EncData
	for i := 0; i < len(paths); i++ {
		result = append(result, generateEncData(paths[i]))
	}
	return result
}
