package utils

import (
	"strings"
)

/**
* 驼峰转蛇形
 */
func SnakeString(s string) string {
	var b strings.Builder
	num := len(s)

	for i := 0; i < num; i++ {
		cur := s[i]

		if cur >= 'A' && cur <= 'Z' {
			cur = cur + 32
			if i != 0 {
				b.WriteByte('_')
			}
		}

		b.WriteByte(cur)
	}

	return b.String()
}

/**
* 蛇形转驼峰
 */
func CamelString(s string) string {
	var b strings.Builder
	num := len(s)

	cur := s[0]
	//first
	if cur != '_' {
		if cur >= 'a' && cur <= 'z' {
			cur = cur - 32
		}
		b.WriteByte(cur)
	}

	for i := 1; i < num; i++ {
		cur = s[i]
		if cur == '_' {
			i++
			if i < num {
				cur = s[i]
				if cur >= 'a' && cur <= 'z' {
					cur = cur - 32
				}
				b.WriteByte(cur)
			}
		}else {
			b.WriteByte(cur)
		}
	}

	return b.String()
}
