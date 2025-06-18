package zhkit

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 数字转换相关常量
var (
	// 中文数字字符
	chineseNumbers = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	// 中文单位
	chineseUnits = []string{"", "十", "百", "千"}
	// 中文大单位
	chineseBigUnits = []string{"", "万", "亿", "兆"}
	
	// 金额专用大写数字
	currencyNumbers = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	// 金额专用单位
	currencyUnits = []string{"", "拾", "佰", "仟"}
	// 金额专用大单位
	currencyBigUnits = []string{"", "万", "亿", "兆"}
	// 金额小数单位
	currencyDecimalUnits = []string{"角", "分"}
)

// ToChineseNumber 阿拉伯数字转中文数字
// number: 要转换的数字（支持整数和小数）
// options: 转换选项
func (c *Chinese) ToChineseNumber(number interface{}, options *NumberOptions) (string, error) {
	if options == nil {
		options = &NumberOptions{}
	}
	
	var numStr string
	switch v := number.(type) {
	case int:
		numStr = strconv.Itoa(v)
	case int64:
		numStr = strconv.FormatInt(v, 10)
	case float64:
		numStr = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		numStr = v
	default:
		return "", errors.New("不支持的数字类型")
	}
	
	// 处理负数
	isNegative := false
	if strings.HasPrefix(numStr, "-") {
		isNegative = true
		numStr = numStr[1:]
	}
	
	// 分离整数和小数部分
	parts := strings.Split(numStr, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}
	
	// 转换整数部分
	integerChinese, err := c.convertIntegerToChinese(integerPart, options)
	if err != nil {
		return "", err
	}
	
	// 转换小数部分
	decimalChinese := ""
	if decimalPart != "" {
		decimalChinese = c.convertDecimalToChinese(decimalPart)
	}
	
	// 组合结果
	result := integerChinese
	if decimalChinese != "" {
		result += "点" + decimalChinese
	}
	
	if isNegative {
		result = "负" + result
	}
	
	return result, nil
}

// ToCurrencyNumber 数字转金额大写
// amount: 金额数字
// unit: 货币单位（如"元"、"圆"等）
func (c *Chinese) ToCurrencyNumber(amount interface{}, unit string) (string, error) {
	if unit == "" {
		unit = "元"
	}
	
	var amountStr string
	switch v := amount.(type) {
	case int:
		amountStr = strconv.Itoa(v)
	case int64:
		amountStr = strconv.FormatInt(v, 10)
	case float64:
		// 保留两位小数
		amountStr = fmt.Sprintf("%.2f", v)
	case string:
		amountStr = v
	default:
		return "", errors.New("不支持的金额类型")
	}
	
	// 处理负数
	isNegative := false
	if strings.HasPrefix(amountStr, "-") {
		isNegative = true
		amountStr = amountStr[1:]
	}
	
	// 分离整数和小数部分
	parts := strings.Split(amountStr, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
		// 确保小数部分为两位
		if len(decimalPart) == 1 {
			decimalPart += "0"
		} else if len(decimalPart) > 2 {
			decimalPart = decimalPart[:2]
		}
	} else {
		decimalPart = "00"
	}
	
	// 转换整数部分
	integerChinese, err := c.convertIntegerToCurrency(integerPart)
	if err != nil {
		return "", err
	}
	
	// 转换小数部分
	decimalChinese := c.convertDecimalToCurrency(decimalPart)
	
	// 组合结果
	result := ""
	if isNegative {
		result = "负"
	}
	
	if integerChinese == "" || integerChinese == "零" {
		if decimalChinese == "" {
			result += "零" + unit + "整"
		} else {
			result += decimalChinese
		}
	} else {
		result += integerChinese + unit
		if decimalChinese == "" {
			result += "整"
		} else {
			result += decimalChinese
		}
	}
	
	return result, nil
}

// convertIntegerToChinese 转换整数部分为中文
func (c *Chinese) convertIntegerToChinese(integerStr string, options *NumberOptions) (string, error) {
	if integerStr == "" || integerStr == "0" {
		return "零", nil
	}
	
	// 移除前导零
	integerStr = strings.TrimLeft(integerStr, "0")
	if integerStr == "" {
		return "零", nil
	}
	
	length := len(integerStr)
	if length > 16 {
		return "", errors.New("数字过大，超出处理范围")
	}
	
	result := ""
	zeroFlag := false
	
	// 按4位分组处理
	for i := 0; i < length; i++ {
		digit := int(integerStr[i] - '0')
		pos := length - i - 1
		bigUnitPos := pos / 4
		unitPos := pos % 4
		
		if digit == 0 {
			// 处理零
			if unitPos == 0 && bigUnitPos > 0 {
				// 万、亿等位置的零需要特殊处理
				if !strings.HasSuffix(result, chineseBigUnits[bigUnitPos]) {
					result += chineseBigUnits[bigUnitPos]
				}
			}
			zeroFlag = true
		} else {
			// 处理非零数字
			if zeroFlag && result != "" {
				result += "零"
			}
			zeroFlag = false
			
			result += chineseNumbers[digit]
			
			// 添加单位
			if unitPos > 0 {
				result += chineseUnits[unitPos]
			}
			
			// 添加大单位
			if unitPos == 0 && bigUnitPos > 0 {
				result += chineseBigUnits[bigUnitPos]
			}
		}
	}
	
	// 处理"一十"的特殊情况
	if options.TenMin && strings.HasPrefix(result, "一十") {
		result = "十" + result[2:]
	}
	
	return result, nil
}

// convertDecimalToChinese 转换小数部分为中文
func (c *Chinese) convertDecimalToChinese(decimalStr string) string {
	if decimalStr == "" {
		return ""
	}
	
	result := ""
	for _, char := range decimalStr {
		digit := int(char - '0')
		result += chineseNumbers[digit]
	}
	
	return result
}

// convertIntegerToCurrency 转换整数部分为金额大写
func (c *Chinese) convertIntegerToCurrency(integerStr string) (string, error) {
	if integerStr == "" || integerStr == "0" {
		return "零", nil
	}
	
	// 移除前导零
	integerStr = strings.TrimLeft(integerStr, "0")
	if integerStr == "" {
		return "零", nil
	}
	
	length := len(integerStr)
	if length > 16 {
		return "", errors.New("金额过大，超出处理范围")
	}
	
	result := ""
	zeroFlag := false
	
	// 按4位分组处理
	for i := 0; i < length; i++ {
		digit := int(integerStr[i] - '0')
		pos := length - i - 1
		bigUnitPos := pos / 4
		unitPos := pos % 4
		
		if digit == 0 {
			// 处理零
			if unitPos == 0 && bigUnitPos > 0 {
				// 万、亿等位置的零需要特殊处理
				if !strings.HasSuffix(result, currencyBigUnits[bigUnitPos]) {
					result += currencyBigUnits[bigUnitPos]
				}
			}
			zeroFlag = true
		} else {
			// 处理非零数字
			if zeroFlag && result != "" {
				result += "零"
			}
			zeroFlag = false
			
			result += currencyNumbers[digit]
			
			// 添加单位
			if unitPos > 0 {
				result += currencyUnits[unitPos]
			}
			
			// 添加大单位
			if unitPos == 0 && bigUnitPos > 0 {
				result += currencyBigUnits[bigUnitPos]
			}
		}
	}
	
	return result, nil
}

// convertDecimalToCurrency 转换小数部分为金额大写
func (c *Chinese) convertDecimalToCurrency(decimalStr string) string {
	if len(decimalStr) != 2 {
		return ""
	}
	
	jiao := int(decimalStr[0] - '0')
	fen := int(decimalStr[1] - '0')
	
	result := ""
	
	if jiao > 0 {
		result += currencyNumbers[jiao] + currencyDecimalUnits[0]
	}
	
	if fen > 0 {
		if jiao == 0 {
			result += "零"
		}
		result += currencyNumbers[fen] + currencyDecimalUnits[1]
	}
	
	return result
}

// ChineseToNumber 中文数字转阿拉伯数字
func (c *Chinese) ChineseToNumber(chineseNum string) (float64, error) {
	if chineseNum == "" {
		return 0, errors.New("输入为空")
	}
	
	// 处理负数
	isNegative := false
	if strings.HasPrefix(chineseNum, "负") {
		isNegative = true
		runes := []rune(chineseNum)
		chineseNum = string(runes[1:])
	}
	
	// 分离整数和小数部分
	parts := strings.Split(chineseNum, "点")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = parts[1]
	}
	
	// 转换整数部分
	integerValue, err := c.convertChineseToInteger(integerPart)
	if err != nil {
		return 0, err
	}
	
	// 转换小数部分
	decimalValue := 0.0
	if decimalPart != "" {
		decimalValue, err = c.convertChineseToDecimal(decimalPart)
		if err != nil {
			return 0, err
		}
	}
	
	result := float64(integerValue) + decimalValue
	if isNegative {
		result = -result
	}
	
	return result, nil
}

// convertChineseToInteger 转换中文整数为数字
func (c *Chinese) convertChineseToInteger(chineseNum string) (int64, error) {
	if chineseNum == "零" {
		return 0, nil
	}
	
	// 数字映射表
	numberMap := map[string]int{
		"零": 0, "一": 1, "二": 2, "三": 3, "四": 4, "五": 5,
		"六": 6, "七": 7, "八": 8, "九": 9, "十": 10,
		"壹": 1, "贰": 2, "叁": 3, "肆": 4, "伍": 5,
		"陆": 6, "柒": 7, "捌": 8, "玖": 9, "拾": 10,
	}
	
	// 单个字符直接映射
	if val, exists := numberMap[chineseNum]; exists {
		return int64(val), nil
	}
	
	// 处理复合数字（如"十二"）
	if strings.Contains(chineseNum, "十") {
		return c.parseChineseTen(chineseNum)
	}
	
	// 处理连续数字（如"二零二四"）
	runes := []rune(chineseNum)
	result := int64(0)
	for _, r := range runes {
		charStr := string(r)
		if val, exists := numberMap[charStr]; exists {
			result = result*10 + int64(val)
		} else {
			return 0, errors.New("无法解析的中文数字: " + chineseNum)
		}
	}
	
	return result, nil
}

// parseChineseTen 解析包含"十"的中文数字
func (c *Chinese) parseChineseTen(chineseNum string) (int64, error) {
	// 简化实现
	if chineseNum == "十" {
		return 10, nil
	}
	if strings.HasPrefix(chineseNum, "十") {
		// 十一、十二等
		if len(chineseNum) > 1 {
			lastChar := string([]rune(chineseNum)[1])
			numberMap := map[string]int{
				"一": 1, "二": 2, "三": 3, "四": 4, "五": 5,
				"六": 6, "七": 7, "八": 8, "九": 9,
			}
			if val, exists := numberMap[lastChar]; exists {
				return int64(10 + val), nil
			}
		}
	}
	
	return 0, errors.New("无法解析的十位数字: " + chineseNum)
}

// convertChineseToDecimal 转换中文小数为数字
func (c *Chinese) convertChineseToDecimal(chineseDecimal string) (float64, error) {
	numberMap := map[string]int{
		"零": 0, "一": 1, "二": 2, "三": 3, "四": 4, "五": 5,
		"六": 6, "七": 7, "八": 8, "九": 9,
	}
	
	result := 0.0
	divisor := 10.0
	
	for _, char := range chineseDecimal {
		charStr := string(char)
		if val, exists := numberMap[charStr]; exists {
			result += float64(val) / divisor
			divisor *= 10
		} else {
			return 0, errors.New("无法解析的中文小数: " + chineseDecimal)
		}
	}
	
	return result, nil
}

// 全局函数

// ToChineseNumber 全局函数：阿拉伯数字转中文数字
func ToChineseNumber(number interface{}, options *NumberOptions) (string, error) {
	return defaultChinese.ToChineseNumber(number, options)
}

// ToCurrencyNumber 全局函数：数字转金额大写
func ToCurrencyNumber(amount interface{}, unit string) (string, error) {
	return defaultChinese.ToCurrencyNumber(amount, unit)
}

// ChineseToNumber 全局函数：中文数字转阿拉伯数字
func ChineseToNumber(chineseNum string) (float64, error) {
	return defaultChinese.ChineseToNumber(chineseNum)
}