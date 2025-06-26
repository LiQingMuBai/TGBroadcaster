package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestGetMinutes(t *testing.T) {
	// 解析时间字符串
	timeStr1 := "2023-05-15 10:30:10"
	timeStr2 := "2023-05-15 11:45:50"

	minutes := GetMinutes(timeStr1, timeStr2)

	fmt.Printf("时间相差 %.0f 分钟\n", minutes)

	// 当前时间
	now := time.Now()

	// 另一个时间（例如1小时30分钟前）
	pastTime := now.Add(-90 * time.Minute)

	// 计算分钟差
	minutes2 := now.Sub(pastTime).Minutes()

	fmt.Printf("时间相差 %.0f 分钟\n", minutes2)

}
