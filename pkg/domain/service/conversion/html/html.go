package html

import (
	"strings"
)

func ContainsHTML(content string) bool {
	return strings.Contains(content, "<") && strings.Contains(content, ">")
}
