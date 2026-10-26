package common

import (
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const DefaultSplitSize = 1000

// PageCount 计算分页总数
func PageCount(count int64) int64 {
	pageCount := count / DefaultSplitSize
	if count%DefaultSplitSize > 0 {
		pageCount++
	}
	if pageCount == 0 {
		return 1
	}
	return pageCount
}

// 判断字符串是否在数组中
func StringInArray(value string, array []string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

// 判断字符串是否在数组中
func AnyInArray(value interface{}, array []interface{}) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

// 判断数字是否在数组中
func IntInArray(value int, array []int) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

// 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 获取httpReader
func GetResponseReader(encoding string, closer io.ReadCloser) (io.Reader, error) {
	// 判断响应是否经过压缩
	var bodyReader io.Reader
	switch encoding {
	case "gzip":
		itemReader, err := gzip.NewReader(closer)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("gzip响应读取出错:%s", err.Error()))
		}
		bodyReader = itemReader
		break
	case "deflate":
		bodyReader = flate.NewReader(closer)
		break
	default:
		bodyReader = closer
	}
	return bodyReader, nil
}

func ZipDir(dir, zipFile string) error {
	fz, err := os.Create(zipFile)
	if err != nil {
		log.Fatalf("创建zip文件失败: %s\n", err.Error())
	}
	defer func() { _ = fz.Close() }()
	w := zip.NewWriter(fz)
	defer w.Close()
	return filepath.Walk(
		dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				fDest, err := w.Create(path[len(dir)+1:])
				if err != nil {
					log.Printf("Zip添加文件失败: %s\n", err.Error())
					return nil
				}
				fSrc, err := os.Open(path)
				if err != nil {
					log.Printf("Zip打开文件失败: %s\n", err.Error())
					return nil
				}
				defer fSrc.Close()
				_, err = io.Copy(fDest, fSrc)
				if err != nil {
					log.Printf("Zip复制文件失败: %s\n", err.Error())
					return nil
				}
			}
			return nil
		},
	)
}
