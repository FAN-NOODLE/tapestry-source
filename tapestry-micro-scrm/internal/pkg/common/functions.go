package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/shakinm/xlsReader/xls"

	"github.com/tealeg/xlsx"
	"micro-scrm/internal/config"
)

func RemoveDuplicate(v []int64) []int64 {
	//为了性能需要尽可能的减小拷贝，最悲观的情况每个元素只移动一次。
	toIndex := 0
	var p int64 = 0
	for i, _ := range v {
		// 为了实际去重结构时减小内存拷贝
		c := &v[i]
		if p == *c && i != 0 {
			// 重复内容，跳过
			continue
		}
		if i != toIndex {
			// 需要移动当前元素
			v[toIndex] = *c
		}
		toIndex++
		p = *c
	}
	return v[:toIndex]
}

// 获取当前日期
func GetDate(sep string) string {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	d := strconv.Itoa(t.Day())

	return y + sep + m + sep + d
}

// GetTmpFileOfURL 获取url指定的内容，生成临时文件
func GetTmpFileOfURL(URL string) (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	res, err := client.Get(URL)
	if err != nil {
		return "", err
	}

	urlBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	//defer res.Body.Close()
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10)

	FilePath := config.CfgData.UploadPath
	ext := path.Ext(URL)
	// URL可能带有query参数，其中包含"/"会导致path.Ext解析错误，需要先解析URL取path部分
	if u, err := url.Parse(URL); err == nil && u.Path != "" {
		ext = path.Ext(u.Path)
	}
	file := FilePath + fileName + ext

	err = ioutil.WriteFile(file, []byte(urlBody), 0777)
	if err != nil {
		return "", err
	}
	return file, nil
}

// removeTmpImgFile 删除临时文件
func RemoveTmpImgFile(file string) error {
	return os.Remove(file)
}

//文本加密
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
var key = "aaabbbcccxaaabbbcccxaaabbbcccx88"

func AesEncrypt(tex string) (string, error) {
	//需要去加密的字符串
	plaintext := []byte(tex)
	// 创建加密算法aes
	c, err := newCipher(key)
	if err != nil {
		return "error", err
	}
	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return hex.EncodeToString(ciphertext), nil
}

func AesDecrypt(enText string) (string, error) {
	enByte, err := hex.DecodeString(enText)
	if err != nil {
		return "error", err
	}
	// 创建加密算法aes
	c, cipherErr := newCipher(key)
	if cipherErr != nil {
		return "error", cipherErr
	}
	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(enByte))
	cfbdec.XORKeyStream(plaintextCopy, enByte)
	return string(plaintextCopy), nil
}

//创建加密算法aes
func newCipher(key string) (c cipher.Block, err error) {
	// 创建加密算法aes
	c, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	return
}

// URLAddQuery url链接拼接query参数
func URLAddQuery(URL string, query string) string {
	// 解析小程序链接
	urlObj, err := url.Parse(URL)
	if err != nil {
		return URL
	}
	// url中拼接参数
	if urlObj.RawQuery == "" {
		urlObj.RawQuery = query
	} else {
		urlObj.RawQuery = fmt.Sprintf("%s&%s", urlObj.RawQuery, query)
	}
	return urlObj.String()
}

// ReadExcelFile 读取表单上传的excel文件内容，支持 .xlsx、.xls格式的文件
// 如果上传的文件格式是.xlsx，则返回*xlsx.File，否则返回*xls.Workbook
func ReadExcelFile(fh *multipart.FileHeader, ft *string) (interface{}, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	if filepath.Ext(fh.Filename) == ".xlsx" {
		*ft = ".xlsx"
		fBody, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return xlsx.OpenBinary(fBody)
	}

	*ft = ".xls"
	wb, err := xls.OpenReader(f)
	if err != nil {
		return nil, err
	}
	return &wb, nil
}
