package check

import (
	"strings"
	"io/ioutil"
)

var (
	Trie *SensitiveTrie
)

func Init(file string) (err error) {
	//敏感词
	//
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	
	//
	sensitiveWords := strings.Split(string(data), "，")
	Trie = NewSensitiveTrie()
	Trie.AddWords(sensitiveWords)
	return nil
}