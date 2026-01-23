package main

import (
	"fmt"
	"strings"
	"sync"
)

func printData(wg *sync.WaitGroup, data []byte) {
	defer wg.Done()

	var result strings.Builder
	for _, b := range data {
		result.WriteString(string(b))
	}

	fmt.Println(result.String())
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	data1 := "hello"
	data2 := "hello world"

	go printData(&wg, []byte(data1))
	go printData(&wg, []byte(data2))

	wg.Wait()
}
