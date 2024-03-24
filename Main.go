package main

import (
	"fmt"
	"time"
)

func main() {
	var pool = generateWorkerpool(6, "./resources", "password")
	start := time.Now()
	pool.run()
	fmt.Println(time.Since(start))
	fmt.Println("______    _     _             _____                            _")
	fmt.Println("|  ___|  | |   | |           |  ___|                          | |")
	fmt.Println("| |_ ___ | | __| | ___ _ __  | |__ _ __   ___ _ __ _   _ _ __ | |_ ___  _ __")
	fmt.Println("|  _/ _ \\| |/ _` |/ _ \\ '__| |  __| '_ \\ / __| '__| | | | '_ \\| __/ _ \\| '__|")
	fmt.Println("| || (_) | | (_| |  __/ |    | |__| | | | (__| |  | |_| | |_) | || (_) | |")
	fmt.Println("\\_| \\___/|_|\\__,_|\\___|_|    \\____/_| |_|\\___|_|   \\__, | .__/ \\__\\___/|_| ")
	fmt.Println("V2 by Robin K.")
}
