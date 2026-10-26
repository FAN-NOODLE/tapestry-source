package blob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"micro-scrm/internal/config"
	"net/url"
	"os"
)

const BLOB = "blob"

type Storage struct{}

// 上传微软云
func (cs *Storage) Upload(ctx context.Context, filePath string, fileName string) (string, error) {
	accountName := config.CfgData.Storage.Blob.AccountName
	accountKey := config.CfgData.Storage.Blob.AccountKey
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	bucket := config.CfgData.Storage.Blob.Bucket
	rawUrl := fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s/", accountName, bucket)
	URL, _ := url.Parse(rawUrl)
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	blobURL := containerURL.NewBlockBlobURL(fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (cs *Storage) Download(ctx context.Context, fromFilePath string, toFilePath string) error {
	return nil
}

// 获取存储域
func (cs *Storage) GetDomain() string {
	return config.CfgData.Storage.Blob.Domain
}

func (cs *Storage) GetURL(objectKey string) string {
	return cs.GetDomain() + "/" + objectKey
}
