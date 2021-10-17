package main

import "fmt"

func main() {
	var x, y float32
	var op string
	for {
		fmt.Println("输入表达式")
		_, err := fmt.Scanln(&x, &op, &y)
		if err != nil {
			return
		}
		switch op {
			case "+":
				fmt.Println(x + y)
			case "-":
				fmt.Println(x - y)
			case "*":
				fmt.Println(x * y)
			case "/":
				fmt.Println(x / y)
			default:
				panic("你输入的你马呢?")
			}
	}
}
