package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var num = rand.Intn(10) + 1 // rand.Intn(10)包括10吗？不包括
	fmt.Println(num)

	num = rand.Intn(10) + 1
	fmt.Println(num)
}
