package main

import (
	"encoding/base64"
	"fmt"
)

const (
	// PublicKey保存在客户端
	PUBLIC_KEY = "MIGJAoGBAL5z+JAG0XbCXe90ASsrsBUsTo03pMzWUe8cteKZJslJjvBxWvtlcyMNklXQgg5N0dkzUrfbt1y6hOykQbTYrRLXwbqUy9FBQT2lsU1jHcPCT/vc3vWU3WHrZPhd09KAFBhZxtHmEw9uW0wjsgqu1au8yhIq3a2M+JKII/kaW+flAgMBAAE="
	// PrivateKey保存在服务端
	PRIVATE_KEY = "MIICXgIBAAKBgQC+c/iQBtF2wl3vdAErK7AVLE6NN6TM1lHvHLXimSbJSY7wcVr7ZXMjDZJV0IIOTdHZM1K327dcuoTspEG02K0S18G6lMvRQUE9pbFNYx3Dwk/73N71lN1h62T4XdPSgBQYWcbR5hMPbltMI7IKrtWrvMoSKt2tjPiSiCP5Glvn5QIDAQABAoGBALbW1UlIEm3F+bJ5lumgHoKlL6BpTCCOnMhGsuMhDthtcvmoiaUR5zA+xj72VvVuhkjT+dSi7ezq79PTeXUqEzR4zo9qJficCTrgjhJ6i0B7IlWweCz2tKv3jQbSianDNhaYznZzva3s/NRwyP0qVUqIQ540GIdCbxynHDHF3iahAkEA/VNM3NiyKWW9ems7uD8j4HLviXVvgl64S4BNyfG1wyes6E34jf/qzVu2t2sa1Y2Yx77kOOZdc9u41v5gQ2VsLQJBAMB2vH0YoyO/TL9sxp+DuezpQWN24Ue2eOlPxwAYsuIa9LkCqPTX0xtSP7Knf8ZkXpEBmXhHrONwXxBiD4ci5ZkCQG/VSmVktKJZ6+ATXvXjye7YTq8cTPH85tdN+Qlhz6Ar78VORqBJjlrCVlN60QndzMjBmPcVm8P+CAfBnLWkHLECQQCj/ghJh06qzPv2OBdeH/2ycmY2/DqkwkRweHuWB3WU12cipbOVPLkylHiWH8buIuO5JuW/6ULVYRB/gy679O4xAkEAu8+slpuekvmZyFQkFGLU0OdlBEi5x9Yq+WyJ8H9KMIy1M2DINhMrTObLk1lHD2Ju0hn6A0J4Ie/c7Qh7OYyAXg=="
)

type LoginPage struct {
	Name_    string
	PwdLock_ string
}

// 以下代码会在客户端运行

func (p *LoginPage) Login_1() {
	fmt.Println("请输入名称: ")
	var name string
	_, err := fmt.Scanln(&name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("请输入密码: ")
	var pwd string
	_, err2 := fmt.Scanln(&pwd)
	if err2 != nil {
		fmt.Println(err2.Error())
		return
	}

	p.Name_ = name
	encrypt, err := RSA_Encrypt([]byte(pwd), PUBLIC_KEY)
	if err != nil {
		fmt.Println("加密时出现了错误!")
		return
	}
	p.PwdLock_ = base64.StdEncoding.EncodeToString(encrypt)
}

// 以下程序在服务端执行

func (p *LoginPage) Login_2(pool DataPool) {
	if pool.ContainsName(p.Name_) {
		dao := pool.getIfContainsName(p.Name_)
		savePwd, _ := RSA_Decrypt([]byte(dao.Pwd), PRIVATE_KEY)
		inputPwd, _ := RSA_Decrypt([]byte(p.PwdLock_), PRIVATE_KEY)
		if string(savePwd) == string(inputPwd) {
			fmt.Println("成功登录!")
			//todo 后台命令系统：密码更改，访问权限
		} else {
			fmt.Println("密码错误! ")
		}
	} else {
		fmt.Println("未查明此用户，重新查找中...")
		if pool.ReloadAll(); pool.ContainsName(p.Name_) {
			fmt.Println("查找成功！")
			p.Login_2(pool)
		} else {
			fmt.Println("未查到此用户，是否注册？y/N")
			var cmd string
			fmt.Scanln(&cmd)
			if cmd == "y" {
				if !checkPwd(p.PwdLock_) {
					fmt.Println("密码不规范!")
					p.Login_1()
					return
				}
				pool.Cache = append(pool.Cache, UserDao{
					Name: p.Name_,
					Pwd:  p.PwdLock_,
				})
				pool.saveToDisk()
				fmt.Println("注册成功!")
			}
		}
	}
}

func checkPwd(cipherText string) bool {
	decodeString, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return false
	}
	decrypt, err := RSA_Decrypt(decodeString, PRIVATE_KEY)
	if err != nil {
		return false
	}
	pwdRaw := string(decrypt)
	if len(pwdRaw) < 6 {
		return false
	}

	return true
}
