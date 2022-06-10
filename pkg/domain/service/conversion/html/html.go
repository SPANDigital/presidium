package html

import (
	"strings"
)

func ContainsHTML(content string) bool {
	return strings.Index(content, "<") > -1 && strings.Index(content, ">") > -1
}
