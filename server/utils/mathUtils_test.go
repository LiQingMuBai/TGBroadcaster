package utils

import (
	"fmt"
	"testing"
)

func TestMultiplyBy100ToInt64(t *testing.T) {
	// 示例用法
	floatStr := "16.88"
	result, err := MultiplyBy100ToInt64Decimal(floatStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", result) // 输出: 12345
}
