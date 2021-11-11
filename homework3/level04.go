package homework3

import (
	"fmt"
	"sort"
)

func Level04() {
	a := []int{20, 1, 45, 123, 3, -30, 45, 31}

	// 函数式编程
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	fmt.Println(a)
}
