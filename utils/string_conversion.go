package utils

import "strconv"

func StrToUnint(str string) uint {
	parsedUint, _ := strconv.ParseUint(str, 10, 32)
	return uint(parsedUint)
}
