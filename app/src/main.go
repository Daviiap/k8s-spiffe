package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Running...")
		time.Sleep(2 * time.Second)
	}
}
