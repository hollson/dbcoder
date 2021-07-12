package utils

import (
	// _ "github.com/lib/pq"
	"strings"
)

const BLANK = ""

// 目标元素tar是否包含在container集合中
func ContainString(tar string, container ...string) bool {
	for _, v := range container {
		if tar == v {
			return true
		}
	}
	return false
}

// 转换为帕斯卡命名
//  如: userName => UserName
//     user_name => UserName
func Pascal(title string) string {
	arr := strings.FieldsFunc(title, func(c rune) bool { return c == '_' })
	RangeStringsFunc(arr, func(s string) string { return strings.Title(s) })
	return strings.Join(arr, BLANK)
}

// 遍历处理集合成员
func RangeStringsFunc(arr []string, f func(string) string) {
	for k, v := range arr {
		arr[k] = f(v)
	}
}

func PathTrim(path string) string {
	return strings.ReplaceAll(path, "//", "/")
}
