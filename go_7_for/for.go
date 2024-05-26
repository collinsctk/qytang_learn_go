package main

import (
	"fmt"
	"time"
)

func main() {
	var count1 = 10

	for count1 > 0 {
		fmt.Println(count1)
		time.Sleep(time.Second)
		count1--
	}
	fmt.Println("Liftoff1!")

	var count2 = 10
	for { //不给条件，就是无限循环
		if count2 < 0 {
			break //满足条件就跳出循环
		}
		fmt.Println(count2)
		time.Sleep(time.Second * 2)
		count2--
	}
	fmt.Println("Liftoff2!")
}
