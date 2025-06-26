package timer

import "time"

func GetMinutes(timeStr1, timeStr2 string) float64 {
	// 定义时间格式
	layout := "2006-01-02 15:04:05"

	t1, _ := time.Parse(layout, timeStr1)
	t2, _ := time.Parse(layout, timeStr2)

	// 计算分钟差
	minutes := t2.Sub(t1).Minutes()

	return minutes

}
