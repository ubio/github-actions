package main

import (
	"fmt"
	"os"
)

func main() {
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}
}
