package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Notification service")
	for {
		time.Sleep(10 * time.Second)
	}
}
