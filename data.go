package zhkit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CharData 字符数据结构
type CharData struct {
	Char        string   `json:"char"`
	Pinyin      []string `json:"pinyin,omitempty"`
	Simplified  []string `json:"simplified,omitempty"`
	Traditional []string `json:"traditional,omitempty"`
}

// DataLoader 数据加载器
type DataLoader struct {
	dataPath string
}

// NewDataLoader 创建数据加载器
func NewDataLoader(dataPath string) *DataLoader {
	return &DataLoader{
		dataPath: dataPath,
	}
}

// LoadFromJSON 从JSON文件加载数据
func (dl *DataLoader) LoadFromJSON(filename string) (map[string]*CharData, error) {
	filePath := filepath.Join(dl.dataPath, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	data := make(map[string]*CharData)
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("解析JSON文件失败: %v", err)
	}

	return data, nil
}

// LoadPinyinData 加载拼音数据
func (c *Chinese) LoadPinyinData(dataPath string) error {
	// 尝试加载charsData.json（兼容原PHP项目）
	if err := c.loadCharsDataJSON(filepath.Join(dataPath, "charsData.json")); err == nil {
		return nil
	}

	// 尝试加载JSON格式的数据
	loader := NewDataLoader(dataPath)
	if data, err := loader.LoadFromJSON("pinyin.json"); err == nil {
		return c.loadPinyinFromCharData(data)
	}

	return fmt.Errorf("无法加载拼音数据，请确保数据文件存在")
}

// LoadSimplifiedTraditionalData 加载简繁转换数据
func (c *Chinese) LoadSimplifiedTraditionalData(dataPath string) error {
	// 尝试加载charsData.json（兼容原PHP项目）
	if err := c.loadCharsDataJSON(filepath.Join(dataPath, "charsData.json")); err == nil {
		return nil
	}

	// 尝试加载JSON格式的数据
	loader := NewDataLoader(dataPath)
	if data, err := loader.LoadFromJSON("simplified_traditional.json"); err == nil {
		return c.loadSimplifiedTraditionalFromCharData(data)
	}

	return fmt.Errorf("无法加载简繁转换数据，请确保数据文件存在")
}

// loadPinyinFromCharData 从CharData加载拼音数据
func (c *Chinese) loadPinyinFromCharData(data map[string]*CharData) error {
	for char, charData := range data {
		if len(charData.Pinyin) > 0 {
			runes := []rune(char)
			if len(runes) == 1 {
				c.pinyinData[runes[0]] = charData.Pinyin
			}
		}
	}
	return nil
}

// loadSimplifiedTraditionalFromCharData 从CharData加载简繁转换数据
func (c *Chinese) loadSimplifiedTraditionalFromCharData(data map[string]*CharData) error {
	for char, charData := range data {
		runes := []rune(char)
		if len(runes) == 1 {
			charRune := runes[0]

			// 加载简体字数据
			if len(charData.Simplified) > 0 {
				simplifiedRunes := make([]rune, 0, len(charData.Simplified))
				for _, s := range charData.Simplified {
					sRunes := []rune(s)
					if len(sRunes) == 1 {
						simplifiedRunes = append(simplifiedRunes, sRunes[0])
					}
				}
				if len(simplifiedRunes) > 0 {
					c.simplifiedData[charRune] = simplifiedRunes
				}
			}

			// 加载繁体字数据
			if len(charData.Traditional) > 0 {
				traditionalRunes := make([]rune, 0, len(charData.Traditional))
				for _, t := range charData.Traditional {
					tRunes := []rune(t)
					if len(tRunes) == 1 {
						traditionalRunes = append(traditionalRunes, tRunes[0])
					}
				}
				if len(traditionalRunes) > 0 {
					c.traditionalData[charRune] = traditionalRunes
				}
			}
		}
	}
	return nil
}

// loadCharsDataJSON 加载charsData.json文件（兼容原PHP项目）
func (c *Chinese) loadCharsDataJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开charsData.json文件: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取charsData.json文件失败: %v", err)
	}

	// 尝试解析为不同的JSON格式
	// 格式1: {"字符": {"pinyin": ["拼音"], "simplified": ["简体"], "traditional": ["繁体"]}}
	var format1 map[string]map[string]interface{}
	if err := json.Unmarshal(content, &format1); err == nil {
		return c.parseCharsDataFormat1(format1)
	}

	// 格式2: {"字符": ["拼音1", "拼音2"]}
	var format2 map[string][]string
	if err := json.Unmarshal(content, &format2); err == nil {
		return c.parseCharsDataFormat2(format2)
	}

	// 格式3: {"字符": "拼音"}
	var format3 map[string]string
	if err := json.Unmarshal(content, &format3); err == nil {
		return c.parseCharsDataFormat3(format3)
	}

	// 格式4: {"字符": ["拼音", "简体", "繁体", 数字1, 数字2]} (原PHP项目格式)
	var format4 map[string][]interface{}
	if err := json.Unmarshal(content, &format4); err == nil {
		return c.parseCharsDataFormat4(format4)
	}

	return fmt.Errorf("无法解析charsData.json文件格式")
}

// parseCharsDataFormat1 解析格式1的charsData.json
func (c *Chinese) parseCharsDataFormat1(data map[string]map[string]interface{}) error {
	for char, charInfo := range data {
		runes := []rune(char)
		if len(runes) != 1 {
			continue
		}
		charRune := runes[0]

		// 解析拼音
		if pinyinInterface, exists := charInfo["pinyin"]; exists {
			if pinyinSlice, ok := pinyinInterface.([]interface{}); ok {
				pinyins := make([]string, 0, len(pinyinSlice))
				for _, py := range pinyinSlice {
					if pyStr, ok := py.(string); ok {
						pinyins = append(pinyins, pyStr)
					}
				}
				if len(pinyins) > 0 {
					c.pinyinData[charRune] = pinyins
				}
			}
		}

		// 解析简体字
		if simplifiedInterface, exists := charInfo["simplified"]; exists {
			if simplifiedSlice, ok := simplifiedInterface.([]interface{}); ok {
				simplified := make([]rune, 0, len(simplifiedSlice))
				for _, s := range simplifiedSlice {
					if sStr, ok := s.(string); ok {
						sRunes := []rune(sStr)
						if len(sRunes) == 1 {
							simplified = append(simplified, sRunes[0])
						}
					}
				}
				if len(simplified) > 0 {
					c.simplifiedData[charRune] = simplified
				}
			}
		}

		// 解析繁体字
		if traditionalInterface, exists := charInfo["traditional"]; exists {
			if traditionalSlice, ok := traditionalInterface.([]interface{}); ok {
				traditional := make([]rune, 0, len(traditionalSlice))
				for _, t := range traditionalSlice {
					if tStr, ok := t.(string); ok {
						tRunes := []rune(tStr)
						if len(tRunes) == 1 {
							traditional = append(traditional, tRunes[0])
						}
					}
				}
				if len(traditional) > 0 {
					c.traditionalData[charRune] = traditional
				}
			}
		}
	}
	return nil
}

// parseCharsDataFormat2 解析格式2的charsData.json
func (c *Chinese) parseCharsDataFormat2(data map[string][]string) error {
	for char, pinyins := range data {
		runes := []rune(char)
		if len(runes) == 1 && len(pinyins) > 0 {
			c.pinyinData[runes[0]] = pinyins
		}
	}
	return nil
}

// parseCharsDataFormat3 解析格式3的charsData.json
func (c *Chinese) parseCharsDataFormat3(data map[string]string) error {
	for char, pinyin := range data {
		runes := []rune(char)
		if len(runes) == 1 && pinyin != "" {
			c.pinyinData[runes[0]] = []string{pinyin}
		}
	}
	return nil
}

// parseCharsDataFormat4 解析格式4的charsData.json (原PHP项目格式)
// 格式: {"字符": ["拼音", "简体", "繁体", 数字1, 数字2]}
func (c *Chinese) parseCharsDataFormat4(data map[string][]interface{}) error {
	for char, charData := range data {
		runes := []rune(char)
		if len(runes) != 1 || len(charData) < 3 {
			continue
		}
		charRune := runes[0]

		// 解析拼音 (第一个元素)
		if pinyinInterface := charData[0]; pinyinInterface != nil {
			if pinyinStr, ok := pinyinInterface.(string); ok && pinyinStr != "" {
				// 处理多个拼音，用逗号分隔
				pinyins := strings.Split(pinyinStr, ",")
				for i, py := range pinyins {
					pinyins[i] = strings.TrimSpace(py)
				}
				c.pinyinData[charRune] = pinyins
			}
		}

		// 解析简体字 (第二个元素)
		if simplifiedInterface := charData[1]; simplifiedInterface != nil {
			if simplifiedStr, ok := simplifiedInterface.(string); ok && simplifiedStr != "" {
				simplifiedRunes := []rune(simplifiedStr)
				if len(simplifiedRunes) > 0 {
					c.simplifiedData[charRune] = simplifiedRunes
				}
			}
		}

		// 解析繁体字 (第三个元素)
		if traditionalInterface := charData[2]; traditionalInterface != nil {
			if traditionalStr, ok := traditionalInterface.(string); ok && traditionalStr != "" {
				traditionalRunes := []rune(traditionalStr)
				if len(traditionalRunes) > 0 {
					c.traditionalData[charRune] = traditionalRunes
				}
			}
		}
	}
	return nil
}

// LoadPinyinSplitData 加载拼音分词数据
func (c *Chinese) LoadPinyinSplitData(dataPath string) error {
	// 尝试加载JSON格式的拼音分词数据
	filePath := filepath.Join(dataPath, "pinyin_split.json")
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开拼音分词数据文件: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c.pinyinSplitData)
	if err != nil {
		return fmt.Errorf("解析拼音分词数据失败: %v", err)
	}

	return nil
}

// SaveDataToJSON 保存数据到JSON文件
func (c *Chinese) SaveDataToJSON(dataPath string) error {
	// 保存拼音数据
	pinyinData := make(map[string][]string)
	for char, pinyins := range c.pinyinData {
		pinyinData[string(char)] = pinyins
	}

	pinyinFile := filepath.Join(dataPath, "pinyin_export.json")
	if err := c.saveJSONFile(pinyinFile, pinyinData); err != nil {
		return fmt.Errorf("保存拼音数据失败: %v", err)
	}

	// 保存简繁转换数据
	simplifiedData := make(map[string][]string)
	for char, simplified := range c.simplifiedData {
		simplifiedStrs := make([]string, len(simplified))
		for i, s := range simplified {
			simplifiedStrs[i] = string(s)
		}
		simplifiedData[string(char)] = simplifiedStrs
	}

	simplifiedFile := filepath.Join(dataPath, "simplified_export.json")
	if err := c.saveJSONFile(simplifiedFile, simplifiedData); err != nil {
		return fmt.Errorf("保存简体数据失败: %v", err)
	}

	traditionalData := make(map[string][]string)
	for char, traditional := range c.traditionalData {
		traditionalStrs := make([]string, len(traditional))
		for i, t := range traditional {
			traditionalStrs[i] = string(t)
		}
		traditionalData[string(char)] = traditionalStrs
	}

	traditionalFile := filepath.Join(dataPath, "traditional_export.json")
	if err := c.saveJSONFile(traditionalFile, traditionalData); err != nil {
		return fmt.Errorf("保存繁体数据失败: %v", err)
	}

	// 保存拼音分词数据
	pinyinSplitFile := filepath.Join(dataPath, "pinyin_split_export.json")
	if err := c.saveJSONFile(pinyinSplitFile, c.pinyinSplitData); err != nil {
		return fmt.Errorf("保存拼音分词数据失败: %v", err)
	}

	return nil
}

// saveJSONFile 保存数据到JSON文件
func (c *Chinese) saveJSONFile(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
