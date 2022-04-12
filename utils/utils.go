package utils

import (
	"fmt"
	"regexp"
	"strings"
)

const BLANK = ""

// ContainString 目标元素tar是否包含在container集合中
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

// 字符串切片去重
func DistinctStringArray(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

// 修剪路径中的切割符
func PathTrim(path string) string {
	return strings.ReplaceAll(path, "//", "/")
}

// raw中是否包含sub...中的任意子串
//  如: HasSub("datetime","date") => true
//     HasSub("datetime","time") => true
func HasAny(raw string, sub ...string) bool {
	for _, v := range sub {
		if strings.Contains(raw, v) {
			return true
		}
	}
	return false
}

// raw中是否包含sub...中的所有子串
//  如: HasSub("datetime","date","time") => true
//     HasSub("datetime","time") => false
func HasAll(raw string, sub ...string) bool {
	for _, v := range sub {
		if !strings.Contains(raw, v) {
			return false
		}
	}
	return true
}

// 判断tar是否与patterns中的任意规则所匹配
func MatchAny(tar string, patterns ...string) bool {
	if ContainString(tar, patterns...) {
		return true
	}

	for _, p := range patterns {
		switch {
		case strings.HasPrefix(p, "*"):
			p = fmt.Sprintf(".%s$", p) // *abc
		case strings.HasSuffix(p, "*"):
			p = fmt.Sprintf("^%s", strings.ReplaceAll(p, "*", ".*")) // abc*
		default:
			p = fmt.Sprintf("^%s$", strings.ReplaceAll(p, "*", ".*")) // ab*cd
		}
		// fmt.Printf("[%s]	", p)
		if b, err := regexp.MatchString(p, tar); err == nil && b {
			return true
		}
	}
	return false
}
