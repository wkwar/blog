package logic

import (
	"unicode"
	"unicode/utf8"
)

//输出评论的summary
func TruncateByWords(s string, maxWords int) string {
	processedWords := 0
	wordStarted := false
	for i := 0; i < len(s); {
		//输出 Unicode字符 --- 中文字符 代表一个Unicode字符， 但是占3个字节，使用这个方法，直接输出中文
		r, width := utf8.DecodeRuneInString(s[i:])
		if !isSeparator(r) {
			i += width
			wordStarted = true
			continue
		}

		if !wordStarted {
			i += width
			continue
		}
		//如果为分隔符 -- 一个单词结束，单词数量+1
		wordStarted = false
		processedWords++
		if processedWords == maxWords {
			const ending = "..."
			if (i + len(ending)) >= len(s) {
				// Source string ending is shorter than "..."
				return s
			}

			return s[:i] + ending
		}

		i += width
	}

	// Source string contains less words count than maxWords.
	return s
}

//判断是否为分隔符
func isSeparator(r rune) bool {
	// ASCII 字母数字 and 下划线 不是分隔符
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// 字母和数字不是分割符
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// 空格是分隔符
	return unicode.IsSpace(r)
}
