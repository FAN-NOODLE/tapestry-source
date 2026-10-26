package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gw-scrm/internal/config"
	"gw-scrm/internal/pkg/cosroute"
	"net/http"
	"net/url"
)

const COS = "cos"

type Storage struct{}

// 上传腾讯云
func (cs *Storage) Upload(ctx context.Context, filePath string, fileName string, toPath string) (string, error) {
	toPath += "/"
	objectKey := toPath + fileName
	bucket := cosroute.BucketForKey(objectKey)

	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosroute.OriginHost(bucket))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})

	// 通过本地文件上传对象
	_, err := c.Object.PutFromFile(ctx, objectKey, filePath, nil)
	if err != nil {
		return "", err
	}

	return objectKey, nil
}

// 获取存储域
func (cs *Storage) GetDomain() string {
	return config.CfgData.Storage.COS.Schema + "://" + config.CfgData.Storage.COS.CdnHost
}

func (cs *Storage) GetURL(objectKey string) string {
	return cosroute.URLForKey(objectKey)
}
