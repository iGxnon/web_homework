package homework3

import (
	"fmt"
	"time"
)

type Person struct {
	name      string
	age       int
	isMarried bool
	id        int
}

type Page struct {
	Title string
	Bv    string
	Date  time.Time
	Video Video
}

type Video struct {
	Time int
}

type 一键三连 interface {
	大拇指()
	大钢镚()
	小星星()
}

func (v Video) 大拇指() {
	panic("大拇指")
}

func (v Video) 大钢镚() {
	fmt.Println("大钢镚")
}

func (v Video) 小星星() {
	fmt.Println("小星星")
}

func Level01AndLevel02() {
	xiaoMin := Person{
		name:      "小明",
		age:       5,
		isMarried: false,
		id:        2021211503,
	}

	bvideo := Page{
		Title: "NMSL",
		Bv:    "Bv114514",
		Date:  time.Now(),
		Video: Video{
			Time: 114514,
		},
	}
	fmt.Println(xiaoMin)
	fmt.Println(bvideo)
	bvideo.Video.大钢镚()
}
