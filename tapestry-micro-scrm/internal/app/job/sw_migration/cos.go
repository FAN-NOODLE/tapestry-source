package swMigrationJob

import (
	"context"
	"fmt"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"io"
	"micro-scrm/internal/app"
	sourceDao "micro-scrm/internal/app/dao/source"
	"micro-scrm/internal/config"
	"micro-scrm/internal/modules/marketingPromotion/storage"
	"micro-scrm/internal/modules/marketingPromotion/storage/cos"
	"micro-scrm/internal/pkg/helper"
	"micro-scrm/internal/pkg/xjob"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type MigrationCosUrl struct {
	name         string
	enterpriseId int64
}

func NewMigrationCosUrl() *MigrationCosUrl {
	return &MigrationCosUrl{
		name:         "迁移cos",
		enterpriseId: int64(helper.RuleThree(app.Cfg().Env == "uat", 3373738709328576, 3374939108972352)),
	}
}

func (job *MigrationCosUrl) Run(ctx ctx.Context, handler *xjob.JobHandler) error {
	//处理图片
	job.batchUpdateImage(ctx)
	//处理视频
	job.batchUpdateVideo(ctx)
	return nil
}

func (job *MigrationCosUrl) batchUpdateFile(ctx ctx.Context) {
	thisDao := sourceDao.NewFileDao()
	var nowUnix = time.Now().Unix()
	//获取素材url
	var signId int64
	var total, cosErrTotal, updateErrTotal int
	for {
		list, _ := thisDao.GetManyLimit(ctx, job.enterpriseId, 1000, "id > ?", signId)
		if len(list) == 0 {
			break
		}
		for _, image := range list {
			if image.URL == "" {
				continue
			}
			newUrl, err := job.migrateFromCOSURL(ctx, image.URL)
			if err != nil {
				cosErrTotal++
				fmt.Printf("cost上传错误,URL%s,err:%+v", image.URL, err)
			} else {
				err = thisDao.Update(ctx, map[string]any{"id": image.ID}, map[string]any{
					"url":               newUrl,
					"media_id":          "",
					"media_expire_time": 0,
					"media_status":      0,
					"media_fail_reason": "",
					"updated_at":        nowUnix,
				})
				if err == nil {
					total++
					continue
				}
				updateErrTotal++
			}
		}
		signId = list[len(list)-1].ID
	}
	fmt.Printf("成功处理image:%d条，cos错误:%d,更新错误%d", total, cosErrTotal, updateErrTotal)
}

func (job *MigrationCosUrl) batchUpdateVideo(ctx ctx.Context) {
	thisDao := sourceDao.NewVideoDao()
	var nowUnix = time.Now().Unix()
	//获取素材url
	var signId int64
	var total, cosErrTotal, updateErrTotal int
	for {
		list, _ := thisDao.GetManyLimit(ctx, job.enterpriseId, 1000, "id > ?", signId)
		if len(list) == 0 {
			break
		}
		for _, image := range list {
			if image.URL == "" {
				continue
			}
			newUrl, err := job.migrateFromCOSURL(ctx, image.URL)
			if err != nil {
				cosErrTotal++
				fmt.Printf("cost上传错误,URL%s,err:%+v", image.URL, err)
			} else {
				err = thisDao.Update(ctx, map[string]any{"id": image.ID}, map[string]any{
					"url":               newUrl,
					"media_id":          "",
					"media_expire_time": 0,
					"media_status":      0,
					"media_fail_reason": "",
					"updated_at":        nowUnix,
				})
				if err == nil {
					total++
					continue
				}
				updateErrTotal++
			}
		}
		signId = list[len(list)-1].ID
	}
	fmt.Printf("成功处理image:%d条，cos错误:%d,更新错误%d", total, cosErrTotal, updateErrTotal)
}
func (job *MigrationCosUrl) batchUpdateImage(ctx ctx.Context) {
	thisDao := sourceDao.NewImageDao()
	var nowUnix = time.Now().Unix()
	//获取素材url
	var signId int64
	var total, cosErrTotal, updateErrTotal int
	for {
		list, _ := thisDao.GetManyLimit(ctx, job.enterpriseId, 1000, "id > ?", signId)
		if len(list) == 0 {
			break
		}
		for _, image := range list {
			if image.URL == "" {
				continue
			}
			newUrl, err := job.migrateFromCOSURL(ctx, image.URL)
			if err != nil {
				cosErrTotal++
				fmt.Printf("cost上传错误,URL%s,err:%+v", image.URL, err)
			} else {
				err = thisDao.Update(ctx, map[string]any{"id": image.ID}, map[string]any{
					"url":               newUrl,
					"media_id":          "",
					"media_expire_time": 0,
					"media_status":      0,
					"media_fail_reason": "",
					"updated_at":        nowUnix,
				})
				if err == nil {
					total++
					continue
				}
				updateErrTotal++
			}
		}
		signId = list[len(list)-1].ID
	}
	fmt.Printf("成功处理image:%d条，cos错误:%d,更新错误%d", total, cosErrTotal, updateErrTotal)
}

// MigrateFromCOSURL COS迁移 - 从一个COS存储桶迁移到另一个，只改变host部分
func (job *MigrationCosUrl) migrateFromCOSURL(ctx context.Context, sourceURL string) (string, error) {
	cs := storage.Init(cos.COS)
	// 解析源URL，提取出路径部分（object key）
	parsedURL, err := url.Parse(sourceURL)
	if err != nil {
		return "", fmt.Errorf("parse source URL failed: %v", err)
	}

	// 获取路径部分（去掉开头的/）
	objectKey := strings.TrimPrefix(parsedURL.Path, "/")
	if objectKey == "" {
		return "", fmt.Errorf("invalid COS URL path: %s", sourceURL)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "cos_migrate_*")
	if err != nil {
		return "", fmt.Errorf("create temp file failed: %v", err)
	}
	tmpFilePath := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFilePath) // 确保清理临时文件
	}()

	// 下载源COS文件到临时文件
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", sourceURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request failed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download from source COS failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// 将内容拷贝到临时文件
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("save to temp file failed: %v", err)
	}

	// 确保内容写入磁盘
	tmpFile.Sync()
	tmpFile.Close()

	// 使用你原有的本地文件上传方法上传到目标COS
	// 这里保持原有的文件路径和文件名
	_, err = cs.Upload(ctx, tmpFilePath, objectKey)
	if err != nil {
		return "", fmt.Errorf("upload to target COS failed: %v", err)
	}

	// 构建目标COS的URL（只改变host部分）
	targetURL := config.CfgData.Storage.COS.Schema + "://" +
		config.CfgData.Storage.COS.Bucket + ".cos." +
		config.CfgData.Storage.COS.Origin + ".myqcloud.com/" +
		objectKey

	return targetURL, nil
}
