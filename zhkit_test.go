package zhkit

import (
	"testing"
)

func TestToPinyin(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		text     string
		mode     ConvertMode
		expected bool // 是否期望成功
	}{
		{
			name:     "基本汉字转拼音",
			text:     "中国",
			mode:     ModePinyin,
			expected: true,
		},
		{
			name:     "空字符串",
			text:     "",
			mode:     ModePinyin,
			expected: true,
		},
		{
			name:     "混合中英文",
			text:     "中国China",
			mode:     ModePinyin,
			expected: true,
		},
		{
			name:     "多种模式组合",
			text:     "中国",
			mode:     ModePinyin | ModePinyinFirst,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ToPinyin(tt.text, tt.mode, " ", false)
			if tt.expected {
				if err != nil {
					t.Errorf("ToPinyin() error = %v, expected success", err)
					return
				}
				if result == nil {
					t.Errorf("ToPinyin() result is nil")
					return
				}
				t.Logf("ToPinyin(%s) = %+v", tt.text, result)
			} else {
				if err == nil {
					t.Errorf("ToPinyin() expected error, got success")
				}
			}
		})
	}
}

func TestSplitPinyin(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		pinyin   string
		expected bool
	}{
		{
			name:     "基本拼音分词",
			pinyin:   "beijing",
			expected: true,
		},
		{
			name:     "复杂拼音分词",
			pinyin:   "xianggang",
			expected: true,
		},
		{
			name:     "空字符串",
			pinyin:   "",
			expected: true,
		},
		{
			name:     "单个拼音",
			pinyin:   "zhong",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.SplitPinyin(tt.pinyin)
			if tt.expected {
				if err != nil {
					t.Errorf("SplitPinyin() error = %v, expected success", err)
					return
				}
				t.Logf("SplitPinyin(%s) = %v", tt.pinyin, result)
			} else {
				if err == nil {
					t.Errorf("SplitPinyin() expected error, got success")
				}
			}
		})
	}
}

func TestSplitPinyinArray(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		pinyin   string
		expected bool
	}{
		{
			name:     "基本拼音分词数组",
			pinyin:   "beijing",
			expected: true,
		},
		{
			name:     "复杂拼音分词数组",
			pinyin:   "xianggang",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.SplitPinyinArray(tt.pinyin)
			if tt.expected {
				if err != nil {
					t.Errorf("SplitPinyinArray() error = %v, expected success", err)
					return
				}
				t.Logf("SplitPinyinArray(%s) = %v", tt.pinyin, result)
			} else {
				if err == nil {
					t.Errorf("SplitPinyinArray() expected error, got success")
				}
			}
		})
	}
}

func TestToSimplified(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "繁体转简体",
			text:     "發財",
			expected: true,
		},
		{
			name:     "已是简体",
			text:     "发财",
			expected: true,
		},
		{
			name:     "空字符串",
			text:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ToSimplified(tt.text)
			if tt.expected {
				if err != nil {
					t.Errorf("ToSimplified() error = %v, expected success", err)
					return
				}
				t.Logf("ToSimplified(%s) = %v", tt.text, result)
			} else {
				if err == nil {
					t.Errorf("ToSimplified() expected error, got success")
				}
			}
		})
	}
}

func TestToTraditional(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "简体转繁体",
			text:     "发财",
			expected: true,
		},
		{
			name:     "已是繁体",
			text:     "發財",
			expected: true,
		},
		{
			name:     "空字符串",
			text:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ToTraditional(tt.text)
			if tt.expected {
				if err != nil {
					t.Errorf("ToTraditional() error = %v, expected success", err)
					return
				}
				t.Logf("ToTraditional(%s) = %v", tt.text, result)
			} else {
				if err == nil {
					t.Errorf("ToTraditional() expected error, got success")
				}
			}
		})
	}
}

func TestToChineseNumber(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		number   interface{}
		options  *NumberOptions
		expected bool
	}{
		{
			name:     "整数转换",
			number:   123,
			options:  nil,
			expected: true,
		},
		{
			name:     "小数转换",
			number:   123.45,
			options:  nil,
			expected: true,
		},
		{
			name:     "字符串数字",
			number:   "456",
			options:  nil,
			expected: true,
		},
		{
			name:     "负数",
			number:   -123,
			options:  nil,
			expected: true,
		},
		{
			name:     "十的特殊处理",
			number:   12,
			options:  &NumberOptions{TenMin: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ToChineseNumber(tt.number, tt.options)
			if tt.expected {
				if err != nil {
					t.Errorf("ToChineseNumber() error = %v, expected success", err)
					return
				}
				t.Logf("ToChineseNumber(%v) = %s", tt.number, result)
			} else {
				if err == nil {
					t.Errorf("ToChineseNumber() expected error, got success")
				}
			}
		})
	}
}

func TestToCurrencyNumber(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name     string
		amount   interface{}
		unit     string
		expected bool
	}{
		{
			name:     "整数金额",
			amount:   123,
			unit:     "元",
			expected: true,
		},
		{
			name:     "小数金额",
			amount:   123.45,
			unit:     "元",
			expected: true,
		},
		{
			name:     "字符串金额",
			amount:   "456.78",
			unit:     "元",
			expected: true,
		},
		{
			name:     "负金额",
			amount:   -123.45,
			unit:     "元",
			expected: true,
		},
		{
			name:     "零金额",
			amount:   0,
			unit:     "元",
			expected: true,
		},
		{
			name:     "默认单位",
			amount:   100,
			unit:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ToCurrencyNumber(tt.amount, tt.unit)
			if tt.expected {
				if err != nil {
					t.Errorf("ToCurrencyNumber() error = %v, expected success", err)
					return
				}
				t.Logf("ToCurrencyNumber(%v, %s) = %s", tt.amount, tt.unit, result)
			} else {
				if err == nil {
					t.Errorf("ToCurrencyNumber() expected error, got success")
				}
			}
		})
	}
}

func TestChineseToNumber(t *testing.T) {
	chinese := NewChinese()

	tests := []struct {
		name        string
		chineseNum  string
		expected    bool
		expectedVal float64
	}{
		{
			name:        "基本数字",
			chineseNum:  "一",
			expected:    true,
			expectedVal: 1,
		},
		{
			name:        "十位数字",
			chineseNum:  "十",
			expected:    true,
			expectedVal: 10,
		},
		{
			name:        "十几",
			chineseNum:  "十二",
			expected:    true,
			expectedVal: 12,
		},
		{
			name:        "零",
			chineseNum:  "零",
			expected:    true,
			expectedVal: 0,
		},
		{
			name:        "负数",
			chineseNum:  "负一",
			expected:    true,
			expectedVal: -1,
		},
		{
			name:        "小数",
			chineseNum:  "一点二三",
			expected:    true,
			expectedVal: 1.23,
		},
		{
			name:        "连续数字",
			chineseNum:  "二零二四",
			expected:    true,
			expectedVal: 2024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := chinese.ChineseToNumber(tt.chineseNum)
			if tt.expected {
				if err != nil {
					t.Errorf("ChineseToNumber() error = %v, expected success", err)
					return
				}
				if result != tt.expectedVal {
					t.Errorf("ChineseToNumber() = %v, expected %v", result, tt.expectedVal)
				}
				t.Logf("ChineseToNumber(%s) = %v", tt.chineseNum, result)
			} else {
				if err == nil {
					t.Errorf("ChineseToNumber() expected error, got success")
				}
			}
		})
	}
}

// 测试全局函数
func TestGlobalFunctions(t *testing.T) {
	// 测试全局拼音转换函数
	result, err := ToPinyin("中国", ModePinyin, " ", false)
	if err != nil {
		t.Errorf("Global ToPinyin() error = %v", err)
	} else {
		t.Logf("Global ToPinyin(中国) = %+v", result)
	}

	// 测试全局拼音分词函数
	splitResult, err := SplitPinyin("beijing")
	if err != nil {
		t.Errorf("Global SplitPinyin() error = %v", err)
	} else {
		t.Logf("Global SplitPinyin(beijing) = %v", splitResult)
	}

	// 测试全局简繁转换函数
	simplified, err := ToSimplified("發財")
	if err != nil {
		t.Errorf("Global ToSimplified() error = %v", err)
	} else {
		t.Logf("Global ToSimplified(發財) = %v", simplified)
	}

	traditional, err := ToTraditional("发财")
	if err != nil {
		t.Errorf("Global ToTraditional() error = %v", err)
	} else {
		t.Logf("Global ToTraditional(发财) = %v", traditional)
	}

	// 测试全局数字转换函数
	chineseNum, err := ToChineseNumber(123, nil)
	if err != nil {
		t.Errorf("Global ToChineseNumber() error = %v", err)
	} else {
		t.Logf("Global ToChineseNumber(123) = %s", chineseNum)
	}

	currency, err := ToCurrencyNumber(123.45, "元")
	if err != nil {
		t.Errorf("Global ToCurrencyNumber() error = %v", err)
	} else {
		t.Logf("Global ToCurrencyNumber(123.45) = %s", currency)
	}

	number, err := ChineseToNumber("一二三")
	if err != nil {
		t.Logf("Global ChineseToNumber(一二三) expected error = %v", err)
	} else {
		t.Logf("Global ChineseToNumber(一二三) = %v", number)
	}
}

// 基准测试
func BenchmarkToPinyin(b *testing.B) {
	chinese := NewChinese()
	text := "中华人民共和国"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chinese.ToPinyin(text, ModePinyin, " ", false)
	}
}

func BenchmarkToChineseNumber(b *testing.B) {
	chinese := NewChinese()
	number := 123456789

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chinese.ToChineseNumber(number, nil)
	}
}

func BenchmarkToCurrencyNumber(b *testing.B) {
	chinese := NewChinese()
	amount := 123456.78

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = chinese.ToCurrencyNumber(amount, "元")
	}
}
