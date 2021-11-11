package level02

import (
	"fmt"
	"time"
)

const (
	Continue = iota
	Once
)

type Clock struct {
	Duration   time.Duration
	Mess       string
	ConType    int
	cancelNext bool
	tick       int64
	cancelled  bool
	ch         chan int
}

func NewClock(duration time.Duration, mess string, conType int, ch chan int) *Clock {
	clock := &Clock{
		Duration:   duration,
		Mess:       mess,
		ConType:    conType,
		cancelNext: false,
		tick:       0,
		cancelled:  false,
		ch:         ch,
	}
	go func(c *Clock) {
		tick := time.Tick(time.Second)
		for {
			<-tick
			if c.cancelled {
				break
			}
			c.run()
		}
	}(clock)
	return clock
}

// 关闭这个闹钟
func (c *Clock) CancelAll() {
	c.cancelled = true
	c.ch <- 1
}

// 闹钟
func (c *Clock) Noise() {
	if c.cancelNext {
		return
	}
	fmt.Println(c.Mess)
}

// 关闭下一次
func (c *Clock) CancelNext() {
	c.cancelNext = true
	// recover
	fmt.Println("懒狗，已经帮你关闭一次闹钟了！")
	go func() {
		<-time.Tick(c.Duration)
		<-time.Tick(time.Millisecond * 10) // 延迟一会防止误判
		c.cancelNext = false
	}()
}

// 不对外暴露
func (c *Clock) run() {
	c.tick++
	if c.tick == int64(c.Duration.Seconds()) {
		c.Noise()
		c.tick = 0
		if c.ConType == Once {
			c.CancelAll()
		}
	}
}

func Level02EX() {
	ch := make(chan int)
	clock := NewClock(time.Second*2, "懒狗，懒狗起床了！", Continue, ch)
	go func(c *Clock) {
		time.Sleep(time.Second * 6)
		c.CancelNext()
		time.Sleep(time.Second * 6)
		c.CancelAll()
	}(clock)
	<-ch
}
