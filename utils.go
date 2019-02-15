package main

import (
	"fmt"
)

func Arr2RegOr(arr []string, prefix string, suffix string) string {
	if arr == nil || len(arr) == 0 {
		return ""
	}
	var ret = prefix + "("
	var start = true
	for _, s := range arr {
		if start {
			ret += s
		} else {
			start = false
			ret += "|" + s
		}
	}
	return fmt.Sprintf("(%s)%s)", ret, suffix)
}
