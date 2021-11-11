package level01

import (
	"errors"
	"fmt"
	"strings"
)

type ReleaseFunc func(string)

var PrefixMap = make(map[string]ReleaseFunc)

func RegisterDefaultPrefix() {
	PrefixMap["全集中·常中"] = func(skillName string) {
		fmt.Println("全集中·常中", skillName)
	}
}

func RegisterPrefix(prefix string) (err error) {
	return RegisterPrefixWithAllArg(prefix, prefix)
}

func RegisterPrefixIfAbsence(prefix, output string) (err error) {
	if _, ok := PrefixMap[prefix]; !ok {
		err = RegisterPrefixWithAllArg(prefix, output)
	}
	return
}

func ChangePrefix(prefix, output string) (err error) {
	err = nil
	if _, ok := PrefixMap[prefix]; ok {
		err = RegisterPrefixWithAllArg(prefix, output)
	} else {
		err = errors.New("列表里好像没有这个技能模版哦")
	}
	return
}

func RegisterPrefixWithAllArg(prefix, output string) (err error) {
	err = nil
	if check(output) {
		err = errors.New("你嘴臭你马呢？")
		return
	}
	PrefixMap[prefix] = func(skillName string) {
		fmt.Println(output, skillName)
	}
	return
}

// RegisterPrefixWithFormat 垃圾代码
func RegisterPrefixWithFormat(prefix, format string) (err error) {
	err = nil
	if check(format) {
		err = errors.New("你嘴臭你马呢？")
		return
	}
	PrefixMap[prefix] = func(skillName string) {
		fmt.Println(strings.Replace(format, "{skill}", skillName, -1))
	}
	return
}
