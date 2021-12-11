package utils

import (
	"fmt"
	"testing"
)

func TestCheckPwdSafe(t *testing.T) {
	fmt.Println(CheckPwdSafe("1234abDc.."))
}
