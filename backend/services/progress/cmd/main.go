package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Progress service")
	for {
		time.Sleep(10 * time.Second)
	}
}
