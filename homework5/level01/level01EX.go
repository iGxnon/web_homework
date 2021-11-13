package main

import (
	"fmt"
	"strings"
)

// 有些人作业写着写着就...
// 貌似指令系统有点问题，懒的改了(
func main() {
	var input string

	RegisterDefaultPrefix()
	RegisterDefaultSkill()
	fmt.Println("这里是Level01EX技能释放台，请执行你的操作: ")
	fmt.Println("1. 查看技能列表 2. 查看技能模版列表\n3. 释放技能 4. 注册技能\n5. 注册模版 6. 退出")
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}
	for input != "6" {
		switch input {
		case "1":
			fmt.Println("===技能列表===")
			for skillName, _ := range SkillMap {
				fmt.Println(skillName)
			}
			fmt.Println("============")
		case "2":
			fmt.Println("===模版列表===")
			for prefix, _ := range PrefixMap {
				fmt.Println(prefix)
			}
			fmt.Println("============")
		case "3":
			ReleaseStatus()
		case "4":
			RegisterSkillStatus()
		case "5":
			RegisterPrefixStatus()
		case "6":
			fmt.Println("hope u enjoy.")
			break
		}
		_, err := fmt.Scanln(&input)
		if err != nil {
			return
		}
	}
}

func ReleaseStatus() {
	fmt.Println("===你已经进入技能释放界面===")
	var cmd string
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		return
	}
	for cmd != "q" {
		skill, ok := SkillMap[cmd]
		if !ok {
			fmt.Println("没有这个技能")
		} else {
			skill.Release()
		}
		_, err := fmt.Scanln(&cmd)
		if err != nil {
			return
		}
	}
	fmt.Println("==已经退出该界面==")
}

func RegisterSkillStatus() {
	fmt.Println("===你已经进入技能注册界面===")
	fmt.Println("格式:  [skillName]=[prefix]")
	var cmd string
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		return
	}
	for cmd != "q" {
		split := strings.Split(cmd, "=")
		if len(split) != 2 {
			fmt.Println("注册格式错误！")
			break
		}
		err1 := RegisterSkill(split[0], split[1])
		if err1 != nil {
			fmt.Println(err1.Error())
			break
		}
		fmt.Println("注册成功！")
		_, err2 := fmt.Scanln(&cmd)
		if err2 != nil {
			return
		}
	}
	fmt.Println("==已经退出该界面==")
}

func RegisterPrefixStatus() {
	fmt.Println("===你已经进入技能模版注册界面===")
	fmt.Println("格式1: prefix=[{skill}表示技能名称]")
	fmt.Println("格式2: prefix=[内容]")
	var cmd string
	_, err := fmt.Scanln(&cmd)
	if err != nil {
		return
	}
	for cmd != "q" {
		split := strings.Split(cmd, "=")
		if len(split) != 2 {
			fmt.Println("注册格式错误！")
			break
		}
		if strings.Contains(split[1], "{skill}") {
			err := RegisterPrefixWithFormat(split[0], split[1])
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		} else {
			err := RegisterPrefixWithAllArg(split[0], split[1])
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
		_, err := fmt.Scanln(&cmd)
		if err != nil {
			return
		}
	}
	fmt.Println("==已经退出该界面==")
}

func ReleaseSkill(skillNames string, releaseSkillFunc func(string)) {
	releaseSkillFunc(skillNames)
}
