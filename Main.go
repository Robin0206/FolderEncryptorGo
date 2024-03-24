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
}
