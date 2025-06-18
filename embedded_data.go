package zhkit

import (
	"encoding/json"
	"fmt"

	_ "embed"
)

// 嵌入数据文件
//
//go:embed data/charsData.json
var embeddedCharsData []byte

//go:embed data/pinyinData.json
var embeddedPinyinData []byte

// loadEmbeddedData 加载嵌入的数据
func (c *Chinese) loadEmbeddedData() error {
	// 加载字符数据
	if err := c.loadEmbeddedCharsData(); err != nil {
		return fmt.Errorf("加载嵌入字符数据失败: %v", err)
	}

	// 加载拼音分词数据
	if err := c.loadEmbeddedPinyinSplitData(); err != nil {
		return fmt.Errorf("加载嵌入拼音分词数据失败: %v", err)
	}

	return nil
}

// loadEmbeddedCharsData 加载嵌入的字符数据
func (c *Chinese) loadEmbeddedCharsData() error {
	var data map[string][]interface{}
	if err := json.Unmarshal(embeddedCharsData, &data); err != nil {
		return err
	}

	// 使用现有的解析函数
	return c.parseCharsDataFormat4(data)
}

// loadEmbeddedPinyinSplitData 加载嵌入的拼音分词数据
func (c *Chinese) loadEmbeddedPinyinSplitData() error {
	var data map[string]interface{}
	if err := json.Unmarshal(embeddedPinyinData, &data); err != nil {
		return err
	}

	// 解析拼音分词数据
	if splitData, ok := data["split"].(map[string]interface{}); ok {
		if relation, ok := splitData["relation"].(map[string]interface{}); ok {
			for key, value := range relation {
				if valueStr, ok := value.(string); ok {
					// 将字符串按空格分割为切片
					if valueStr != "" {
						c.pinyinSplitData[key] = []string{valueStr}
					}
				}
			}
		}
	}

	return nil
}

// NewChineseWithFullData 创建包含完整数据的Chinese实例
func NewChineseWithFullData() *Chinese {
	c := NewChinese()
	if err := c.loadEmbeddedData(); err != nil {
		// 如果加载嵌入数据失败，返回空实例
		// 用户应该检查实例是否正常工作
	}
	return c
}

// LoadGlobalData 为全局实例加载数据
func LoadGlobalData(dataPath string) error {
	if err := defaultChinese.LoadPinyinData(dataPath); err != nil {
		return err
	}
	if err := defaultChinese.LoadSimplifiedTraditionalData(dataPath); err != nil {
		return err
	}
	if err := defaultChinese.LoadPinyinSplitData(dataPath); err != nil {
		return err
	}
	return nil
}

// LoadGlobalEmbeddedData 为全局实例加载嵌入数据
func LoadGlobalEmbeddedData() error {
	return defaultChinese.loadEmbeddedData()
}
