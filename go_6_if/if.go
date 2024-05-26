package main

import "fmt"

func main() {
	// ---------------- 一个简单的if语句 ----------------
	var command = "go east"

	if command == "go east" {
		println("You head further up the mountain.")
	} else if command == "go inside" {
		println("You enter the cave where you live out the rest of your life.")
	} else {
		println("Didn't quite get that.")
	}

	// ---------------- and 和 or ----------------
	fmt.Println("The years is 2100, should you leap? ")
	var year = 2100
	var leap = year%400 == 0 || (year%4 == 0 && year%100 != 0)

	if leap {
		fmt.Println("Look before you leap!")
	} else {
		fmt.Println("Keep your feet on the ground.")
	}

	//---------------------! 取反---------------------
	var haveTorch = true
	var litTorch = false

	if !haveTorch || !litTorch {
		fmt.Println("Nothing to see here.")
	}

	//---------------------switch---------------------
	fmt.Println("There is a cavern entrance here and a path to the east.")

	var command2 = "go inside"
	switch command2 {
	case "go east":
		fmt.Println("You head further up the mountain.")
	case "enter cave", "go inside":
		fmt.Println("You find yourself in a dimly lit cavern.")
	case "read sign":
		fmt.Println("The sign reads 'No Minors'.")
	default:
		fmt.Println("Didn't quite get that.")
	}
	//---------------------switch fallthrough---------------------
	var room = "lake"
	switch {
	case room == "cave":
		fmt.Println("You find yourself in a dimly lit cavern.")
	case room == "lake":
		fmt.Println("The ice seems solid enough.")
		fallthrough // fallthrough关键字会继续执行下一个case
	case room == "underwater":
		fmt.Println("The water is freezing cold.")
	}
}
