package main

import "fmt"

func main() {
	var paths []string
	paths = getAllPaths("/home/robin/GolandProjects", paths)
	fmt.Println("===================================================================================")
	for i := 0; i < len(paths); i++ {
		fmt.Println(paths[i])
	}
}
