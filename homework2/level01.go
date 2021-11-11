package homework2

import (
	"fmt"
	"strings"
)

func Level01() {
	var str string
	_, err := fmt.Scanln(&str)
	if err != nil {
		return
	}
	strs := strings.Split(str, "")
	newStrs := ""
	for i := len(strs) - 1; i >= 0; i-- {
		newStrs += strs[i]
	}
	fmt.Println(newStrs)
}
