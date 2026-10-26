package cos

import (
	"errors"
	"fmt"
	cos2 "github.com/tencentyun/cos-go-sdk-v5"
	"micro-download/internal/config"
	"net/http"
	url2 "net/url"
)

var cos *cos2.Client

func Init(cfg config.QCos) error {
	bUrl := fmt.Sprintf("%s://%s.cos.%s.myqcloud.com", cfg.Schema, cfg.Bucket, cfg.Region)
	bucket, err := url2.Parse(bUrl)
	if err != nil {
		return err
	}
	url := &cos2.BaseURL{
		BucketURL: bucket,
	}
	cos = cos2.NewClient(url, &http.Client{
		Transport: &cos2.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	if cos == nil {
		return errors.New("cos初始化失败")
	}
	return nil
}

func Cos() *cos2.Client {
	return cos
}
