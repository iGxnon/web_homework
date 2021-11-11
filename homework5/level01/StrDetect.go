package level01

import "strings"

func check(s string) bool {
	m := []string{
		"nmsl",
		"猎马",
		"马",
		"妈",
		"鼠",
		"nigga",
		"nigger",
	}
	return strings.Contains(strings.Join(m, "-"), s)
}
