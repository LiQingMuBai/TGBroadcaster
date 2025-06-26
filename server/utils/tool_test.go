package utils

import (
	"fmt"
	"testing"
)

func TestSubtractDaysFromUnix(t *testing.T) {
	timestamp := int64(1718323200) // 假设这是某个时间戳
	result := SubtractDaysFromUnix(timestamp, 7)
	fmt.Printf("结果: %v (Unix: %d)\n", result, result.Unix())
}

func TestGetTimeDaysAgo(t *testing.T) {
	GetTimeDaysAgo(7)
}
