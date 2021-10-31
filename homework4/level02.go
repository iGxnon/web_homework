package main

import "os"

func main() {
	file, err := os.Create("./txt/plan.txt")
	if err != nil {
		return
	}
	_, err2 := file.Write([]byte("I’m not afraid of difficulties and insist on learning programming"))
	if err2 != nil {
		return
	}

}
