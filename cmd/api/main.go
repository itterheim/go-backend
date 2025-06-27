package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("Done in", elapsed.String())
	}()

	fmt.Println("API")
}
