package web

import (
	"net/url"
	"strconv"
	"strings"
)

// TemplateFuncs 模板辅助函数
var TemplateFuncs = map[string]interface{}{
	"buildQuery": buildQuery,
	"add":        add,
	"sub":        sub,
	"seq":        seq,
	"contains":   contains,
	"join":       join,
	"deref":      deref,
	"gt":         gt,
	"lt":         lt,
	"eq":         eq,
}

// buildQuery 构建查询参数
func buildQuery(data interface{}, key string, value interface{}) string {
	// 这里简化处理，实际应该从当前请求中获取参数
	params := url.Values{}

	// 添加新参数
	params.Set(key, toString(value))

	return params.Encode()
}

// add 加法运算
func add(a, b int) int {
	return a + b
}

// sub 减法运算
func sub(a, b int) int {
	return a - b
}

// seq 生成序列
func seq(start, end int) []int {
	if start > end {
		return []int{}
	}
	result := make([]int, end-start+1)
	for i := start; i <= end; i++ {
		result[i-start] = i
	}
	return result
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// join 连接字符串
func join(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

// deref 解引用指针
func deref(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}

	// 这里简化处理，实际应该根据类型进行解引用
	// 对于 *float64 类型，我们直接返回指针的值
	if f, ok := ptr.(*float64); ok {
		if f == nil {
			return nil
		}
		return *f
	}

	return ptr
}

// gt 大于比较
func gt(a, b interface{}) bool {
	return compareNumbers(a, b) > 0
}

// lt 小于比较
func lt(a, b interface{}) bool {
	return compareNumbers(a, b) < 0
}

// eq 等于比较
func eq(a, b interface{}) bool {
	// 首先尝试字符串比较
	if aStr, ok := a.(string); ok {
		if bStr, ok := b.(string); ok {
			return aStr == bStr
		}
	}
	// 然后尝试数字比较
	return compareNumbers(a, b) == 0
}

// compareNumbers 比较数字
func compareNumbers(a, b interface{}) int {
	// 转换为float64进行比较
	var aFloat, bFloat float64

	switch v := a.(type) {
	case float64:
		aFloat = v
	case int:
		aFloat = float64(v)
	case int64:
		aFloat = float64(v)
	default:
		return 0
	}

	switch v := b.(type) {
	case float64:
		bFloat = v
	case int:
		bFloat = float64(v)
	case int64:
		bFloat = float64(v)
	default:
		return 0
	}

	if aFloat > bFloat {
		return 1
	} else if aFloat < bFloat {
		return -1
	}
	return 0
}

// toString 转换为字符串
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}
