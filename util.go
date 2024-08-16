// convertToBytes 将给定的值和单位转换为字节数
package main

import "fmt"

//
func convertToBytes(value int, unit string) int {
	switch unit {
	case UNIT_B:
		return value
	case UNIT_KB:
		return value * 1024
	case UNIT_MB:
		return value * 1024 * 1024
	case UNIT_GB:
		return value * 1024 * 1024 * 1024
	default:
		return -1 // 错误的单位
	}
}

// bytesToString 将字节数转换为字符串
func bytesToString(value int64) string {
	if value < 1024 {
		return fmt.Sprintf("%d %s", value, UNIT_B)
	} else if value < 1024*1024 {
		return fmt.Sprintf("%d %s", value/1024, UNIT_KB)
	} else if value < 1024*1024*1024 {
		return fmt.Sprintf("%d %s", value/1024/1024, UNIT_MB)
	} else {
		return fmt.Sprintf("%d %s", value/1024/1024/1024, UNIT_GB)
	}
}
