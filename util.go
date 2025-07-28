// convertToBytes 将给定的值和单位转换为字节数
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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

// getFileInfo 获取文件信息
func getFileInfo(filename string) (fileLine int, fileSize string, err error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return
	}
	fileSize = bytesToString(fileInfo.Size())

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 4096) // 设置缓冲区大小
	for {
		_, err := reader.ReadBytes('\n') // 读取到下一个换行符
		if err != nil {
			if err == io.EOF {
				break // 文件结束
			}
			return fileLine, fileSize, err
		}
		fileLine++
	}
	return
}
