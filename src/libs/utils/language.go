package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var (
	languages = make(map[string]map[string]string, 0)
)

// InitLanguage unix风格的配置文件
func InitLanguage(dir, language string) map[string]string {
	// 再查找出config文件路劲
	languagePath := GetFilePath("src/"+dir, language)
	if len(languagePath) <= 0 {
		LogFatal("language", "not found language file")
	}
	// 打开config文件
	f, err := os.Open(languagePath)
	defer f.Close()
	LogFatal("language", err)
	// 开始读取config
	r := bufio.NewReader(f)
	key := ""
	for {
		// 逐行读取language文件
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break // EOF为读取完毕
			}
		}
		line := strings.TrimSpace(string(b))
		// 整行注释
		if strings.Index(line, "#") == 0 {
			continue
		}
		// 查找键值对
		index := strings.Index(line, "=")
		if index < 0 {
			continue
		}
		key = strings.TrimSpace(line[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(line[index+1:])
		// 注释处理
		pos := strings.Index(value, "\t#")
		if pos < 0 {
			pos = strings.Index(value, " #")
			if pos < 0 {
				pos = strings.Index(value, "\t//")
				if pos < 0 {
					pos = strings.Index(value, " //")
				}
			}
		}
		if pos >= 0 {
			value = strings.TrimSpace(value[0:pos])
		}
		// 无效配置，也给加上
		if len(value) == 0 {
			value = ""
		}
		// 创建language对象
		if languages[language] == nil {
			languages[language] = make(map[string]string, 0)
		}
		// 加入配置信息
		languages[language][key] = strings.TrimSpace(value)
	}
	return languages[language]
}

func GetLanguage(language, key string) string {
	if len(language) <= 0 || languages[language] == nil || len(languages[language]) <= 0 {
		language = "zh-cn"
	}
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" || language == "zh" {
		language = "zh-cn"
	}
	value := languages[language][key]
	if len(value) <= 0 {
		value = key
	}
	return value
}
