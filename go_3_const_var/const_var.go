package main

import "fmt"

// const声明常量
// const声明的常量是不可变的

// var声明变量
// var声明的变量是可变的

func main() {
	const lightSpeed = 299792 // km/s
	var distance = 56000000   // km

	fmt.Println(distance/lightSpeed, "seconds")
	distance = 401000000 // 变量可以重新赋值
	fmt.Println(distance/lightSpeed, "seconds")

	var speed = 100
	fmt.Println(speed)

	distance, speed = 56000000, 100.0

	//speed = speed * 0.5  这个会错误，因为speed是int类型，0.5是float类型
	//在 Go 语言中，乘法运算的两个操作数必须是相同的类型。也就是说，如果一个操作数是整数，那么另一个操作数也必须是整数。同样，
	//如果一个操作数是浮点数，那么另一个操作数也必须是浮点数。

	speed = speed * 50
	speed *= 50
	speed++ // speed = speed + 1
	fmt.Println(speed)
	fmt.Println(distance/speed, "seconds")
}
