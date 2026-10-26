package file

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
	"micro-download/internal/config"
	"micro-download/internal/modules/download/delivery"
	"micro-download/internal/pkg/blob"
	"micro-download/internal/pkg/cos"
)

const RootPath = "download"

type GenerateFile interface {
	Generate(c ctx.Context, filename string, startLine int, data [][]any, overview *delivery.SourceOverview, complete bool) error
}

func factory(fileType int8) GenerateFile {
	var once sync.Once
	// var excelObject *ExcelGenerate
	var csvObject *CsvGenerate
	var zipObject *ZipGenerate
	var xlsObject *XlsxGenerate
	switch fileType {
	case delivery.TypeExcel:
		once.Do(func() {
			// excelObject = &ExcelGenerate{}
			xlsObject = &XlsxGenerate{}
		})
		return xlsObject
	case delivery.TypeCsv:
		once.Do(func() {
			csvObject = &CsvGenerate{}
		})
		return csvObject
	case delivery.TypeZip:
		once.Do(func() {
			zipObject = &ZipGenerate{}
		})
		return zipObject
	}
	return nil
}

// GenerateFileName 统一生成下载的名称
func GenerateFileName(name string) string {
	regex, err := regexp.Compile("[{}\\[\\]<>|\\^~`#%]")
	if nil == err {
		name = string(regex.ReplaceAll([]byte(name), []byte("_")))
	}

	return name + "-" + time.Now().Format("20060102150405")
}

// Generate 返回错误和临时文件路径
func Generate(c ctx.Context, fileType int8, filename string, startLine int, data [][]any, overview *delivery.SourceOverview, complete bool) error {
	object := factory(fileType)
	if object != nil {
		return object.Generate(c, filename, startLine, data, overview, complete)
	}
	return errors.New("临时文件生成失败")
}

// UploadCos 上传到COS返回cos存储地址
func UploadCos(c ctx.Context, filePath string) (error, string) {
	client := cos.Cos()
	now := time.Now()
	cosFilename := fmt.Sprintf("%s/%d/%02d/%02d%s", RootPath, now.Year(), now.Month(), now.Day(), strings.Replace(filePath, config.CfgData.Downloader.TmpFilePath, "", 1))
	putResponse, err := client.Object.PutFromFile(c, cosFilename, filePath, nil)
	if err != nil {
		return err, ""
	}
	if putResponse.StatusCode != 200 {
		return errors.New("上传遇到未知错误"), ""
	}
	// 上传完毕删除文件
	_ = os.Remove(filePath)

	return nil, cosFilename
}

// UploadBlob 上传到blob返回blob存储地址
func UploadBlob(c ctx.Context, filePath string) (err error, path string) {
	blobUpload := new(blob.Blob)
	now := time.Now()
	pixFileName := fmt.Sprintf("%s/%d/%02d/%02d%s", RootPath, now.Year(), now.Month(), now.Day(), strings.Replace(filePath, config.CfgData.Downloader.TmpFilePath, "", 1))
	path, err = blobUpload.Upload(c, filePath, pixFileName)
	if err != nil {
		log.Info(c).Msgf("uploadCloud Error: %+v", err)
		return
	}
	// 上传完毕删除文件
	_ = os.Remove(filePath)

	return nil, path
}

// GeneratePathFileName 生成文件带路径文件名称
func GeneratePathFileName(name string, fileType int8) string {
	var suffix string
	switch fileType {
	case delivery.TypeExcel:
		suffix = "xlsx"
		break
	case delivery.TypeCsv:
		suffix = "csv"
		break
	case delivery.TypeZip:
		suffix = "zip"
		break
	}
	return config.CfgData.Downloader.TmpFilePath + "/" + name + "." + suffix
}
