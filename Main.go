package main

import "fmt"

func main() {
	var paths []string
	getAllPaths("/home/robin/programme/", paths)
	fmt.Println("===================================================================================")
	for i := 0; i < len(paths); i++ {
		fmt.Println(paths[i])
	}
}
