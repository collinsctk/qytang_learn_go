package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//var count = 0 // 变量的作用域是在{}内部
	count := 0 // 短变量声明

	for count < 10 {
		var num = rand.Intn(10) + 1
		fmt.Println(num)

		count++
	}
	fmt.Println("----------")

	for qyt_count := 0; qyt_count < 10; qyt_count++ {
		var num = rand.Intn(10) + 1
		fmt.Println(num)
	}

	fmt.Println("----------")

	switch num := rand.Intn(10); num {
	case 0:
		fmt.Println("Space Adventures")
	case 1:
		fmt.Println("SpaceX")
	case 2:
		fmt.Println("Virgin Galactic")
	default:
		fmt.Println("Random spaceline #", num)
	}
}
