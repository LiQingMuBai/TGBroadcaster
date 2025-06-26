package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestHoursSince(t *testing.T) {
	// 示例：假设给定的时间戳是 1718640000000（2024-06-17 16:00:00 UTC）
	unixMillis := int64(1718640000000)
	hours := HoursSince(unixMillis)
	fmt.Printf("相差 %.2f 小时\n", hours)

}

func TestDaysSince(t *testing.T) {
	// 示例：2024-06-10 00:00:00 UTC (1717977600000)
	unixMillis := int64(1717977600000)
	days := DaysSince(unixMillis)
	fmt.Printf("相差 %.2f 天\n", days) // 如：8.54 天
}
func TestFloatEqual(t *testing.T) {
	x := 0.1 + 0.2
	y := 0.3

	fmt.Println("直接比较 (x == y):", x == y)               // false（因为浮点误差）
	fmt.Println("安全比较 (FloatEqual):", FloatEqual(x, y)) // true
}
func TestCompareBalances(t *testing.T) {

	flag := CompareBalances(19.22, 120)

	log.Println(flag)

}
func TestSuoJinSuanFa(t *testing.T) {

	fmt.Println(SuoJinSuanFa("2071.82")) // 输出: 0.05001
	fmt.Println(SuoJinSuanFa("5000"))    // 输出: 0.05001
	fmt.Println(SuoJinSuanFa("80"))      // 输出: 0.08
	fmt.Println(SuoJinSuanFa("1000"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("1010"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("1002"))    // 输出: 0.01001
	fmt.Println(SuoJinSuanFa("8"))       // 输出: 0.01001

}
func TestSuoJinFa2(t *testing.T) {

	fmt.Println(SuoJinSuanFa2("19785000")) // 输出: 1.1978501
	fmt.Println(SuoJinSuanFa2("5000"))     // 输出: 0.501
	fmt.Println(SuoJinSuanFa2("80"))       // 输出: 0.80
	fmt.Println(SuoJinSuanFa2("69863"))    // 输出: 0.69863

}
func TestSuoJinSuanFaReverse(t *testing.T) {
	testCases := []string{
		"0.05001",  // 5000
		"0.08",     // 80
		"0.069863", // 69863
		"0.01001",  // 1000
		"0.123",    // 123（无补1情况）
	}

	for _, s := range testCases {
		original, err := SuoJinSuanFaReverse(s)
		if err != nil {
			fmt.Printf("错误：%v\n", err)
			continue
		}
		fmt.Printf("格式化：%-10s => 原始值：%d\n", s, original)
	}
}
func TestSuoJinSuanFa2Reverse(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{"19785000", "1.9785001"},
		{"5000", "0.5001"},
		{"80", "0.8"},
		{"69863", "0.69863"},
		{"100", "1.11"},
		{"1234500", "1.23451"},
		{"200", "0.21"},
	}

	fmt.Println("=== 格式化测试 ===")
	for _, tc := range testCases {
		result := SuoJinSuanFa2(tc.input)
		fmt.Printf("%s → %s (期望: %s)\n", tc.input, result, tc.output)
		//if result != tc.output {
		//	fmt.Println("  错误！不匹配")
		//}
	}

	fmt.Println("\n=== 反向解析测试 ===")
	for _, tc := range testCases {
		reversed, err := SuoJinSuanFa2Reverse(tc.output)
		if err != nil {
			fmt.Printf("解析 %s 错误: %v\n", tc.output, err)
			continue
		}
		fmt.Printf("%s → %s (原始: %s)\n", tc.output, reversed, tc.input)
		//if reversed != tc.input {
		//	fmt.Println("  错误！不匹配")
		//}
	}

	// 额外测试一些边缘情况
	extraTests := []struct {
		formatted string
		expected  string
	}{
		{"1.2345671", "123456700"},
		{"0.1231", "12300"},
		{"1.11", "100"},
		{"0.21", "200"},
	}

	fmt.Println("\n=== 额外反向解析测试 ===")
	for _, et := range extraTests {
		reversed, err := SuoJinSuanFa2Reverse(et.formatted)
		if err != nil {
			fmt.Printf("解析 %s 错误: %v\n", et.formatted, err)
			continue
		}
		fmt.Printf("%s → %s (期望: %s)\n", et.formatted, reversed, et.expected)
		//if reversed != et.expected {
		//	fmt.Println("  错误！不匹配")
		//}
	}
}
