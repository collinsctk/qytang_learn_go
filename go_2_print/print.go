package main

import "fmt"

func main() {
	// -------------- Print不换行 ----------------
	fmt.Print("My weight on the surface of Mars is ")
	// -------------- Println换行 ----------------
	fmt.Println(149 * 0.3783)
	fmt.Println("My weight on the surface of Mars is", 149*0.3783)
	// ---------------下面是格式化-----------------
	fmt.Printf("My weight on the surface of Mars is %v", 149*0.3783)
	fmt.Printf(" and I would be %v years old.\n", 41*365/687)
	fmt.Printf("My weight on the surface of %v is %v \n", "Earth", 149*0.3783)

	// ---------------负号表示左对齐, 15表示占15个字符的宽度, 4表示占4个字符的宽度-----------------
	fmt.Printf("%-15v $%4v\n", "SpaceX", 94)
	fmt.Printf("%-15v $%4v\n", "Virgin Galactic", 100)
}
