package cos

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/cosroute"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const COS = "cos"

type Storage struct{}

// 上传腾讯云
func (cs *Storage) Upload(ctx context.Context, filePath string, toFilePath string) (string, error) {
	bucket := cosroute.BucketForKey(toFilePath)
	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosroute.OriginHost(bucket))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})

	// 通过本地文件上传对象
	// 图片类型文件不设置自动下载，其他文件设置 Content-Disposition 头部让下载的时候保持原文件名
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExtSlice := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true, ".svg": true}

	var options *cos.ObjectPutOptions
	if imageExtSlice[ext] {
		options = &cos.ObjectPutOptions{
			ACLHeaderOptions: nil,
		}
	} else {
		disposition := fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filepath.Base(filePath)))
		options = &cos.ObjectPutOptions{
			ACLHeaderOptions: nil,
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentDisposition: disposition,
			},
		}
	}
	_, err := c.Object.PutFromFile(ctx, toFilePath, filePath, options)

	if err != nil {
		return "", err
	}

	return toFilePath, nil
}

// Download 下载腾讯云 COS 文件 含私有文件
func (cs *Storage) Download(ctx context.Context, fileUrl string, toFilePath string) error {
	objectInfo, err := cosroute.ParseObjectURL(fileUrl)
	if err != nil {
		return err
	}
	// 初始化腾讯云 COS 客户端
	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosroute.OriginHost(objectInfo.Bucket))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})
	// 1.通过响应体获取对象

	resp, err := c.Object.Get(ctx, objectInfo.Key, nil)
	if err != nil {
		return err
	}
	bs, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
	// 2.获取对象到本地文件
	_, err = c.Object.GetToFile(ctx, objectInfo.Key, toFilePath, nil)
	if err != nil {
		return err
	}

	return nil
}

// 获取存储域
func (cs *Storage) GetDomain() string {
	return config.CfgData.Storage.COS.Schema + "://" + config.CfgData.Storage.COS.CdnHost
}

func (cs *Storage) GetURL(objectKey string) string {
	return cosroute.URLForKey(objectKey)
}

// GetPresignedURL 获取可在浏览器中访问的私有预签名的URL
func (cs *Storage) GetPresignedURL(ctx context.Context, fileUrl string) (string, error) {
	objectInfo, err := cosroute.ParseObjectURL(fileUrl)
	if err != nil {
		return "", err
	}
	// 初始化腾讯云 COS 客户端
	u, _ := url.Parse(config.CfgData.Storage.COS.Schema + "://" + cosroute.OriginHost(objectInfo.Bucket))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.CfgData.Storage.COS.SecretID,
			SecretKey: config.CfgData.Storage.COS.SecretKey,
		},
	})
	// 1.通过响应体获取对象
	expiredTime := time.Now().Add(time.Hour * 24 * 365).Unix()
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodGet, objectInfo.Key, config.CfgData.Storage.COS.SecretID, config.CfgData.Storage.COS.SecretKey, time.Until(time.Unix(expiredTime, 0)), nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
