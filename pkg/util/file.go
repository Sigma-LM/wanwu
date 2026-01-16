package util

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/UnicomAI/wanwu/pkg/log"
)

const (
	kb               = 1024
	mb               = kb * 1024
	MaxScanTokenSize = 1024 * 1024 // Set the maximum token size to 1 MB
)

var specialFileExtList = []string{".tar.gz"}

type FileMergeResult struct {
	TotalSuccessCount int64
	TotalLineCount    int64
	TotalByteCount    int64
	FilePath          string
}

func FileExt(filePath string) string {
	if len(filePath) == 0 {
		return ""
	}
	for _, ext := range specialFileExtList {
		if strings.HasSuffix(filePath, ext) {
			return ext
		}
	}
	return filepath.Ext(filePath)
}

// ToFileSizeStr fileSize单位是B，转换规则：小于1M为KB，大于等于1M，单位为M，保留两位小数
func ToFileSizeStr(fileSize int64) string {
	if fileSize < mb {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(kb))
	} else {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(mb))
	}
}

func FileExist(filePath string) (bool, error) {
	if len(filePath) == 0 {
		return false, nil
	}
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func DirFileList(dir string, subDir bool, fullPath bool) ([]string, error) {
	var fileNameList []string
	// 读取目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir (%v) err: %v", dir, err)
	}

	// 遍历目录下的所有文件和子目录
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// 处理错误
			log.Errorf("read dir (%v) entry err: %v", dir, err)
			continue
		}

		// 判断是否是文件
		if !info.IsDir() {
			if fullPath {
				fileNameList = append(fileNameList, dir+"/"+entry.Name())
			} else {
				fileNameList = append(fileNameList, entry.Name())
			}
		} else if !subDir { //不需要校验底层目录
			continue
		} else {
			list, err := DirFileList(dir+"/"+entry.Name(), subDir, fullPath)
			if err != nil {
				return nil, err
			} else {
				fileNameList = append(fileNameList, list...)
			}
		}
	}

	return fileNameList, nil
}

// MergeFile 合并文件
func MergeFile(filePathList []string, mergeFilePath string) (*FileMergeResult, error) {
	// 创建或打开文件
	//0644，表示文件所有者可读写，同组用户及其他用户只可读
	dir := filepath.Dir(mergeFilePath)
	exist, err := FileExist(dir)
	if err != nil {
		return nil, err
	}
	if !exist {
		err = os.MkdirAll(filepath.Dir(mergeFilePath), 0755)
		if err != nil {
			return nil, err
		}
	}

	destinationFile, err := os.OpenFile(mergeFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open merge file (%v) err: %v", mergeFilePath, err)
	}
	defer func() {
		if err := destinationFile.Close(); err != nil {
			log.Errorf("close merge file (%v) err: %v", mergeFilePath, err)
		}
	}()

	var totalByteCount int64
	for _, fileInfo := range filePathList {
		byteCount, err := AppendFileStream(fileInfo, destinationFile)
		if err != nil {
			return nil, fmt.Errorf("merge file (%v) err: %v", mergeFilePath, err)
		}
		totalByteCount += byteCount
	}
	return &FileMergeResult{
		TotalByteCount: totalByteCount,
		FilePath:       mergeFilePath,
	}, nil
}

func DeleteDirFile(fileDir string) error {
	err := os.RemoveAll(fileDir)
	if err != nil {
		return fmt.Errorf("delete dir (%v) err: %v", fileDir, err)
	}
	return nil
}

func DeleteFile(file string) error {
	err := os.Remove(file)
	if err != nil {
		return fmt.Errorf("delete file (%v) err: %v", file, err)
	}
	return nil
}

func AppendFileStream(filePath string, destinationFile *os.File) (int64, error) {
	// Open the source file for reading
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("open append file (%v) err: %v", filePath, err)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			log.Errorf("close append file (%v) err: %v", filePath, err)
		}
	}()
	fileReader := bufio.NewReader(sourceFile)
	byteCount, err := appendFile(fileReader, destinationFile)
	if err != nil {
		return 0, fmt.Errorf("append file (%v) to (%v) err: %v", filePath, destinationFile.Name(), err)
	}
	log.Infof("append file (%v) to (%v) succeed, bytes: %v", filePath, destinationFile.Name(), byteCount)
	return byteCount, nil
}

func appendFile(reader *bufio.Reader, destinationFile *os.File) (byteCount int64, error error) {
	buf := make([]byte, MaxScanTokenSize)
	for {
		n, err := reader.Read(buf)
		if FileEOF(err) { // 检查是否到达文件末尾
			break
		}
		if err != nil {
			log.Errorf("Error reading file: %s", err)
			return -1, err
		}
		line := buf[:n]
		bytesWritten, err := destinationFile.Write(line)
		if err != nil {
			log.Errorf("appendFile error %s", err)
			return -1, err
		}
		byteCount += int64(bytesWritten)
	}
	return byteCount, nil
}

func FileEOF(err error) bool {
	return errors.Is(err, io.EOF) || (err != nil && err.Error() == "EOF")
}

// Img2base64
//
//	@Description:
//	@Author zhangzekai
//	@Time 2026-01-14 17:35:29
//	@param imgPath 图片路径
//	@return string 完整的DataURI格式Base64字符串 (带 data:image/xxx;base64, 前缀，前端可直接渲染图片)
//	@return string 纯净的Base64编码内容 (无任何前缀，纯文件二进制编码结果)
//	@return error
func Img2base64(imgPath string) (string, string, error) {
	// 读取图片文件
	data, err := os.ReadFile(imgPath)
	if err != nil {
		return "", "", err
	}

	// 获取文件扩展名（不含点）
	ext := strings.TrimPrefix(filepath.Ext(imgPath), ".")

	// 对文件内容进行base64编码
	encodedImage := base64.StdEncoding.EncodeToString(data)

	// 构建完整的base64数据URI
	imgBase64Str := "data:image/" + ext + ";base64," + encodedImage
	return imgBase64Str, encodedImage, nil
}

// Img2base64ByBytes
//
//	@Description: 二进制文件流转Base64编码
//	@Author zhangzekai
//	@Time 2026-01-16 11:13:47
//	@param fileData 文件二进制字节流
//	@param ext 文件后缀(带/不带点均可)
//	@return string 带data:image/xxx;base64,前缀的完整Base64，前端可直接渲染
//	@return string 无前缀的纯净Base64编码内容
//	@return error
func Img2base64ByBytes(fileData []byte, ext string) (string, string, error) {
	if len(fileData) == 0 {
		return "", "", fmt.Errorf("文件二进制数据为空")
	}
	// 处理文件后缀，统一去除前缀点
	ext = strings.TrimPrefix(ext, ".")
	// 生成纯净Base64编码
	encodedImage := base64.StdEncoding.EncodeToString(fileData)
	// 拼接完整的DataURI格式Base64
	imgBase64Str := "data:image/" + ext + ";base64," + encodedImage
	return imgBase64Str, encodedImage, nil
}
