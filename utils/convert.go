package utils

import "strconv"

// ConvertStringToInt 字符串转int
func ConvertAStringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}
