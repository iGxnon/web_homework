package homework4

import "fmt"

func Level03() {
	// 9个缓存区，避免没必要的阻塞
	over := make(chan bool, 9)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
			over <- true
		}(i)
		// 应该在最后一个goroutine里往channel里写入true
		// 或者往goroutine里写入九次，并在主协程里读取九次
		//if i == 9 {
		//over <- true
		//}
	}
	// 接收9次，少了最后一次的话就阻塞
	for i := 0; i < 10; i++ {
		<-over
	}
	fmt.Println("over!!!")
}
