package utils

import (
	//"bufio"
	//"errors"
	//"io"
	"os"
	"path/filepath"
	//"regexp"
)

const (
	B  int64 = 1 << (iota * 10) // 1B，位运算
	KB                          // 1*1024 B
	MB                          // 1*1024*1024 B
	GB                          // 1*1024*1024*1024 B
	TB                          // 1*1024*1024*1024*1024 B
)

type (
// FileChunk 文件片定义
//FileChunk struct {
//	Number int   // 块序号
//	Offset int64 // 块在文件中的偏移量
//	Size   int64 // 块大小
//}
)

// FileExists reports whether the named file or directory exists.
func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 根据目录和文件名 获取文件真实绝对路径
func GetFilePath(dir string, name string) string {
	workRootPath, _ := os.Getwd()
	// 再查找出config文件路劲
	filePath := filepath.Join(workRootPath, dir, name)
	if !FileExists(filePath) {
		AppPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		filePath = filepath.Join(AppPath, dir, name)
		if !FileExists(filePath) {
			return ""
		}
	}
	return filePath
}

// SplitFileByChunkSize Split big file to part by the size of chunk
// 按块大小分割文件。返回值FileChunk为分割结果，error为nil时有效。
//func SplitFileByChunkSize(fileSize int64, chunkSize int64) []FileChunk {
//	var chunkN = fileSize / chunkSize //整除块数
//	var chunks []FileChunk
//	var chunk = FileChunk{}
//	for i := int64(0); i < chunkN; i++ { //整除块数的添加
//		chunk.Number = int(i + 1)
//		chunk.Offset = i * chunkSize
//		chunk.Size = chunkSize
//		chunks = append(chunks, chunk)
//	}
//	if fileSize%chunkSize > 0 { //最后一块，未整除的
//		chunk.Number = len(chunks) + 1
//		chunk.Offset = int64(len(chunks)) * chunkSize
//		chunk.Size = fileSize % chunkSize
//		chunks = append(chunks, chunk)
//	}
//	return chunks
//}

// SearchFile Search a file in paths.
// this is often used in search config file in /etc ~/
//func SearchFile(filename string, paths ...string) (fullPath string, err error) {
//	for _, path := range paths {
//		if fullPath = filepath.Join(path, filename); FileExists(fullPath) {
//			return
//		}
//	}
//	err = errors.New(fullPath + " not found in paths")
//	return
//}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
//func GrepFile(patten string, filename string) (lines []string, err error) {
//	re, err := regexp.Compile(patten)
//	if err != nil {
//		return
//	}
//
//	fd, err := os.Open(filename)
//	if err != nil {
//		return
//	}
//	lines = make([]string, 0)
//	reader := bufio.NewReader(fd)
//	prefix := ""
//	var isLongLine bool
//	for {
//		byteLine, isPrefix, er := reader.ReadLine()
//		if er != nil && er != io.EOF {
//			return nil, er
//		}
//		if er == io.EOF {
//			break
//		}
//		line := string(byteLine)
//		if isPrefix {
//			prefix += line
//			continue
//		} else {
//			isLongLine = true
//		}
//
//		line = prefix + line
//		if isLongLine {
//			prefix = ""
//		}
//		if re.MatchString(line) {
//			lines = append(lines, line)
//		}
//	}
//	return lines, nil
//}
