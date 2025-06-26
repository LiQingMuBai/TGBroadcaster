package utils

import (
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
	"time"
)

func AddDecimalStringsWithPrecision(a, b string, precision int32) (string, error) {
	num1, err := decimal.NewFromString(a)
	if err != nil {
		return "", err
	}

	num2, err := decimal.NewFromString(b)
	if err != nil {
		return "", err
	}

	result := num1.Add(num2)
	return result.StringFixed(precision), nil
}

// GenerateTronOrderID 生成波场订单号（年月日时分 + 波场地址后4位）
func GenerateOrderID(userId string) string {
	// 1. 校验波场地址格式
	userId = strings.TrimSpace(userId)
	//if len(tronAddress) != 34 || !strings.HasPrefix(tronAddress, "T") {
	//	return "", fmt.Errorf("无效的波场地址（必须34位且以T开头）")
	//}

	// 2. 获取当前时间的 "年月日时分"（格式：200601021504）
	timestamp := time.Now().Format("200601021504")

	// 4. 拼接时间 + 地址后4位
	orderID := timestamp + userId
	return orderID
}
func StringToUint32(input string) uint32 {
	// First convert to uint64 to check for overflow
	val, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(val)
}

func MultiplyBy100ToInt64Decimal(floatStr string) (int64, error) {
	d, err := decimal.NewFromString(floatStr)
	if err != nil {
		return 0, err
	}

	result := d.Mul(decimal.NewFromInt(100))
	return result.IntPart(), nil
}
