package blob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"micro-download/internal/config"
	"net/url"
	"os"
)

const BLOB = "blob"

type Blob struct{}

// 上传微软云
func (c *Blob) Upload(ctx context.Context, filePath string, fileName string) (string, error) {
	accountName := config.CfgData.Blob.AccountName
	accountKey := config.CfgData.Blob.AccountKey
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	bucket := config.CfgData.Blob.Bucket
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

// 获取存储域
func (c *Blob) GetDomain() string {
	return config.CfgData.Blob.Domain
}
