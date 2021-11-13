package main

import (
	"fmt"
	"strings"
	"time"
)

func Level02() {
	sync := time.Tick(time.Second)
	// 同步
	for {
		t := <-sync
		if strings.Contains(t.String(), "00:00:00") {
			break
		}
	}
	hourCnt := 0
	minCnt := 0
	for {
		<-sync
		if hourCnt == 2 {
			fmt.Println("谁能比我还卷!")
		}
		if hourCnt == 6 {
			fmt.Println("早八算什么，早六才是吾辈应起之时")
		}
		if minCnt == 30 {
			fmt.Println("芜湖！起飞！")
		}
		minCnt++
		if minCnt == 60 {
			minCnt = 0
			hourCnt++
		}
		if hourCnt == 24 {
			hourCnt = 0
		}
	}
}
