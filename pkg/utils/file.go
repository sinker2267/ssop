package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// SaveFile 保存上传的文件到指定路径
func SaveFile(file multipart.File, dst string) (int64, error) {
	// 确保目标目录存在
	dir := filepath.Dir(dst)
	if err := EnsureDir(dir); err != nil {
		return 0, err
	}
	
	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	
	// 复制文件内容
	return io.Copy(out, file)
}

// GetFileExt 获取文件扩展名
func GetFileExt(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// IsAllowedFileType 检查文件类型是否在允许列表中
func IsAllowedFileType(filename string, allowedExts []string) bool {
	ext := GetFileExt(filename)
	for _, allowed := range allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
} 