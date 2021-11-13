package main

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
	for _, badword := range m {
		if strings.Contains(s, badword) {
			return true
		}
	}
	return false
}
