package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("run local...")
	for {
		fmt.Println("run1")
		time.Sleep(5 * time.Second)
	}
}
