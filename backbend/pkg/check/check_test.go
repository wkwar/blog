package check

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	file := "./a.txt"
	err := Init(file)
	if err != nil {
		t.Error("err", err)
	}
	var str string = "胡紧套,皇冠投注,adhibiu"
	s1, s2 := Trie.Match(str)
	fmt.Println(s1, " ", s2)
}

func TestContains(t *testing.T) {
	var s string = "傻逼"
	var s1 string = "傻&*逼/"
	r, _ := regexp.Compile(s)
	res := r.FindString(s1)
	fmt.Println("Hello, World!", res)
	ok := strings.Contains(s1, s)
   fmt.Println("Hello, World!", ok)
}