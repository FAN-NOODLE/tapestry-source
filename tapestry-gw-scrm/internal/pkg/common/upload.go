package common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

//UploadCfg 上传设置
type UploadCfg struct {
}

//UploadFile 上传文件
func (receiver UploadCfg) UploadFile(c *gin.Context, path string) (string, error) {
	header, _ := c.FormFile("file")
	fileName := header.Filename
	var build strings.Builder
	build.WriteString(path)
	build.WriteString(fileName)
	dst := build.String()
	if err := c.SaveUploadedFile(header, dst); err != nil {
		return "", err
	}
	return dst, nil
}
