package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

func main() {
	randomList := make([]int, 0)
	for i := 0; i < 100; i ++ {
		randomList = append(randomList, int(rand.Int63n(100))) //伪随机数? 每次都一样。。装起来了
	}
	randomList = sortSleep(randomList)
	fmt.Println(randomList)
}

// Bubble is me
func sortBubble(randList []int) []int {
	decline := 0 // 提升效率
	for range randList {
		for j := 0; j < len(randList) - decline - 1; j ++ {
			maxNum := int(math.Max(float64(randList[j]), float64(randList[j + 1])))
			minNum := int(math.Min(float64(randList[j]), float64(randList[j + 1])))
			randList[j] = minNum
			randList[j + 1] = maxNum
		}
		decline ++
	}
	return randList
}

// 艹，卷个p，睡起来了
// 睡排
func sortSleep(randList []int) (result []int)  {
	result = make([]int, 0)
	var wg sync.WaitGroup
	wg.Add(len(randList))
	for _, v := range randList {
		go func(t int) {
			time.Sleep(time.Duration(t) * time.Millisecond)
			result = append(result, t)
			wg.Done()
		}(v)
	}
	wg.Wait()
	return
}