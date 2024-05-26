package main

import (
	"fmt"
	"strings"
)

// --------------------主要介绍bool类型的判断--------------------
func main() {
	fmt.Println("You find yourself in a dimly lit cavern.")
	// --------------------字符串包含的判断--------------------
	var command = "walk outside"
	var exit = strings.Contains(command, "outside") // Contains函数用于判断字符串是否包含另一个字符串

	fmt.Println("You leave the cave:", exit)

	fmt.Println("There is a sign near the entrance that reads 'No Minors'.")

	// --------------------小于号的判断--------------------
	var age = 41
	var minor = age < 18

	fmt.Println("At age", age, "are you a minor?", minor)
}
