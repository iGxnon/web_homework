package main

import (
	"log"
	"os/exec"
	"time"
)

var SessionsMap = map[string]*Session{}

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
		<-time.Tick(time.Duration(s.maxAge * int64(time.Second)))
		delete(SessionsMap, s.sid)
	}()
}

func PutSessionIfAbsence(session *Session) {
	if _, ok := SessionsMap[session.sid]; !ok {
		SessionsMap[session.sid] = session
		session.onEnable()
	}
}

func PutSession(session *Session) {
	SessionsMap[session.sid] = session
	session.onEnable()
}

// 想不到吧，Windows测试不了的，嘿嘿

func GenerateRandomSid() string {
	// 生成一个uuid作为sid, golang貌似并没有把uuid纳入标准库，投机取巧的获得方法 :p
	// 其实uuid做为键并不好 :<
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
