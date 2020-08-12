package utils

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	configs = make(map[string]map[string]map[string]string, 0)
)

// InitConfig unix风格的配置文件
func InitConfig(dir, name string) map[string]map[string]string {
	// 再查找出config文件路劲
	configPath := GetFilePath("src/"+dir, name)
	if len(configPath) <= 0 {
		LogFatal("config", "not found config file")
	}
	// 打开config文件
	f, err := os.Open(configPath)
	defer f.Close()
	LogFatal("config", err)
	// 开始读取config
	r := bufio.NewReader(f)
	node := ""
	for {
		// 逐行读取config文件
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
		// 节点配置
		n1 := strings.Index(line, "[")
		n2 := strings.LastIndex(line, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			// 有新节点则更新当前节点
			node = strings.TrimSpace(line[n1+1 : n2])
			continue
		}
		// 节点配置不能为空
		if len(node) == 0 {
			continue
		}
		// 查找节点下的键值对
		index := strings.Index(line, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(line[:index])
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
		// 创建config对象
		if configs[dir+"/"+name] == nil {
			configs[dir+"/"+name] = make(map[string]map[string]string)
			//config = &Config{
			//	configs: make(map[string]map[string]string),
			//	path:    configPath,
			//}
		}
		//if config.configs == nil {
		//	config.configs = make(map[string]map[string]string)
		//}
		if configs[dir+"/"+name][node] == nil {
			configs[dir+"/"+name][node] = make(map[string]string)
		}
		// 加入配置信息
		configs[dir+"/"+name][node][key] = strings.TrimSpace(value)
	}
	return configs[dir+"/"+name]
}

func GetConfigStr(dir, name, node, key string) string {
	value := configs[dir+"/"+name][node][key]
	return value
}

func GetConfigInt(dir, name, node, key string) int {
	value := GetConfigStr(dir, name, node, key)
	i, _ := strconv.Atoi(value)
	return i
}

func GetConfigInt64(dir, name, node, key string) int64 {
	value := GetConfigStr(dir, name, node, key)
	i, _ := strconv.ParseInt(value, 10, 64)
	return i
}

func GetConfigFloat64(dir, name, node, key string) float64 {
	value := GetConfigStr(dir, name, node, key)
	i, _ := strconv.ParseFloat(value, 64)
	return i
}

func GetConfigBool(dir, name, node, key string) bool {
	value := GetConfigStr(dir, name, node, key)
	b, _ := strconv.ParseBool(value)
	return b
}
