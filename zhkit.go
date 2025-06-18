// Package zhkit 提供中文处理工具包
// 支持汉字转拼音、拼音分词、简繁互转、数字转换、金额转换等功能
package zhkit

import (
	"strings"
)

// ConvertMode 转换模式
type ConvertMode int

const (
	// ModePinyin 全拼模式
	ModePinyin ConvertMode = 1 << iota
	// ModePinyinFirst 首字母模式
	ModePinyinFirst
	// ModePinyinSound 读音模式（带声调符号）
	ModePinyinSound
	// ModePinyinSoundNumber 读音数字模式（数字声调）
	ModePinyinSoundNumber
)

// PinyinResult 拼音转换结果
type PinyinResult struct {
	Pinyin            [][]string `json:"pinyin,omitempty"`
	PinyinFirst       [][]string `json:"pinyinFirst,omitempty"`
	PinyinSound       [][]string `json:"pinyinSound,omitempty"`
	PinyinSoundNumber [][]string `json:"pinyinSoundNumber,omitempty"`
}

// NumberOptions 数字转换选项
type NumberOptions struct {
	TenMin bool // "一十二" => "十二"
}

// Chinese 中文工具类
type Chinese struct {
	pinyinData      map[rune][]string
	simplifiedData  map[rune][]rune
	traditionalData map[rune][]rune
	pinyinSplitData map[string][]string
}

// NewChinese 创建新的中文工具实例
func NewChinese() *Chinese {
	c := &Chinese{
		pinyinData:      make(map[rune][]string),
		simplifiedData:  make(map[rune][]rune),
		traditionalData: make(map[rune][]rune),
		pinyinSplitData: make(map[string][]string),
	}
	return c
}

// ToPinyin 汉字转拼音
// text: 要转换的文本
// mode: 转换模式，可以使用位运算组合多种模式
// separator: 分隔符，默认为空格
// splitNonChinese: 是否分割非中文字符
func (c *Chinese) ToPinyin(text string, mode ConvertMode, separator string, splitNonChinese bool) (*PinyinResult, error) {
	if text == "" {
		return &PinyinResult{}, nil
	}

	if separator == "" {
		separator = " "
	}

	runes := []rune(text)
	result := &PinyinResult{}

	// 根据模式初始化结果数组
	if mode&ModePinyin != 0 {
		result.Pinyin = make([][]string, 0)
	}
	if mode&ModePinyinFirst != 0 {
		result.PinyinFirst = make([][]string, 0)
	}
	if mode&ModePinyinSound != 0 {
		result.PinyinSound = make([][]string, 0)
	}
	if mode&ModePinyinSoundNumber != 0 {
		result.PinyinSoundNumber = make([][]string, 0)
	}

	// 处理每个字符
	for _, r := range runes {
		if pinyins, exists := c.pinyinData[r]; exists {
			// 中文字符
			if mode&ModePinyin != 0 {
				result.Pinyin = append(result.Pinyin, pinyins)
			}
			if mode&ModePinyinFirst != 0 {
				firsts := make([]string, len(pinyins))
				for i, py := range pinyins {
					if len(py) > 0 {
						firsts[i] = string(py[0])
					}
				}
				result.PinyinFirst = append(result.PinyinFirst, firsts)
			}
			if mode&ModePinyinSound != 0 {
				sounds := make([]string, len(pinyins))
				for i, py := range pinyins {
					sounds[i] = c.addToneMarks(py)
				}
				result.PinyinSound = append(result.PinyinSound, sounds)
			}
			if mode&ModePinyinSoundNumber != 0 {
				numbers := make([]string, len(pinyins))
				for i, py := range pinyins {
					numbers[i] = c.addToneNumbers(py)
				}
				result.PinyinSoundNumber = append(result.PinyinSoundNumber, numbers)
			}
		} else {
			// 非中文字符
			charStr := string(r)
			if mode&ModePinyin != 0 {
				result.Pinyin = append(result.Pinyin, []string{charStr})
			}
			if mode&ModePinyinFirst != 0 {
				result.PinyinFirst = append(result.PinyinFirst, []string{charStr})
			}
			if mode&ModePinyinSound != 0 {
				result.PinyinSound = append(result.PinyinSound, []string{charStr})
			}
			if mode&ModePinyinSoundNumber != 0 {
				result.PinyinSoundNumber = append(result.PinyinSoundNumber, []string{charStr})
			}
		}
	}

	return result, nil
}

// SplitPinyin 拼音分词（返回字符串）
func (c *Chinese) SplitPinyin(pinyin string) ([]string, error) {
	if pinyin == "" {
		return []string{}, nil
	}

	pinyin = strings.ToLower(strings.TrimSpace(pinyin))
	results := c.splitPinyinRecursive(pinyin, []string{})

	if len(results) == 0 {
		return []string{pinyin}, nil
	}

	// 转换为字符串格式
	stringResults := make([]string, len(results))
	for i, result := range results {
		stringResults[i] = strings.Join(result, " ")
	}

	return stringResults, nil
}

// SplitPinyinArray 拼音分词（返回数组）
func (c *Chinese) SplitPinyinArray(pinyin string) ([][]string, error) {
	if pinyin == "" {
		return [][]string{}, nil
	}

	pinyin = strings.ToLower(strings.TrimSpace(pinyin))
	results := c.splitPinyinRecursive(pinyin, []string{})

	if len(results) == 0 {
		return [][]string{{pinyin}}, nil
	}

	return results, nil
}

// ToSimplified 繁体转简体
func (c *Chinese) ToSimplified(text string) ([]string, error) {
	if text == "" {
		return []string{""}, nil
	}

	runes := []rune(text)
	result := make([]rune, len(runes))

	for i, r := range runes {
		if simplified, exists := c.simplifiedData[r]; exists && len(simplified) > 0 {
			result[i] = simplified[0] // 取第一个简体字
		} else {
			result[i] = r // 保持原字符
		}
	}

	return []string{string(result)}, nil
}

// ToTraditional 简体转繁体
func (c *Chinese) ToTraditional(text string) ([]string, error) {
	if text == "" {
		return []string{""}, nil
	}

	runes := []rune(text)
	result := make([]rune, len(runes))

	for i, r := range runes {
		if traditional, exists := c.traditionalData[r]; exists && len(traditional) > 0 {
			result[i] = traditional[0] // 取第一个繁体字
		} else {
			result[i] = r // 保持原字符
		}
	}

	return []string{string(result)}, nil
}

// addToneMarks 添加声调符号
func (c *Chinese) addToneMarks(pinyin string) string {
	// 简化实现，实际需要根据声调规则添加声调符号
	return pinyin
}

// addToneNumbers 添加数字声调
func (c *Chinese) addToneNumbers(pinyin string) string {
	// 简化实现，实际需要根据声调规则添加数字声调
	return pinyin + "1"
}

// splitPinyinRecursive 递归分词
func (c *Chinese) splitPinyinRecursive(pinyin string, current []string) [][]string {
	if pinyin == "" {
		return [][]string{current}
	}

	var results [][]string

	// 尝试不同长度的拼音
	for i := 1; i <= len(pinyin) && i <= 6; i++ {
		prefix := pinyin[:i]
		if c.isValidPinyin(prefix) {
			newCurrent := append(current, prefix)
			subResults := c.splitPinyinRecursive(pinyin[i:], newCurrent)
			results = append(results, subResults...)
		}
	}

	return results
}

// isValidPinyin 检查是否为有效拼音
func (c *Chinese) isValidPinyin(pinyin string) bool {
	// 简化实现，实际需要检查拼音字典
	validPinyins := []string{
		"a", "ai", "an", "ang", "ao",
		"ba", "bai", "ban", "bang", "bao", "bei", "ben", "beng", "bi", "bian", "biao", "bie", "bin", "bing", "bo", "bu",
		"ca", "cai", "can", "cang", "cao", "ce", "cen", "ceng", "cha", "chai", "chan", "chang", "chao", "che", "chen", "cheng", "chi", "chong", "chou", "chu", "chuai", "chuan", "chuang", "chui", "chun", "chuo", "ci", "cong", "cou", "cu", "cuan", "cui", "cun", "cuo",
		"da", "dai", "dan", "dang", "dao", "de", "deng", "di", "dian", "diao", "die", "ding", "diu", "dong", "dou", "du", "duan", "dui", "dun", "duo",
		"e", "en", "er",
		"fa", "fan", "fang", "fei", "fen", "feng", "fo", "fou", "fu",
		"ga", "gai", "gan", "gang", "gao", "ge", "gei", "gen", "geng", "gong", "gou", "gu", "gua", "guai", "guan", "guang", "gui", "gun", "guo",
		"ha", "hai", "han", "hang", "hao", "he", "hei", "hen", "heng", "hong", "hou", "hu", "hua", "huai", "huan", "huang", "hui", "hun", "huo",
		"ji", "jia", "jian", "jiang", "jiao", "jie", "jin", "jing", "jiong", "jiu", "ju", "juan", "jue", "jun",
		"ka", "kai", "kan", "kang", "kao", "ke", "ken", "keng", "kong", "kou", "ku", "kua", "kuai", "kuan", "kuang", "kui", "kun", "kuo",
		"la", "lai", "lan", "lang", "lao", "le", "lei", "leng", "li", "lia", "lian", "liang", "liao", "lie", "lin", "ling", "liu", "long", "lou", "lu", "luan", "lue", "lun", "luo", "lv",
		"ma", "mai", "man", "mang", "mao", "me", "mei", "men", "meng", "mi", "mian", "miao", "mie", "min", "ming", "miu", "mo", "mou", "mu",
		"na", "nai", "nan", "nang", "nao", "ne", "nei", "nen", "neng", "ni", "nian", "niang", "niao", "nie", "nin", "ning", "niu", "nong", "nu", "nuan", "nue", "nuo", "nv",
		"o", "ou",
		"pa", "pai", "pan", "pang", "pao", "pei", "pen", "peng", "pi", "pian", "piao", "pie", "pin", "ping", "po", "pou", "pu",
		"qi", "qia", "qian", "qiang", "qiao", "qie", "qin", "qing", "qiong", "qiu", "qu", "quan", "que", "qun",
		"ran", "rang", "rao", "re", "ren", "reng", "ri", "rong", "rou", "ru", "ruan", "rui", "run", "ruo",
		"sa", "sai", "san", "sang", "sao", "se", "sen", "seng", "sha", "shai", "shan", "shang", "shao", "she", "shen", "sheng", "shi", "shou", "shu", "shua", "shuai", "shuan", "shuang", "shui", "shun", "shuo", "si", "song", "sou", "su", "suan", "sui", "sun", "suo",
		"ta", "tai", "tan", "tang", "tao", "te", "teng", "ti", "tian", "tiao", "tie", "ting", "tong", "tou", "tu", "tuan", "tui", "tun", "tuo",
		"wa", "wai", "wan", "wang", "wei", "wen", "weng", "wo", "wu",
		"xi", "xia", "xian", "xiang", "xiao", "xie", "xin", "xing", "xiong", "xiu", "xu", "xuan", "xue", "xun",
		"ya", "yan", "yang", "yao", "ye", "yi", "yin", "ying", "yo", "yong", "you", "yu", "yuan", "yue", "yun",
		"za", "zai", "zan", "zang", "zao", "ze", "zei", "zen", "zeng", "zha", "zhai", "zhan", "zhang", "zhao", "zhe", "zhen", "zheng", "zhi", "zhong", "zhou", "zhu", "zhua", "zhuai", "zhuan", "zhuang", "zhui", "zhun", "zhuo", "zi", "zong", "zou", "zu", "zuan", "zui", "zun", "zuo",
	}

	for _, valid := range validPinyins {
		if valid == pinyin {
			return true
		}
	}
	return false
}

// 全局实例 - 使用完整数据
var defaultChinese = NewChineseWithFullData()

// ToPinyin 全局函数：汉字转拼音
func ToPinyin(text string, mode ConvertMode, separator string, splitNonChinese bool) (*PinyinResult, error) {
	return defaultChinese.ToPinyin(text, mode, separator, splitNonChinese)
}

// SplitPinyin 全局函数：拼音分词（返回字符串）
func SplitPinyin(pinyin string) ([]string, error) {
	return defaultChinese.SplitPinyin(pinyin)
}

// SplitPinyinArray 全局函数：拼音分词（返回数组）
func SplitPinyinArray(pinyin string) ([][]string, error) {
	return defaultChinese.SplitPinyinArray(pinyin)
}

// ToSimplified 全局函数：繁体转简体
func ToSimplified(text string) ([]string, error) {
	return defaultChinese.ToSimplified(text)
}

// ToTraditional 全局函数：简体转繁体
func ToTraditional(text string) ([]string, error) {
	return defaultChinese.ToTraditional(text)
}
