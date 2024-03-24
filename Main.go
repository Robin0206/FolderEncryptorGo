package main

func main() {
	var pool = generateWorkerpool(6, "./resources/RG-main", "password")
	pool.run()
}
