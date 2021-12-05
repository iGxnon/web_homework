package main

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var SessionsMap = map[string]*Session{}

// 互斥锁，用于保证SessionsMap操作的原子性
var mutex sync.Mutex

type Session struct {
	sid    string
	maxAge int64
	id     string
}

func (s *Session) onEnable() {
	// 持久的不用管
	if s.maxAge <= 0 {
		return
	}
	// 不能在访问量大的时候这么操作,只是跑跑测试
	go func() {
		time.Sleep(time.Duration(s.maxAge * int64(time.Second)))
		delete(SessionsMap, s.sid)
	}()
}

func PutSessionIfAbsence(session *Session) {
	mutex.Lock()
	if _, ok := SessionsMap[session.sid]; !ok {
		SessionsMap[session.sid] = session
		session.onEnable()
	}
	mutex.Unlock()
}

func PutSession(session *Session) {
	mutex.Lock()
	SessionsMap[session.sid] = session
	session.onEnable()
	mutex.Unlock()
}

// GenerateRandomSid 基于时间生成的uuid
func GenerateRandomSid() string {

	rand.Seed(666)
	one := strconv.FormatInt(int64(rand.Intn(100)*time.Now().Nanosecond()%0xFFFFFFFF), 16)
	two := strconv.FormatInt(int64(rand.Intn(100)*time.Now().Nanosecond()%0xFFFF), 16)
	three := strconv.FormatInt(int64(rand.Intn(100)*time.Now().Nanosecond()%0xFFFF), 16)
	four := strconv.FormatInt(int64(rand.Intn(100)*time.Now().Nanosecond()%0xFFFF), 16)
	five := strconv.FormatInt(int64(rand.Intn(100)*time.Now().Nanosecond()%0xFFFFFFFFFFFF), 16)

	return one + "-" + two + "-" + three + "-" + four + "-" + five
}
