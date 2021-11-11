package homework4

import "fmt"

var x int64
var flag = 0

/**
使用 channel 的阻塞作用来代替锁进行同步
*/

func add(identifier string, ch chan int) {
	<-ch
	for i := 0; i < 5000; i++ {
		x = x + 1
		fmt.Println(identifier)
	}
	flag++
	ch <- flag
}

//两个goroutine同时抢一个channel,先抢到的先运行完毕再往channel里写数据让另一个接收去运行

func Level01() {

	ch := make(chan int)
	go add("goroutine-1 add 1 to x", ch)
	go add("goroutine-2 add 1 to x", ch)
	ch <- flag
	// 循环同步
	for <-ch != 2 {

	}
	fmt.Println(x)

}
