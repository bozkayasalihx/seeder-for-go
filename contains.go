package main

import (
	"strings"
)

func IsContains(str string, substr string) bool {
	isContains := strings.Contains(str, substr)
	return isContains
}
