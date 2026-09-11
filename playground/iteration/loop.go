package iteration

import "strings"

func Repeat(character string, num int) string {
	repeated := strings.Repeat(character, num)
	return repeated
}
