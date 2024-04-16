package logic

import (
	"backbend/pkg/check"
) 


func Check(text string) (sensitiveWords []string, replaceText string) {
	return check.Trie.Match(text)
}