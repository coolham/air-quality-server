package utils

import (
	"strconv"
)

// StringToInt 将字符串转换为整数
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// IntToString 将整数转换为字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// BoolToString 将布尔值转换为字符串
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// StringToBool 将字符串转换为布尔值
func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// IntPtr 返回整数指针
func IntPtr(i int) *int {
	return &i
}

// BoolPtr 返回布尔值指针
func BoolPtr(b bool) *bool {
	return &b
}
