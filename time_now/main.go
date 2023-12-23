package main

import (
	"fmt"
	"time"
	_ "unsafe"
)

func main() {
	fmt.Println("Current Time", time.Now().UTC())
	fmt.Println("Is Valid?", isValid())
}
