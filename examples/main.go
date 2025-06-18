package main

import (
	"fmt"

	"github.com/arixbit/zhkit"
)

func main() {
	fmt.Println("=== ZhKit Go 中文处理工具包示例 ===")
	fmt.Println()

	// 创建完整数据实例（推荐方式）
	chinese := zhkit.NewChineseWithFullData()
	fmt.Println("✓ 完整数据实例创建成功")
	fmt.Println()

	// 1. 汉字转拼音示例
	fmt.Println("=== 1. 汉字转拼音 ===")

	// 基本拼音转换
	result, _ := chinese.ToPinyin("中华人民共和国", zhkit.ModePinyin, " ", false)
	fmt.Printf("ToPinyin(中华人民共和国) -> %v\n", result.Pinyin)

	// 首字母转换
	result, _ = chinese.ToPinyin("中华人民共和国", zhkit.ModePinyinFirst, " ", false)
	fmt.Printf("ToPinyinFirst(中华人民共和国) -> %v\n", result.PinyinFirst)

	// 测试更多汉字
	result, _ = chinese.ToPinyin("北京大学", zhkit.ModePinyin, " ", false)
	fmt.Printf("ToPinyin(北京大学) -> %v\n", result.Pinyin)

	result, _ = chinese.ToPinyin("测试加载数据", zhkit.ModePinyin, " ", false)
	fmt.Printf("ToPinyin(测试加载数据) -> %v\n", result.Pinyin)
	fmt.Println()

	// 2. 拼音分词示例
	fmt.Println("=== 2. 拼音分词 ===")
	pinyinSplit, _ := chinese.SplitPinyin("beijing")
	fmt.Printf("SplitPinyin(beijing) -> %v\n", pinyinSplit)

	pinyinSplit, _ = chinese.SplitPinyin("qinghuadaxue")
	fmt.Printf("SplitPinyin(qinghuadaxue) -> %v\n", pinyinSplit)
	fmt.Println()

	// 3. 简繁转换示例
	fmt.Println("=== 3. 简繁转换 ===")
	simplified, _ := chinese.ToSimplified("繁體中文")
	fmt.Printf("ToSimplified(繁體中文) -> %v\n", simplified)

	traditional, _ := chinese.ToTraditional("简体中文")
	fmt.Printf("ToTraditional(简体中文) -> %v\n", traditional)

	// 测试更多简繁转换
	simplified, _ = chinese.ToTraditional("恭喜发财")
	fmt.Printf("ToTraditional(恭喜发财) -> %v\n", simplified)

	simplified, _ = chinese.ToTraditional("红包拿来")
	fmt.Printf("ToTraditional(红包拿来) -> %v\n", simplified)

	traditional, _ = chinese.ToTraditional("八方來財")
	fmt.Printf("ToSimplified(八方来财) -> %v\n", traditional)
	fmt.Println()

	// 4. 数字转换示例
	fmt.Println("=== 4. 数字转换 ===")
	chineseNum, _ := chinese.ToChineseNumber(2025, nil)
	fmt.Printf("ToChineseNumber(2025) -> %s\n", chineseNum)

	chineseNum, _ = chinese.ToChineseNumber(9876.543201, nil)
	fmt.Printf("ToChineseNumber(9876.543201) -> %s\n", chineseNum)

	chineseNum, _ = chinese.ToChineseNumber(-100.345, nil)
	fmt.Printf("ToChineseNumber(-100.345) -> %s\n", chineseNum)
	fmt.Println()

	// 5. 金额转换示例
	fmt.Println("=== 5. 金额转换 ===")
	currency, _ := chinese.ToCurrencyNumber(2024.50, "元")
	fmt.Printf("ToCurrencyNumber(2024.50) -> %s\n", currency)

	currency, _ = chinese.ToCurrencyNumber(1000, "元")
	fmt.Printf("ToCurrencyNumber(1000) -> %s\n", currency)

	// 使用默认单位
	currency, _ = chinese.ToCurrencyNumber(888.88, "")
	fmt.Printf("ToCurrencyNumber(888.88) -> %s\n", currency)
	fmt.Println()

	// 6. 中文数字转阿拉伯数字示例
	fmt.Println("=== 6. 中文数字转换 ===")
	number, _ := chinese.ChineseToNumber("二零二五")
	fmt.Printf("ChineseToNumber(二零二五) -> %.0f\n", number)

	number, _ = chinese.ChineseToNumber("一千二百三十四")
	fmt.Printf("ChineseToNumber(一千二百三十四) -> %.0f\n", number)

	number, _ = chinese.ChineseToNumber("十二点五")
	fmt.Printf("ChineseToNumber(十二点五) -> %.1f\n", number)

	number, _ = chinese.ChineseToNumber("负一百")
	fmt.Printf("ChineseToNumber(负一百) -> %.0f\n", number)
	fmt.Println()

	// 7. 全局函数示例
	fmt.Println("=== 7. 全局函数示例 ===")
	globalResult, _ := zhkit.ToPinyin("全局函数测试", zhkit.ModePinyin, " ", false)
	fmt.Printf("zhkit.ToPinyin(全局函数测试) -> %v\n", globalResult.Pinyin)

	globalSimplified, _ := zhkit.ToSimplified("繁體字")
	fmt.Printf("zhkit.ToSimplified(繁體字) -> %v\n", globalSimplified)

	globalChinese, _ := zhkit.ToChineseNumber(999, nil)
	fmt.Printf("zhkit.ToChineseNumber(999) -> %s\n", globalChinese)
	fmt.Println()

	fmt.Println("=== 示例结束 ===")
}
