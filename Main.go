package main

import "fmt"

func main() {
	encData := getEncDataArr("/home/robin/Dokumente")
	fmt.Println("===================================================================================")
	for i := 0; i < len(encData); i++ {
		fmt.Println(encData[i].toString())
		fmt.Println("===================================================================================")
	}
}
