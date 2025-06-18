# ZhKit - Golang 中文工具包

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/arixbit/zhkit)](https://goreportcard.com/report/github.com/arixbit/zhkit)

ZhKit 是一个功能强大的 Golang 中文处理工具包，提供汉字转拼音、拼音分词、简繁互转、数字转换、金额转换等功能。本项目参考了 PHP 项目 [Yurunsoft/ChineseUtil](https://github.com/Yurunsoft/ChineseUtil) 的功能设计，使用 Go 语言重新实现，充分利用 Go 语言的特性，提供高性能的中文处理能力。

## 使用建议

### 最佳实践

1. **一般用途**：直接使用全局函数（如 `zhkit.ToPinyin()`），自动包含完整数据
2. **需要实例**：使用 `zhkit.NewChineseWithFullData()` 创建完整功能实例

### 数据说明

- **内嵌数据**：项目内置了完整的汉字拼音和简繁转换数据，无需额外文件
- **开箱即用**：无需手动加载数据文件，直接调用函数即可使用

## 功能特性

- ✅ **汉字转拼音**: 支持多种拼音格式（全拼、首字母、带声调等）
- ✅ **拼音分词**: 将连续的拼音字符串分割成独立的拼音
- ✅ **简繁互转**: 简体中文与繁体中文相互转换
- ✅ **数字转换**: 阿拉伯数字转中文数字，支持小数和负数
- ✅ **金额转换**: 数字转金额大写，支持多种货币单位
- ✅ **中文数字转换**: 中文数字转阿拉伯数字
- ✅ **高性能**: 基于内存的快速查找，无需数据库依赖
- ✅ **易于使用**: 提供简洁的 API 和全局函数

## 安装

```bash
go get github.com/arixbit/zhkit
```

## 快速开始

### 推荐使用方式（全局函数）

```go
package main

import (
    "fmt"
    "github.com/arixbit/zhkit"
)

func main() {
    // 直接使用全局函数（自动包含完整数据）
    result, _ := zhkit.ToPinyin("中华人民共和国", zhkit.ModePinyin, " ", false)
    fmt.Printf("拼音: %v\n", result.Pinyin)
    
    // 简繁转换
    simplified, _ := zhkit.ToSimplified("繁體中文")
    fmt.Printf("简体: %v\n", simplified)
    
    // 数字转换
    chineseNum, _ := zhkit.ToChineseNumber(2024, nil)
    fmt.Printf("中文数字: %s\n", chineseNum)
    
    // 拼音分词
    pinyinSplit, _ := zhkit.SplitPinyin("zhongguoren")
    fmt.Printf("拼音分词: %v\n", pinyinSplit)
}
```

### 实例使用方式

```go
package main

import (
    "fmt"
    "github.com/arixbit/zhkit"
)

func main() {
    // 创建完整数据实例
    chinese := zhkit.NewChineseWithFullData()
    
    result, _ := chinese.ToPinyin("测试文本", zhkit.ModePinyin, " ", false)
    fmt.Printf("拼音: %v\n", result.Pinyin)
}
```

### 使用全局函数

```go
package main

import (
    "fmt"
    "github.com/arixbit/zhkit"
)

func main() {
    // 直接使用全局函数
    result, _ := zhkit.ToPinyin("北京", zhkit.ModePinyin, " ", false)
    fmt.Println(result.Pinyin) // [[bei] [jing]]
    
    simplified, _ := zhkit.ToSimplified("繁體中文")
    fmt.Println(simplified) // [繁体中文]
    
    traditional, _ := zhkit.ToTraditional("简体中文")
    fmt.Println(traditional) // [簡體中文]
}
```

## 详细功能说明

### 1. 汉字转拼音

支持多种拼音转换模式：

```go
chinese := zhkit.NewChineseWithFullData()

// 基本拼音
result, _ := chinese.ToPinyin("中国人", zhkit.ModePinyin, " ", false)
fmt.Println(result.Pinyin) // [[zhong] [guo] [ren]]

// 首字母
result, _ = chinese.ToPinyin("中国人", zhkit.ModePinyinFirst, " ", false)
fmt.Println(result.PinyinFirst) // [[z] [g] [r]]

// 组合模式
result, _ = chinese.ToPinyin("中国人", zhkit.ModePinyin|zhkit.ModePinyinFirst, " ", false)
fmt.Println("拼音:", result.Pinyin)
fmt.Println("首字母:", result.PinyinFirst)
```

**转换模式说明：**
- `ModePinyin`: 全拼模式
- `ModePinyinFirst`: 首字母模式
- `ModePinyinSound`: 读音模式（带声调符号）
- `ModePinyinSoundNumber`: 读音数字模式（数字声调）

### 2. 拼音分词

将连续的拼音字符串分割成独立的拼音：

```go
chinese := zhkit.NewChineseWithFullData()

// 返回字符串格式
result, _ := chinese.SplitPinyin("beijing")
fmt.Println(result) // ["bei jing"]

// 返回数组格式
resultArray, _ := chinese.SplitPinyinArray("xianggang")
fmt.Println(resultArray) // [["xiang" "gang"]]
```

### 3. 简繁互转

```go
chinese := zhkit.NewChineseWithFullData()

// 繁体转简体
simplified, _ := chinese.ToSimplified("發財")
fmt.Println(simplified) // ["发财"]

// 简体转繁体
traditional, _ := chinese.ToTraditional("发财")
fmt.Println(traditional) // ["發財"]
```

### 4. 数字转换

```go
chinese := zhkit.NewChineseWithFullData()

// 基本数字转换
result, _ := chinese.ToChineseNumber(123, nil)
fmt.Println(result) // "一百二十三"

// 小数转换
result, _ = chinese.ToChineseNumber(123.45, nil)
fmt.Println(result) // "一百二十三点四五"

// 使用选项（十的特殊处理）
options := &zhkit.NumberOptions{TenMin: true}
result, _ = chinese.ToChineseNumber(12, options)
fmt.Println(result) // "十二" 而不是 "一十二"

// 负数转换
result, _ = chinese.ToChineseNumber(-123, nil)
fmt.Println(result) // "负一百二十三"
```

### 5. 金额转换

```go
chinese := zhkit.NewChineseWithFullData()

// 基本金额转换
result, _ := chinese.ToCurrencyNumber(123.45, "元")
fmt.Println(result) // "壹佰贰拾叁元肆角伍分"

// 整数金额
result, _ = chinese.ToCurrencyNumber(1000, "元")
fmt.Println(result) // "壹仟元整"

// 零金额
result, _ = chinese.ToCurrencyNumber(0, "元")
fmt.Println(result) // "零元整"

// 自定义货币单位
result, _ = chinese.ToCurrencyNumber(100, "圆")
fmt.Println(result) // "壹佰圆整"
```

### 6. 中文数字转阿拉伯数字

```go
chinese := zhkit.NewChineseWithFullData()

// 基本转换
result, _ := chinese.ChineseToNumber("一二三")
fmt.Println(result) // 123.0

// 十位数字
result, _ = chinese.ChineseToNumber("十二")
fmt.Println(result) // 12.0

// 小数转换
result, _ = chinese.ChineseToNumber("一点二三")
fmt.Println(result) // 1.23

// 负数转换
result, _ = chinese.ChineseToNumber("负一")
fmt.Println(result) // -1.0
```



## API 参考

### 类型定义

```go
// 转换模式
type ConvertMode int
const (
    ModePinyin ConvertMode = 1 << iota
    ModePinyinFirst
    ModePinyinSound
    ModePinyinSoundNumber
)

// 拼音转换结果
type PinyinResult struct {
    Pinyin           [][]string `json:"pinyin,omitempty"`
    PinyinFirst      [][]string `json:"pinyinFirst,omitempty"`
    PinyinSound      [][]string `json:"pinyinSound,omitempty"`
    PinyinSoundNumber [][]string `json:"pinyinSoundNumber,omitempty"`
}

// 数字转换选项
type NumberOptions struct {
    TenMin bool // "一十二" => "十二"
}
```

### 主要方法

```go
// 创建完整数据实例
func NewChineseWithFullData() *Chinese

// 汉字转拼音
func (c *Chinese) ToPinyin(text string, mode ConvertMode, separator string, splitNonChinese bool) (*PinyinResult, error)

// 拼音分词
func (c *Chinese) SplitPinyin(pinyin string) ([]string, error)
func (c *Chinese) SplitPinyinArray(pinyin string) ([][]string, error)

// 简繁转换
func (c *Chinese) ToSimplified(text string) ([]string, error)
func (c *Chinese) ToTraditional(text string) ([]string, error)

// 数字转换
func (c *Chinese) ToChineseNumber(number interface{}, options *NumberOptions) (string, error)
func (c *Chinese) ToCurrencyNumber(amount interface{}, unit string) (string, error)
func (c *Chinese) ChineseToNumber(chineseNum string) (float64, error)

```

### 全局函数

```go
// 全局拼音转换
func ToPinyin(text string, mode ConvertMode, separator string, splitNonChinese bool) (*PinyinResult, error)

// 全局拼音分词
func SplitPinyin(pinyin string) ([]string, error)
func SplitPinyinArray(pinyin string) ([][]string, error)

// 全局简繁转换
func ToSimplified(text string) ([]string, error)
func ToTraditional(text string) ([]string, error)

// 全局数字转换
func ToChineseNumber(number interface{}, options *NumberOptions) (string, error)
func ToCurrencyNumber(amount interface{}, unit string) (string, error)
func ChineseToNumber(chineseNum string) (float64, error)
```

## 性能特点

- **内存优化**: 使用高效的数据结构，减少内存占用
- **并发安全**: 支持多协程并发访问
- **快速查找**: 基于哈希表的 O(1) 查找性能
- **无外部依赖**: 不依赖数据库，纯内存操作

## 测试

运行测试：

```bash
# 运行所有测试
go test

# 运行测试并显示详细输出
go test -v

# 运行基准测试
go test -bench=.

# 运行特定测试
go test -run TestToPinyin
```

## 示例程序

查看 `examples/main.go` 文件获取完整的使用示例：

```bash
cd examples
go run main.go
```

## 与原 PHP 项目的兼容性

本项目在设计时充分考虑了与原 PHP 项目 [Yurunsoft/ChineseUtil](https://github.com/Yurunsoft/ChineseUtil) 的兼容性：

1. **功能对等**: 实现了原项目的所有主要功能
2. **API 设计**: 参考原项目的 API 设计，便于迁移
3. **内嵌数据**: 使用完整的内嵌数据，无需外部文件

### 主要差异

1. **语言特性**: 充分利用 Go 语言的类型系统和并发特性
2. **性能优化**: 针对 Go 语言进行了性能优化
3. **错误处理**: 使用 Go 语言的错误处理机制
4. **包管理**: 使用 Go Modules 进行依赖管理

## 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范

- 遵循 Go 语言编码规范
- 添加适当的测试用例
- 更新相关文档
- 确保所有测试通过

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- 感谢 [Yurunsoft/ChineseUtil](https://github.com/Yurunsoft/ChineseUtil) 项目提供的设计思路
- 感谢所有贡献者的支持

## 更新日志

### v1.0.0 (2025-06-18)

- 🎉 初始版本发布
- ✅ 实现汉字转拼音功能
- ✅ 实现拼音分词功能
- ✅ 实现简繁互转功能
- ✅ 实现数字转换功能
- ✅ 实现金额转换功能
- ✅ 内置完整数据，开箱即用
- ✅ 提供完整的测试用例
- ✅ 提供详细的文档和示例

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue: [GitHub Issues](https://github.com/arixbit/zhkit/issues)
- 邮箱: your-email@example.com

---

**ZhKit** - 让中文处理更简单！ 🚀
