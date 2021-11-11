package homework3

import "fmt"

func Level03() {
	Receiver(true)
	Receiver("Hi")
	Receiver(19)
}

// 所有类型都可以看作是空接口的实现
func Receiver(v interface{}) {
	switch v.(type) {
	case int:
		fmt.Println("这是一个int")
	case string:
		fmt.Println("这是一个string")
	case bool:
		fmt.Println("这是一个bool")
	}
}
