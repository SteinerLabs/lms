package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Billing service")
	for {
		time.Sleep(10 * time.Second)
	}
}
